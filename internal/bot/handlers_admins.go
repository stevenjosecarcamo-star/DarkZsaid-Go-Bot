package bot

import (
	"fmt"
	"os/exec"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/db"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/sys"
	tele "gopkg.in/telebot.v3"
)

func handleMenuAdmins(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	if !isAdmin(chatID) {
		return c.Send("⛔ Solo administradores pueden usar esta función.", tele.ModeHTML)
	}

	data, _ := db.Load()
	accStatus := "🔓 Público"
	if !data.PublicAccess {
		accStatus = "🔒 Privado"
	}

	markup := &tele.ReplyMarkup{}
	btnToggle := markup.Data("🔄 Acceso: "+accStatus, "toggle_public_access")
	btnList := markup.Data("📋 Listar Admins", "list_admins")
	btnAdd := markup.Data("➕ Agregar Admin", "add_admin")
	btnDel := markup.Data("➖ Quitar Admin", "del_admin_menu")
	btnRename := markup.Data("✏️ Renombrar Admin", "rename_admin_menu")
	btnInfo := markup.Data("📝 Editar Info Extra", "edit_extrainfo")
	btnCloudflare := markup.Data("☁️ Cloudflare Domain", "edit_cloudflare")
	btnCloudfront := markup.Data("🚀 Cloudfront Domain", "edit_cloudfront")
	btnBanner := markup.Data("📜 Banner SSH", "edit_banner")
	btnReset := markup.Data("🧹 Limpiar Historial", "reset_history")

	scanPubStatus := "🔓 ON"
	if !data.PublicScanner {
		scanPubStatus = "🔒 OFF"
	}
	btnScanToggle := markup.Data("🔍 Escaner Público: "+scanPubStatus, "toggle_public_scanner")

	btnReboot := markup.Data("🔄 Reiniciar VPS", "reboot_vps_confirm")
	btnAutoReboot := markup.Data("🕒 Auto Reboot", "menu_autoreboot")
	btnBackup := markup.Data("🔄 Backup a Drive", "drive_backup")
	btnRestore := markup.Data("📥 Restaurar Backup", "drive_restore")
	btnBack := markup.Data("🔙 Volver", "back_main")

	btnQuotas := markup.Data("📊 Cuotas Creación", "edit_quotas")
	btnBans := markup.Data("🚫 Gestión Bans", "menu_bans")
	btnUpdater := markup.Data("🔄 Sistema Updater", "menu_updater")

	markup.Inline(
		markup.Row(btnToggle),
		markup.Row(btnList, btnAdd),
		markup.Row(btnDel, btnRename),
		markup.Row(btnInfo),
		markup.Row(btnCloudflare, btnCloudfront),
		markup.Row(btnBanner, btnQuotas),
		markup.Row(btnBans, btnScanToggle),
		markup.Row(btnBackup, btnRestore),
		markup.Row(btnUpdater),
		markup.Row(btnReset),
		markup.Row(btnAutoReboot, btnReboot),
		markup.Row(btnBack),
	)

	texto := "⚙️ <b>CONFIGURACIÓN PRO (ADMIN)</b>\n"
	texto += "━━━━━━━━━━━━━━\n"
	texto += fmt.Sprintf("🛡️ <b>Acceso:</b> %s\n", accStatus)
	texto += fmt.Sprintf("🔍 <b>Escaner Público:</b> %s\n", scanPubStatus)
	texto += fmt.Sprintf("👤 <b>Admins:</b> %d\n", len(data.Admins)+1)
	texto += fmt.Sprintf("👥 <b>Historial:</b> %d IDs\n", len(data.UserHistory))
	texto += fmt.Sprintf("📊 <b>Cuotas Público:</b> %d días / %d disp.\n", data.GetMaxDaysPublic(), data.GetMaxLimitPublic())
	texto += fmt.Sprintf("📊 <b>Cuotas Admin:</b> %d días / %d disp.\n", data.GetMaxDaysAdmin(), data.GetMaxLimitAdmin())
	texto += fmt.Sprintf("💎 <b>VMess Público:</b> %d cuentas max\n", data.GetMaxXrayPublic())
	texto += fmt.Sprintf("💎 <b>VMess Admin:</b> %d cuentas max\n", data.GetMaxXrayAdmin())
	texto += "━━━━━━━━━━━━━━\n"
	texto += "<i>Selecciona una opción avanzada:</i>"

	return SafeEditCtx(c, b, texto, markup)
}

func handleTogglePublicAccess(c tele.Context, b *tele.Bot) error {
	db.Update(func(data *db.ConfigData) error {
		data.PublicAccess = !data.PublicAccess
		return nil
	})
	return handleMenuAdmins(c, b)
}

func handleEditQuotas(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()

	markup := &tele.ReplyMarkup{}
	btnDaysPub := markup.Data(fmt.Sprintf("📅 Días Público: %d", data.GetMaxDaysPublic()), "quota_days_public")
	btnLimitPub := markup.Data(fmt.Sprintf("📱 Disp. Público: %d", data.GetMaxLimitPublic()), "quota_limit_public")
	btnDaysAdm := markup.Data(fmt.Sprintf("📅 Días Admin: %d", data.GetMaxDaysAdmin()), "quota_days_admin")
	btnLimitAdm := markup.Data(fmt.Sprintf("📱 Disp. Admin: %d", data.GetMaxLimitAdmin()), "quota_limit_admin")
	btnSSHPublic := markup.Data(fmt.Sprintf("👤 Max SSH Público: %d", data.GetMaxSSHPublic()), "quota_ssh_public")
	btnSSHAdmin := markup.Data(fmt.Sprintf("👤 Max SSH Admin: %d", data.GetMaxSSHAdmin()), "quota_ssh_admin")
	btnZivpnPublic := markup.Data(fmt.Sprintf("🛰️ Max ZiVPN Público: %d", data.GetMaxZivpnPublic()), "quota_zivpn_public")
	btnZivpnAdmin := markup.Data(fmt.Sprintf("🛰️ Max ZiVPN Admin: %d", data.GetMaxZivpnAdmin()), "quota_zivpn_admin")
	btnXrayPub := markup.Data(fmt.Sprintf("💎 VMess Público: %d", data.GetMaxXrayPublic()), "quota_xray_public")
	btnXrayAdm := markup.Data(fmt.Sprintf("💎 VMess Admin: %d", data.GetMaxXrayAdmin()), "quota_xray_admin")
	btnBack := markup.Data("🔙 Volver", "menu_admins")

	markup.Inline(
		markup.Row(btnDaysPub, btnLimitPub),
		markup.Row(btnDaysAdm, btnLimitAdm),
		markup.Row(btnSSHPublic, btnSSHAdmin),
		markup.Row(btnZivpnPublic, btnZivpnAdmin),
		markup.Row(btnXrayPub, btnXrayAdm),
		markup.Row(btnBack),
	)

	texto := "📊 <b>Cuotas de Creación de Usuarios</b>\n"
	texto += "━━━━━━━━━━━━━━\n"
	texto += fmt.Sprintf("👥 <b>Público SSH (Params):</b> %d días / %d dispositivos\n", data.GetMaxDaysPublic(), data.GetMaxLimitPublic())
	texto += fmt.Sprintf("👤 <b>Admin SSH (Params):</b> %d días / %d dispositivos\n", data.GetMaxDaysAdmin(), data.GetMaxLimitAdmin())
	texto += "━━━━━━━━━━━━━━\n"
	texto += fmt.Sprintf("👤 <b>Límite Cuentas SSH Público:</b> máx %d\n", data.GetMaxSSHPublic())
	texto += fmt.Sprintf("👤 <b>Límite Cuentas SSH Admin:</b> máx %d\n", data.GetMaxSSHAdmin())
	texto += "━━━━━━━━━━━━━━\n"
	texto += fmt.Sprintf("🛰️ <b>Límite Cuentas ZiVPN Público:</b> máx %d\n", data.GetMaxZivpnPublic())
	texto += fmt.Sprintf("🛰️ <b>Límite Cuentas ZiVPN Admin:</b> máx %d\n", data.GetMaxZivpnAdmin())
	texto += "━━━━━━━━━━━━━━\n"
	texto += fmt.Sprintf("💎 <b>VMess Público:</b> máx %d cuentas\n", data.GetMaxXrayPublic())
	texto += fmt.Sprintf("💎 <b>VMess Admin:</b> máx %d cuentas\n", data.GetMaxXrayAdmin())
	texto += "━━━━━━━━━━━━━━\n"
	texto += "<i>Estos valores se aplican al crear usuarios SSH, ZiVPN y VMess.\nEl SuperAdmin no tiene límites.</i>"

	return SafeEditCtx(c, b, texto, markup)
}

func handleQuotaPrompt(c tele.Context, b *tele.Bot, step string, label string) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, step)
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_quotas")))
	return SafeEditCtx(c, b, fmt.Sprintf("✏️ <b>%s</b>\n\n<i>Escribe el nuevo valor (número):</i>", label), markup)
}

func handleListAdmins(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	res := "📋 <b>LISTADO DE ADMINISTRADORES</b>\n\n"
	res += fmt.Sprintf("⭐ <b>SuperAdmin (Root):</b> <code>%s</code>\n", superAdmin)

	if len(data.Admins) == 0 {
		res += "\n<i>No hay administradores adicionales.</i>"
	} else {
		i := 1
		for id, info := range data.Admins {
			res += fmt.Sprintf("\n%d. 👤 <b>%s</b>\n   └ ID: <code>%s</code>\n", i, info.Alias, id)
			i++
		}
	}

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "menu_admins")))
	return SafeEditCtx(c, b, res, markup)
}

func handleAddAdminPrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_vpn_admin_id")

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_admins")))

	return SafeEditCtx(c, b, "➕ <b>Agregar Nuevo Administrador</b>\n\n📝 <b>Paso 1/2:</b> Escribe el <b>ID numérico</b> del usuario de Telegram:\n\nEjemplo: <code>123456789</code>", markup)
}

func handleDelAdminMenu(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	if len(data.Admins) == 0 {
		return c.Respond(&tele.CallbackResponse{Text: "No hay administradores para quitar.", ShowAlert: true})
	}

	markup := &tele.ReplyMarkup{}
	var rows []tele.Row
	for id, info := range data.Admins {
		rows = append(rows, markup.Row(markup.Data("❌ "+info.Alias+" ("+id+")", "del_adm_exec", id)))
	}
	rows = append(rows, markup.Row(markup.Data("🔙 Volver", "menu_admins")))
	markup.Inline(rows...)

	return SafeEditCtx(c, b, "➖ <b>Quitar Administrador</b>\n\nSelecciona a quién deseas retirar los permisos:", markup)
}

func handleDelAdminExec(c tele.Context, b *tele.Bot) error {
	id := c.Data()

	// Buscar alias antes de borrar
	data, _ := db.Load()
	alias := "Admin"
	if info, ok := data.Admins[id]; ok {
		alias = info.Alias
	}

	db.Update(func(data *db.ConfigData) error {
		delete(data.Admins, id)
		return nil
	})

	// Responder al callback para desbloquear el botón
	c.Respond(&tele.CallbackResponse{Text: "✅ Admin eliminado"})

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver a Ajustes", "menu_admins")))
	return SafeEditCtx(c, b, fmt.Sprintf("✅ <b>Admin Eliminado</b>\n\n👤 <b>%s</b>\n🆔 ID: <code>%s</code>\n\n<i>Ya no tiene permisos de administrador.</i>", alias, id), markup)
}

func handleRenameAdminMenu(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	if len(data.Admins) == 0 {
		return c.Respond(&tele.CallbackResponse{Text: "No hay administradores para renombrar.", ShowAlert: true})
	}

	markup := &tele.ReplyMarkup{}
	var rows []tele.Row
	for id, info := range data.Admins {
		rows = append(rows, markup.Row(markup.Data("✏️ "+info.Alias+" ("+id+")", "rename_adm_sel", id)))
	}
	rows = append(rows, markup.Row(markup.Data("🔙 Volver", "menu_admins")))
	markup.Inline(rows...)

	return SafeEditCtx(c, b, "✏️ <b>Renombrar Administrador</b>\n\nSelecciona al admin que deseas renombrar:", markup)
}

func handleRenameAdminSelect(c tele.Context, b *tele.Bot) error {
	id := c.Data()
	chatID := c.Chat().ID

	data, _ := db.Load()
	info, exists := data.Admins[id]
	if !exists {
		return c.Respond(&tele.CallbackResponse{Text: "Admin no encontrado.", ShowAlert: true})
	}

	SetTempValue(chatID, "rename_admin_id", id)
	SetUserStep(chatID, "awaiting_rename_admin_alias")

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_admins")))

	return SafeEditCtx(c, b, fmt.Sprintf("✏️ <b>Renombrar Admin</b>\n\n👤 <b>Actual:</b> %s\n🆔 <b>ID:</b> <code>%s</code>\n\nEscribe el <b>nuevo alias</b>:", info.Alias, id), markup)
}

func handleEditExtraInfoPrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_vpn_extrainfo")

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_admins")))

	return SafeEditCtx(c, b, "📝 <b>Editar Información Extra</b>\n\nEsta información aparecerá en el menú /info.\n\n✏️ <i>Escribe el nuevo texto (soporta HTML):</i>", markup)
}

func handleEditCloudflarePrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_vpn_cloudflare")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_admins")))
	return SafeEditCtx(c, b, "☁️ <b>Configurar Dominio Cloudflare</b>\n\n✏️ <i>Escribe el dominio :</i>\n\nEjemplo: <code>mi.host.com</code>", markup)
}

func handleEditCloudfrontPrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_vpn_cloudfront")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_admins")))
	return SafeEditCtx(c, b, "🚀 <b>Configurar Dominio Cloudfront</b>\n\n✏️ <i>Escribe el dominio:</i>\n\nEjemplo: <code>xyz123.cloudfront.net</code>", markup)
}

// Banner predeterminado de DarkZsaid
const defaultBanner = `<html>
<h5 style="text-align:center;">
<font face="monospace" color="#00ff00">
⠀⠀⢀⣶⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣶⡀⠀⠀
⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀
⠀⠀⢸⣿⡇⠀⠀⠀⣠⣶⣄⠀⠀⠀⢸⣿⡇⠀⠀
⠀⠀⢸⣿⡇⠀⠀⢰⣿⣿⣿⡆⠀⠀⢸⣿⡇⠀⠀
⠀⠀⠈⣿⣿⡄⢀⣿⣿⠻⣿⣿⡀⢠⣿⣿⠁⠀⠀
⠀⠀⠀⠹⣿⣿⣾⣿⡏⠀⢹⣿⣷⣿⣿⠏⠀⠀⠀
⠀⠀⠀⠀⠙⢿⣿⡿⠀⠀⠀⢿⣿⡿⠋⠀⠀⠀⠀
</font>
</h5>
<h1 style="text-align:center;">
<font face="monospace" color="#00ff00"><b>DARKZSAID</b></font>
</h1>
<h5 style="text-align:center;">
<font color='#29b6f6'>==============================</font>
<font color='#29b6f6'><b>✈ TELEGRAM ✈</b></font>
<font color='#29b6f6'>==============================</font>
</h5>
<h5 style="text-align:center;">
<font color='#ffffff'>Dev: </font><a href="https://t.me/Dan3651"><font color='#f1c40f'>@DarkZsaid</font></a>
<font color='#ffffff'>Canal: </font><a href="https://t.me/DarkZsaid2"><font color='#f1c40f'>@DarkZsaid2</font></a>
</h5>
<h4 style="text-align:center;">
<font color='#FF00FF'><b>🔥 ¡SE VENDEN SERVIDORES PREMIUM 35 DÍAS A 8.5 SOLES! 🔥</b></font>
</h4>
<h5 style="text-align:center;">
<font color='#ff0000'>==============================</font>
<font color='#ff0000'><b>⚡ SERVIDORES FREE ⚡</b></font>
<font color='#ff0000'>==============================</font>
</h5>
<h6 style="text-align:center;">
<font color='#ff9800'><b>⚠️ REGLAS DEL SERVIDOR ⚠️</b></font>
<font color='#ffffff'>🚫 NO Torrent / P2P</font>
<font color='#ffffff'>🚫 NO Spam / Fraude</font>
<font color='#ffffff'>🚫 NO Ataques DDoS</font>
<font color='#ff5252'><i>El incumplimiento genera ban automático</i></font>
</h6>
<h5 style="text-align:center;">
<font color='#00e676'><b>CREADO EN : @DarkZsaid_bot</b></font>
</h5>
</html>`

func handleEditBannerPrompt(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()

	status := "👤 Banners Individuales (Activo)"
	bannerType := ""
	if data.SSHBanner != "" {
		status = "🌐 Banner Global (Activo)"
		bannerType = "\n\n⚠️ <i>El sistema individual está desactivado. Todas las cuentas usarán el mismo banner global.</i>"
	} else {
		bannerType = "\n\n✅ <i>Cada usuario tiene su propio banner con días y límites.</i>"
	}

	markup := &tele.ReplyMarkup{}
	btnPromo := markup.Data("📝 Editar Textos Promo", "edit_promo_menu")
	btnCustom := markup.Data("🌐 Activar Banner Global", "banner_set_custom")
	btnDeactivate := markup.Data("🚫 Desactivar Global (Usar Individual)", "banner_deactivate")
	btnBack := markup.Data("🔙 Volver", "menu_admins")

	markup.Inline(
		markup.Row(btnPromo),
		markup.Row(btnCustom),
		markup.Row(btnDeactivate),
		markup.Row(btnBack),
	)

	texto := fmt.Sprintf("📜 <b>Gestión de Banners SSH</b>\n\n📊 <b>Modo Actual:</b> %s%s\n\n¿Qué deseas hacer?", status, bannerType)
	return SafeEditCtx(c, b, texto, markup)
}

func handleEditPromoMenu(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()

	promoText := "🔥 ¡SERVIDORES PREMIUM A 8.5 SOLES! 🔥"
	if data.BannerPromoText != "" {
		promoText = data.BannerPromoText
	}

	promoChannel := "@DarkZsaid2"
	if data.BannerPromoChannel != "" {
		promoChannel = data.BannerPromoChannel
	}

	promoSupport := "@DarkZsaid"
	if data.BannerPromoSupport != "" {
		promoSupport = data.BannerPromoSupport
	}

	promoBotName := "@DarkZsaid_bot"
	if data.BannerPromoBotName != "" {
		promoBotName = data.BannerPromoBotName
	}

	markup := &tele.ReplyMarkup{}
	btnText := markup.Data("📝 Editar Mensaje", "edit_promo_text")
	btnChannel := markup.Data("📢 Editar Canal", "edit_promo_channel")
	btnSupport := markup.Data("👤 Editar Soporte", "edit_promo_support")
	btnBotName := markup.Data("🤖 Editar Nombre Bot", "edit_promo_botname")
	btnBack := markup.Data("🔙 Volver", "edit_banner")

	markup.Inline(
		markup.Row(btnText, btnChannel),
		markup.Row(btnSupport, btnBotName),
		markup.Row(btnBack),
	)

	texto := "📝 <b>Editar Textos Promocionales (Banners Individuales)</b>\n\n"
	texto += "Estos textos aparecerán en la parte inferior de los banners de cada usuario.\n\n"
	texto += fmt.Sprintf("💬 <b>Mensaje Promo:</b>\n<code>%s</code>\n\n", promoText)
	texto += fmt.Sprintf("📢 <b>Canal:</b>\n<code>%s</code>\n\n", promoChannel)
	texto += fmt.Sprintf("👤 <b>Soporte:</b>\n<code>%s</code>\n\n", promoSupport)
	texto += fmt.Sprintf("🤖 <b>Creado En:</b>\n✅ CREADO EN : <code>%s</code>", promoBotName)

	return SafeEditCtx(c, b, texto, markup)
}

func handleBannerSetCustom(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_vpn_ssh_banner")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_banner")))
	return SafeEditCtx(c, b, "📜 <b>Banner SSH Personalizado</b>\n\n✏️ <i>Escribe el texto del banner (admite HTML básico):</i>\n\nEsto se mostrará al conectar por SSH.", markup)
}

func handleEditPromoText(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_promo_text")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_promo_menu")))
	return SafeEditCtx(c, b, "💬 <b>Editar Mensaje Promo</b>\n\n✏️ <i>Escribe el nuevo texto promocional (ej: 🔥 ¡OFERTA SERVIDORES A 5$! 🔥):</i>", markup)
}

func handleEditPromoChannel(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_promo_channel")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_promo_menu")))
	return SafeEditCtx(c, b, "📢 <b>Editar Canal Promo</b>\n\n✏️ <i>Escribe el @usuario de tu canal (ej: @MiCanalVIP):</i>", markup)
}

func handleEditPromoSupport(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_promo_support")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_promo_menu")))
	return SafeEditCtx(c, b, "👤 <b>Editar Soporte Promo</b>\n\n✏️ <i>Escribe tu @usuario de Telegram para soporte (ej: @TuUsuario):</i>", markup)
}

func handleEditPromoBotName(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_promo_botname")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "edit_promo_menu")))
	return SafeEditCtx(c, b, "🤖 <b>Editar Nombre del Bot</b>\n\n✏️ <i>Escribe el @usuario de tu bot (ej: @MiSuperVPN_bot):</i>\n\nEl banner mantendrá el prefijo \"✅ CREADO EN : \".", markup)
}

func handleBannerDeactivate(c tele.Context, b *tele.Bot) error {
	db.Update(func(data *db.ConfigData) error {
		data.SSHBanner = ""
		return nil
	})

	// Quitar banner global del sistema
	exec.Command("sh", "-c", "rm -f /etc/sshd_banner").Run()
	exec.Command("sed", "-i", "/^Banner/d", "/etc/ssh/sshd_config").Run()

	// Restaurar banners individuales (Match User)
	go sys.RefreshAllBanners()

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "edit_banner")))
	return SafeEditCtx(c, b, "✅ <b>Banner Global desactivado.</b>\n\n<i>Se ha vuelto al sistema de banners individuales.</i>", markup)
}

func handleResetHistoryConfirm(c tele.Context, b *tele.Bot) error {
	markup := &tele.ReplyMarkup{}
	btnYes := markup.Data("✅ Sí, Limpiar", "reset_history_exec")
	btnNo := markup.Data("❌ No, Cancelar", "menu_admins")
	markup.Inline(markup.Row(btnYes, btnNo))

	return SafeEditCtx(c, b, "⚠️ <b>¿Estás seguro de limpiar el historial?</b>\n\nSe borrarán todos los IDs de usuarios registrados (el broadcast ya no les llegará hasta que vuelvan a iniciar el bot).", markup)
}

func handleResetHistoryExec(c tele.Context, b *tele.Bot) error {
	db.Update(func(data *db.ConfigData) error {
		data.UserHistory = []int64{}
		return nil
	})
	return c.Respond(&tele.CallbackResponse{Text: "Historial de IDs reseteado.", ShowAlert: true})
}

func handleServerRebootConfirm(c tele.Context, b *tele.Bot) error {
	markup := &tele.ReplyMarkup{}
	btnYes := markup.Data("🔄 Reiniciar AHORA", "reboot_vps_exec")
	btnNo := markup.Data("🔙 Cancelar", "menu_admins")
	markup.Inline(markup.Row(btnYes, btnNo))

	return SafeEditCtx(c, b, "🚨 <b>ADVERTENCIA: REINICIO DEL SERVIDOR</b>\n\n¿Estás seguro de que quieres reiniciar la VPS? Todas las conexiones actuales se cortarán.", markup)
}

func handleServerRebootExec(c tele.Context, b *tele.Bot) error {
	c.Edit("⏳ <b>Reiniciando VPS...</b> el bot estará offline unos minutos.", tele.ModeHTML)
	exec.Command("reboot").Run()
	return nil
}

// === SISTEMA DE ACTUALIZACIONES (UPDATER) ===

func handleMenuUpdater(c tele.Context, b *tele.Bot) error {
	if !isAdmin(c.Chat().ID) {
		return c.Send("⛔ Solo administradores pueden usar esta función.")
	}

	data, _ := db.Load()
	autoStatus := "🔴 Desactivada"
	if data.AutoUpdate {
		autoStatus = "🟢 Activada"
	}

	text := "🔄 <b>Sistema de Actualizaciones (GitHub)</b>\n\n"
	text += "Versión Actual: <b>" + sys.CurrentVersion + "</b>\n"
	text += "Auto-Actualización: <b>" + autoStatus + "</b>\n\n"
	text += "Puedes buscar si hay una nueva versión disponible o activar la actualización automática (el bot revisará cada 12 horas)."

	markup := &tele.ReplyMarkup{}
	btnCheck := markup.Data("🔍 Buscar Actualización", "updater_check")
	btnAuto := markup.Data("⚙️ Auto-Update: "+autoStatus, "updater_toggle_auto")
	btnForce := markup.Data("⚠️ Forzar Reinstalación (Dev)", "updater_run")
	btnBack := markup.Data("🔙 Volver a Ajustes", "menu_admins")

	markup.Inline(
		markup.Row(btnCheck),
		markup.Row(btnAuto),
		markup.Row(btnForce),
		markup.Row(btnBack),
	)

	return SafeEditCtx(c, b, text, markup)
}

func handleUpdaterToggleAuto(c tele.Context, b *tele.Bot) error {
	if !isAdmin(c.Chat().ID) {
		return nil
	}

	db.Update(func(d *db.ConfigData) error {
		d.AutoUpdate = !d.AutoUpdate
		return nil
	})

	return handleMenuUpdater(c, b)
}

func handleUpdaterCheck(c tele.Context, b *tele.Bot) error {
	if !isAdmin(c.Chat().ID) {
		return nil
	}

	hasUpdate, newVer, err := sys.CheckForUpdate()

	markup := &tele.ReplyMarkup{}
	btnBack := markup.Data("🔙 Volver", "menu_updater")

	if err != nil {
		markup.Inline(markup.Row(btnBack))
		return SafeEditCtx(c, b, "❌ <b>Error al buscar actualizaciones:</b>\n"+err.Error(), markup)
	}

	if !hasUpdate {
		btnForceNow := markup.Data("⚠️ Forzar Reinstalación", "updater_run")
		markup.Inline(
			markup.Row(btnForceNow),
			markup.Row(btnBack),
		)
		return SafeEditCtx(c, b, "✅ <b>Estás en la última versión.</b>\nVersión actual: "+sys.CurrentVersion+"\nVersión remota: "+newVer, markup)
	}

	btnUpdateNow := markup.Data("⚡ Actualizar a v"+newVer, "updater_run")
	markup.Inline(
		markup.Row(btnUpdateNow),
		markup.Row(btnBack),
	)

	return SafeEditCtx(c, b, "🎉 <b>¡Nueva actualización encontrada!</b>\n\nVersión actual: "+sys.CurrentVersion+"\nNueva versión: <b>"+newVer+"</b>\n\n¿Deseas actualizar el bot ahora mismo? El servicio se reiniciará por unos 15 segundos.", markup)
}

func handleUpdaterRun(c tele.Context, b *tele.Bot) error {
	if !isAdmin(c.Chat().ID) {
		return nil
	}

	c.Send("⚡ <b>Iniciando actualización...</b>\nDescargando y compilando desde GitHub. El bot no responderá durante unos 15 segundos.", tele.ModeHTML)
	
	err := sys.RunUpdate()
	if err != nil {
		return c.Send("❌ Error al iniciar el actualizador: " + err.Error())
	}
	return nil
}

func handleTogglePublicScanner(c tele.Context, b *tele.Bot) error {
	db.Update(func(data *db.ConfigData) error {
		data.PublicScanner = !data.PublicScanner
		return nil
	})
	return handleMenuAdmins(c, b)
}

func handleAutoRebootMenu(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	status := "❌ Desactivado"
	if data.AutoReboot {
		status = "✅ Activado"
	}

	markup := &tele.ReplyMarkup{}
	btnToggle := markup.Data("🔄 Switch: "+status, "toggle_autoreboot")
	btnBack := markup.Data("🔙 Volver", "menu_admins")

	markup.Inline(
		markup.Row(btnToggle),
		markup.Row(btnBack),
	)

	texto := "🕒 <b>CONFIGURACIÓN DE AUTO-REINICIO</b>\n"
	texto += "━━━━━━━━━━━━━━\n"
	texto += "<i>El servidor se reiniciará automáticamente cuando alcance 24 Horas de Uptime continuo.</i>\n\n"
	texto += fmt.Sprintf("📊 <b>Estado:</b> %s\n", status)
	texto += "━━━━━━━━━━━━━━\n"
	texto += "<i>Selecciona una opción:</i>"

	return SafeEditCtx(c, b, texto, markup)
}

func handleToggleAutoReboot(c tele.Context, b *tele.Bot) error {
	db.Update(func(data *db.ConfigData) error {
		data.AutoReboot = !data.AutoReboot
		return nil
	})
	return handleAutoRebootMenu(c, b)
}

func handleMenuBans(c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	markup := &tele.ReplyMarkup{}
	
	btnBanUser := markup.Data("➕ Banear Usuario", "ban_user_prompt")
	btnBack := markup.Data("🔙 Volver", "menu_admins")
	
	var rows []tele.Row
	rows = append(rows, markup.Row(btnBanUser))
	
	texto := "🚫 <b>GESTIÓN DE USUARIOS BANEADOS</b>\n━━━━━━━━━━━━━━\n"
	if len(data.BannedUsers) == 0 {
		texto += "<i>No hay usuarios baneados.</i>\n\n"
	} else {
		texto += "<i>Selecciona un usuario para quitarle el Ban:</i>\n\n"
		for id, info := range data.BannedUsers {
			rows = append(rows, markup.Row(markup.Data(fmt.Sprintf("✅ Desbanear a %s", info.Name), "unban_user", id)))
			texto += fmt.Sprintf("👤 <b>%s</b>\n🆔 ID: <code>%s</code>\n📝 Motivo: <i>%s</i>\n📅 Fecha: %s\n\n", info.Name, id, info.Reason, info.Date)
		}
	}
	
	rows = append(rows, markup.Row(btnBack))
	markup.Inline(rows...)
	
	return SafeEditCtx(c, b, texto, markup)
}

func handleBanUserPrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_ban_id")
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_bans")))
	return SafeEditCtx(c, b, "➕ <b>Banear Usuario</b>\n\n📝 <b>Paso 1/3:</b> Escribe el <b>ID numérico</b> del usuario de Telegram que deseas banear:", markup)
}

func handleUnbanUser(c tele.Context, b *tele.Bot) error {
	id := c.Data()
	db.Update(func(data *db.ConfigData) error {
		delete(data.BannedUsers, id)
		return nil
	})
	c.Respond(&tele.CallbackResponse{Text: "✅ Usuario desbaneado", ShowAlert: true})
	return handleMenuBans(c, b)
}

