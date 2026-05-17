package bot

import (
	"fmt"
	"os"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/db"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/drive"
	tele "gopkg.in/telebot.v3"
)

func handleAuthDrive(c tele.Context, b *tele.Bot) error {
	if !isAdmin(c.Chat().ID) {
		return c.Send("⛔ Solo administradores pueden usar este comando.", tele.ModeHTML)
	}

	args := c.Args()
	if len(args) == 0 {
		url, err := drive.GetAuthURL()
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Error obteniendo URL de Google:\n%v\n\n¿Has puesto credentials.json (ID OAuth) en la carpeta del bot?", err))
		}
		texto := "🔐 <b>Autorización Google Drive (Solo 1 vez en la vida)</b>\n\n" +
			"1. Abre este enlace: <a href='" + url + "'>Clica Aquí para dar Permiso</a>\n" +
			"2. Inicia sesión con tu cuenta personal normal.\n" +
			"3. Copia el código que te da la pantalla.\n" +
			"4. Envíamelo de vuelta usando el comando:\n" +
			"<code>/authdrive TU_CODIGO_AQUI</code>"
		return c.Send(texto, tele.ModeHTML)
	}

	code := args[0]
	err := drive.SaveToken(code)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Error guardando el Token:\n%v", err))
	}

	return c.Send("✅ <b>¡Identificación con Google correcta!</b>\n\nEl bot ya tiene los permisos vitalicios guardados en memoria y subirá tus copias hacia TU carpeta de Drive.\n\nPrueba los botones de Ajustes Pro.", tele.ModeHTML)
}

func notifyNotConfigured(c tele.Context) error {
	return c.Respond(&tele.CallbackResponse{
		Text:      "⚠️ No has vinculado tu cuenta con OAuth.\n\nPor favor envía /authdrive al bot para iniciar el enlace con Google antes de usar los Backups.",
		ShowAlert: true,
	})
}

func handleDriveBackup(c tele.Context, b *tele.Bot) error {
	if !drive.IsAuthenticated() {
		return notifyNotConfigured(c)
	}
	
	SafeEditCtx(c, b, "⏳ <i>Subiendo copia de seguridad a tu Drive...\nEl proceso tomará unos segundos.</i>", nil)
	
	err := drive.UploadBackup(db.GetDataPath())
	
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver a Ajustes Pro", "menu_admins")))
	
	if err != nil {
		if err.Error() == "no estas logueado" {
			os.Remove(drive.TokenFile) 
			return SafeEditCtx(c, b, "❌ <b>Token Ausente o Revocado:</b>\nPor favor usa /authdrive para autorizar de nuevo.", markup)
		}
		return SafeEditCtx(c, b, fmt.Sprintf("❌ <b>Error al subir backup:</b>\n%v", err), markup)
	}
	
	return SafeEditCtx(c, b, "✅ <b>Copia de Seguridad subida exitosamente a TU Drive.</b>\nSe creó/actualizó la carpeta BotVPN_Backups.", markup)
}

func handleDriveRestore(c tele.Context, b *tele.Bot) error {
	if !drive.IsAuthenticated() {
		return notifyNotConfigured(c)
	}
	
	SafeEditCtx(c, b, "📥 <i>Descargando y aplicando copia de seguridad desde tu Drive...</i>", nil)
	
	err := drive.RestoreBackup(db.GetDataPath())
	
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver al Inicio", "menu_admins")))
	
	if err != nil {
		return SafeEditCtx(c, b, fmt.Sprintf("❌ <b>Error al restaurar:</b>\n%v", err), markup)
	}
	
	return SafeEditCtx(c, b, "✅ <b>Base de Datos Restaurada Exitosamente!</b>\nLos IDs de usuario y configuraciones se han cargado.\n\n⚠️ <i>Te recomiendo presionar 'Reiniciar VPS' en Ajustes Pro para aplicar configuraciones.</i>", markup)
}
