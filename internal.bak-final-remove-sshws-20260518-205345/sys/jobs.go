package sys

import (
	"fmt"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/stevenjosecarcamo-star/DarkZsaid-Go-Bot/internal/db"
	"github.com/stevenjosecarcamo-star/DarkZsaid-Go-Bot/internal/vpn"
	tele "gopkg.in/telebot.v3"
)

// CountZivpnActive returns true if any UDP session exists for zivpn
func CountZivpnActive() bool {
	out, err := exec.Command("sh", "-c", "ss -u -n -p | grep 'zivpn' | wc -l").Output()
	if err != nil {
		return false
	}
	count := strings.TrimSpace(string(out))
	return count != "" && count != "0"
}

// AutoCleanupLoop corre en un hilo separado ejecutando la limpieza de Iptables
// y usuarios excedidos cada cierto tiempo.
func AutoCleanupLoop(b *tele.Bot) {

	tick := 0
	lastUpdateCheck := time.Now()
	for {
		// Revisar límites de conexión activa cada 14 segundos (2 ticks)
		if tick%2 == 0 {
			EnforceConnectionLimits()
		}

		// 1. Limpieza de usuarios vencidos y AutoReboot de forma periódica
		if tick >= 9 { // Cada 60-70 segundos aprox
			// Guardar el tráfico en DB para que persista tras reiniciar la VPS
			GetGlobalTraffic()

			var sshExpired bool
			db.Update(func(data *db.ConfigData) error {
				now := time.Now()
				nowStr := now.Format("2006-01-02")

				// REBOOT AUTOMÁTICO POR UPTIME (24 HORAS)
				if data.AutoReboot {
					outUptime, err := exec.Command("awk", "{print $1}", "/proc/uptime").Output()
					if err == nil {
						uptimeSecStr := strings.TrimSpace(string(outUptime))
						var uptimeSecFloat float64
						// Parsear manualmente de modo simplificado si se necesita.
						// Para evitar dependencias extra si goimports falla, uso format standard.
						// Pero mejor agregar la lógica que será limpiada con goimports:
						fmt.Sscanf(uptimeSecStr, "%f", &uptimeSecFloat)
						
						if uptimeSecFloat >= 86400 {
							go func() {
								time.Sleep(2 * time.Second)
								exec.Command("reboot").Run()
							}()
						}
					}
				}

				// Revisar SSH
				for user, expire := range data.SSHTimeUsers {
					if nowStr > expire {
						DeleteSSHUser(user)
						delete(data.SSHTimeUsers, user)
						delete(data.SSHOwners, user)
						delete(data.SSHLastActive, user)
						delete(data.SSHBannerTitles, user)
						sshExpired = true
					}
				}

				// Revisar ZiVPN - auto-expiración por fecha
				for pass, expire := range data.ZivpnUsers {
					if nowStr > expire {
						vpn.RemoveZivpnUser(pass)
						delete(data.ZivpnUsers, pass)
						delete(data.ZivpnOwners, pass)
						delete(data.ZivpnLastActive, pass)
					}
				}

				// Revisar Xray - auto-expiración por fecha
				for uid, user := range data.XrayUsers {
					if nowStr > user.Expire {
						vpn.RemoveXrayUser(uid)
						delete(data.XrayUsers, uid)
					}
				}

				return nil
			})

			if sshExpired {
				SyncSSHDBanners()
			}

			// === AUTO-UPDATE CHECK (CADA 12 HORAS) ===
			if time.Since(lastUpdateCheck) > 12*time.Hour {
				lastUpdateCheck = time.Now()
				
				var autoUpdateEnabled bool
				db.Update(func(d *db.ConfigData) error {
					autoUpdateEnabled = d.AutoUpdate
					return nil
				})

				if autoUpdateEnabled {
					hasUpdate, _, err := CheckForUpdate()
					if err == nil && hasUpdate {
						// Notificar a los admins y ejecutar actualización
						data, _ := db.Load()
						for adminID := range data.Admins {
							id, _ := strconv.ParseInt(adminID, 10, 64)
							b.Send(&tele.User{ID: id}, "🔄 <b>ACTUALIZACIÓN AUTOMÁTICA</b>\nSe ha detectado una nueva versión del bot en GitHub. Aplicando actualización en segundo plano...")
						}
						RunUpdate()
					}
				}
			}
			// ==========================================

			// Liberar memoria RAM inactiva al Sistema Operativo
			debug.FreeOSMemory()

			// Regenerar banners de usuarios SSH para actualizar días restantes
			RefreshAllBanners()

			// Nueva Ejecución: Limpieza cada 60s terminada
			tick = 0
		}

		// Ejecución Crítica: Cuotas y Conexiones (Más frecuente)
		if tick%2 == 0 {
			EnforceConnectionLimits()
		}

		tick++
		time.Sleep(7 * time.Second)
	}
}
