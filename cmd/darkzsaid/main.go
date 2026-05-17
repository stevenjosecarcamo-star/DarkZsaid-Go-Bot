package main

import (
	"log"
	"time"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/bot"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/db"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/drive"
)

func main() {
	log.Println("Iniciando DarkZsaid SSH VPN Manager...")

	// Hilo de vigilancia de Backups (cada hora evalúa independientemente de reboots)
	go func() {
		// Le daremos 1 minuto de retraso al encender para que el bot y la red se estabilicen
		time.Sleep(1 * time.Minute)
		
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for {
			if drive.IsAuthenticated() {
				data, err := db.Load()
				if err == nil {
					needsBackup := false

					if data.DriveLastBackup == "" {
						needsBackup = true // Primer backup en la historia de esta DB
					} else {
						lastBackup, errParse := time.Parse(time.RFC3339, data.DriveLastBackup)
						if errParse == nil {
							// Si han pasado 8 horas o la fecha guardada es bizarramente del futuro
							if time.Since(lastBackup) >= 8*time.Hour || time.Since(lastBackup) < 0 {
								needsBackup = true
							}
						} else {
							// Formato dañado, forzamos backup por precaución
							needsBackup = true
						}
					}

					if needsBackup {
						log.Println("Ejecutando ciclo de respaldo automático de 8 horas en Drive...")
						errUpload := drive.UploadBackup(db.GetDataPath())
						if errUpload != nil {
							log.Printf("❌ Error en backup persistente: %v\n", errUpload)
						} else {
							log.Println("✅ Backup automático subido a Drive correctamente.")
							// Guardar el tiempo exacto
							db.Update(func(d *db.ConfigData) error {
								d.DriveLastBackup = time.Now().Format(time.RFC3339)
								return nil
							})
						}
					}
				}
			}
			<-ticker.C
		}
	}()

	// Iniciar servidor del bot (bloqueante)
	bot.StartBot()
}
