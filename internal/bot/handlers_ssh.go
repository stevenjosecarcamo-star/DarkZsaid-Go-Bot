package bot

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/db"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/sys"
	tele "gopkg.in/telebot.v3"
)

func handleCrearSSH(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()

	if !isSuperAdminID(chatID) {
		maxAccounts := data.GetMaxSSHPublic()
		if isAdmin(chatID) {
			maxAccounts = data.GetMaxSSHAdmin()
		}

		currentCount := 0
		for _, ownerID := range data.SSHOwners {
			if ownerID == fmt.Sprintf("%d", chatID) {
				currentCount++
			}
		}

		if currentCount >= maxAccounts {
			markup := &tele.ReplyMarkup{}
			markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
			return SafeEditCtx(c, b, fmt.Sprintf("⚠️ <b>Límite Alcanzado</b>\n\nYa tienes <code>%d/%d</code> cuentas SSH activas.\nNo puedes crear más hasta que se elimine o expire alguna.", currentCount, maxAccounts), markup)
		}
	}

	// 1. Iniciar registro de estado
	SetUserStep(chatID, "awaiting_ssh_username")
	SetTempData(chatID, make(map[string]string))

	markup := &tele.ReplyMarkup{}
	btnCancel := markup.Data("❌ Cancelar", "cancelar_accion")
	markup.Inline(markup.Row(btnCancel))

	lastMsg := GetLastBotMsg(chatID)
	msg, _ := SafeEdit(chatID, b, lastMsg, "👤 <b>Crear Nuevo Usuario SSH</b>\n\n✏️ <i>Escribe el nombre de usuario que deseas (ej. pepito):</i>", markup)
	SetLastBotMsg(chatID, msg)
	return nil
}

func handleTextInputs(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	text := c.Text()
	step, ok := GetUserStepWithOk(chatID)
	if !ok {
		return nil
	}

	// Borrar el mensaje del usuario de inmediato para mantener el chat limpio (Sink Global)
	_ = c.Delete()

	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))

	// Dispatcher para otros protocolos
	if strings.HasPrefix(step, "awaiting_zivpn_") {
		lastMsg := GetLastBotMsg(chatID)
		return processZivpnSteps(step, text, chatID, c, b, lastMsg)
	}
	if strings.HasPrefix(step, "awaiting_vpn_") || strings.HasPrefix(step, "awaiting_quota_") || strings.HasPrefix(step, "awaiting_rename_") || strings.HasPrefix(step, "awaiting_promo_") || strings.HasPrefix(step, "awaiting_ban_") {
		lastMsg := GetLastBotMsg(chatID)
		return processVPNSteps(step, text, chatID, c, b, lastMsg)
	}
	if strings.HasPrefix(step, "awaiting_scanner_") {
		lastMsg := GetLastBotMsg(chatID)
		return processScannerSteps(step, text, chatID, c, b, lastMsg)
	}
	if strings.HasPrefix(step, "awaiting_xray_") {
		lastMsg := GetLastBotMsg(chatID)
		return processXraySteps(step, text, chatID, c, b, lastMsg)
	}

	lastMsg := GetLastBotMsg(chatID)
	textLower := strings.ToLower(strings.TrimSpace(text))

	// Interceptar comandos de navegación para cancelar estado
	if strings.HasPrefix(text, "/") || textLower == "menu" || textLower == "salir" || textLower == "atrás" || textLower == "atras" || textLower == "cancelar" {
		DeleteUserStep(chatID)
		return handleStart(c, b)
	}

	switch step {
	case "awaiting_ssh_username":
		if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(text) {
			_, err := SafeEdit(chatID, b, lastMsg, "⚠️ El usuario solo puede contener letras, números y guiones bajos.\n✏️ <i>Intenta con otro:</i>", markupCancel)
			return err
		}
		SetTempValue(chatID, "username", text)
		SetUserStep(chatID, "awaiting_ssh_password")
		markupPass := &tele.ReplyMarkup{}
		btnRandom := markupPass.Data("🎲 Generar Aleatoria", "ssh_rnd_pass")
		btnCancel := markupPass.Data("❌ Cancelar", "cancelar_accion")
		markupPass.Inline(markupPass.Row(btnRandom), markupPass.Row(btnCancel))
		_, err := SafeEdit(chatID, b, lastMsg, fmt.Sprintf("✅ Usuario <code>%s</code> guardado.\n\n🔑 <i>Escribe la contraseña:</i>", html.EscapeString(text)), markupPass)
		return err

	case "awaiting_ssh_password":
		SetTempValue(chatID, "password", text)
		if !isSuperAdminID(chatID) {
			data, _ := db.Load()
			if isAdmin(chatID) {
				SetTempValue(chatID, "days", strconv.Itoa(data.GetMaxDaysAdmin()))
				SetTempValue(chatID, "limit", strconv.Itoa(data.GetMaxLimitAdmin()))
			} else {
				SetTempValue(chatID, "days", strconv.Itoa(data.GetMaxDaysPublic()))
				SetTempValue(chatID, "limit", strconv.Itoa(data.GetMaxLimitPublic()))
			}
			if data.SSHBanner != "" {
				return finishSSHCreation(c, b, chatID, lastMsg)
			}
			// Pedir título del banner
			SetUserStep(chatID, "awaiting_ssh_banner_title")
			markupTitle := &tele.ReplyMarkup{}
			btnDefault := markupTitle.Data("🏷️ Título Predeterminado", "ssh_default_title")
			btnCancel2 := markupTitle.Data("❌ Cancelar", "cancelar_accion")
			markupTitle.Inline(markupTitle.Row(btnDefault), markupTitle.Row(btnCancel2))
			_, err := SafeEdit(chatID, b, lastMsg, "🏷️ <b>Título del Banner SSH</b>\n\n✏️ <i>Escribe un título personalizado para el banner del usuario (ej: INTERNET ILIMITADO, SPEED PREMIUM VIP):</i>", markupTitle)
			return err
		}
		SetUserStep(chatID, "awaiting_ssh_days")
		_, err := SafeEdit(chatID, b, lastMsg, "⏳ <i>¿Cuántos días de duración (ej: 30)?</i>", markupCancel)
		return err

	case "awaiting_ssh_days":
		days, err := strconv.Atoi(text)
		if err != nil || days <= 0 {
			_, err := SafeEdit(chatID, b, lastMsg, "⚠️ Valor inválido.\n⏳ <i>Días:</i>", markupCancel)
			return err
		}
		SetTempValue(chatID, "days", text)
		SetUserStep(chatID, "awaiting_ssh_limit")
		_, err = SafeEdit(chatID, b, lastMsg, "💻 <i>Límite de conexiones (0=infinito):</i>", markupCancel)
		return err

	case "awaiting_ssh_limit":
		limit, err := strconv.Atoi(text)
		if err != nil || limit < 0 {
			_, err := SafeEdit(chatID, b, lastMsg, "⚠️ Valor inválido.\n💻 <i>Límite:</i>", markupCancel)
			return err
		}
		SetTempValue(chatID, "limit", text)
		data, _ := db.Load()
		if data.SSHBanner != "" {
			return finishSSHCreation(c, b, chatID, lastMsg)
		}
		// Pedir título del banner
		SetUserStep(chatID, "awaiting_ssh_banner_title")
		markupTitle := &tele.ReplyMarkup{}
		btnDefault := markupTitle.Data("🏷️ Título Predeterminado", "ssh_default_title")
		btnCancel2 := markupTitle.Data("❌ Cancelar", "cancelar_accion")
		markupTitle.Inline(markupTitle.Row(btnDefault), markupTitle.Row(btnCancel2))
		_, err = SafeEdit(chatID, b, lastMsg, "🏷️ <b>Título del Banner SSH</b>\n\n✏️ <i>Escribe un título personalizado para el banner del usuario (ej: INTERNET ILIMITADO, SPEED PREMIUM VIP):</i>", markupTitle)
		return err

	case "awaiting_ssh_banner_title":
		SetTempValue(chatID, "banner_title", strings.ToUpper(strings.TrimSpace(text)))
		return finishSSHCreation(c, b, chatID, lastMsg)

	case "awaiting_broadcast":
		msg := "📢 <b>MENSAJE GLOBAL DE ADMIN:</b>\n\n" + text
		data, _ := db.Load()
		success := 0
		for _, id := range data.UserHistory {
			_, err := b.Send(tele.ChatID(id), msg, tele.ModeHTML)
			if err == nil {
				success++
			}
		}
		DeleteUserStep(chatID)
		return c.Send(fmt.Sprintf("✅ Broadcast enviado a %d usuarios.", success))

	case "awaiting_edit_user_selection":
		user := text
		userData, _ := db.Load()
		sa, _ := strconv.ParseInt(superAdmin, 10, 64)
		if chatID != sa {
			if ownerID, ok := userData.SSHOwners[user]; !ok || ownerID != fmt.Sprintf("%d", chatID) {
				SafeEdit(chatID, b, lastMsg, "❌ <b>No permitido o no existe.</b>\n✏️ <i>Intenta otro:</i>", markupCancel)
				return nil
			}
		} else if _, ok := userData.SSHOwners[user]; !ok {
			SafeEdit(chatID, b, lastMsg, "❌ <b>No existe.</b>\n✏️ <i>Intenta otro:</i>", markupCancel)
			return nil
		}
		SetTempValue(chatID, "edit_target", user)
		SetUserStep(chatID, "") // Clear step but retain TempData for subsequent edits
		return showEditUserMenu(c, b, user)

	case "awaiting_info_cuenta":
		DeleteUserStep(chatID)
		return processInfoCuenta(text, chatID, c, b)

	case "awaiting_edit_pass_val":
		user := GetTempValue(chatID, "edit_target")
		err := sys.UpdateSSHUserPassword(user, text)
		DeleteUserStep(chatID)
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Menú Editar", "menu_editar")))
		if err != nil {
			SafeEdit(chatID, b, lastMsg, "❌ Error: "+err.Error(), markup)
		} else {
			SafeEdit(chatID, b, lastMsg, "✅ Pass cambiado para "+user, markup)
		}
		return nil

	case "awaiting_edit_renew_val":
		user := GetTempValue(chatID, "edit_target")
		days, _ := strconv.Atoi(text)
		err := sys.RenewSSHUser(user, days)
		DeleteUserStep(chatID)
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Menú Editar", "menu_editar")))
		if err != nil {
			SafeEdit(chatID, b, lastMsg, "❌ Error: "+err.Error(), markup)
		} else {
			// Regenerar banner con nueva fecha
			db.Update(func(data *db.ConfigData) error {
				newExpire := time.Now().AddDate(0, 0, days).Format("2006-01-02")
				data.SSHTimeUsers[user] = newExpire
				title := data.SSHBannerTitles[user]
				limit := sys.GetUserMaxLogins(user)
				sys.WriteUserBanner(user, title, limit, newExpire, data)
				sys.SyncSSHDBanners()
				return nil
			})
			SafeEdit(chatID, b, lastMsg, fmt.Sprintf("✅ Renovado %d días para %s", days, user), markup)
		}
		return nil

	case "awaiting_edit_limit_val":
		user := GetTempValue(chatID, "edit_target")
		limit, _ := strconv.Atoi(text)
		err := sys.SetConnectionLimit(user, limit)
		DeleteUserStep(chatID)
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Menú Editar", "menu_editar")))
		if err != nil {
			SafeEdit(chatID, b, lastMsg, "❌ Error: "+err.Error(), markup)
		} else {
			// Regenerar banner con nuevo límite
			data, _ := db.Load()
			if expire, ok := data.SSHTimeUsers[user]; ok {
				title := data.SSHBannerTitles[user]
				sys.WriteUserBanner(user, title, limit, expire, data)
				sys.SyncSSHDBanners()
			}
			SafeEdit(chatID, b, lastMsg, fmt.Sprintf("✅ Límite cambiado para %s", user), markup)
		}
		return nil

	case "awaiting_delete_user_selection":
		return processDeleteSteps(text, chatID, c, b)
	}

	return nil
}

func handleMenuEditar(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()
	sa, _ := strconv.ParseInt(superAdmin, 10, 64)
	isSA := chatID == sa
	res := "✏️ <b>EDITAR USUARIO</b>\n━━━━━━━━━━━━━━\n"
	count := 0
	for user, ownerID := range data.SSHOwners {
		if isSA || ownerID == fmt.Sprintf("%d", chatID) {
			handle := data.SSHHandles[user]
			if handle != "" {
				res += fmt.Sprintf("👤 <code>%s</code> (%s)\n", user, handle)
			} else {
				res += "👤 <code>" + user + "</code>\n"
			}
			count++
		}
	}
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
	if count == 0 {
		return c.Edit("❌ No hay usuarios.", markup, tele.ModeHTML)
	}
	res += "━━━━━━━━━━━━━━\n✏️ Escribe el nombre del usuario:"
	SetUserStep(chatID, "awaiting_edit_user_selection")
	SetTempData(chatID, make(map[string]string))
	return c.Edit(res, markup, tele.ModeHTML)
}

func showEditUserMenu(c tele.Context, b *tele.Bot, user string) error {
	markup := &tele.ReplyMarkup{}
	btnPass := markup.Data("🔑 Pass", "edit_pass")
	btnRenew := markup.Data("📅 Renov", "edit_renew")
	btnLimit := markup.Data("📱 Lim", "edit_limit")
	btnBack := markup.Data("🔙 Volver", "menu_editar")
	markup.Inline(markup.Row(btnPass, btnRenew), markup.Row(btnLimit), markup.Row(btnBack))
	texto := fmt.Sprintf("✏️ <b>EDITAR:</b> <code>%s</code>", user)

	lastMsg := GetLastBotMsg(c.Chat().ID)
	_, err := SafeEdit(c.Chat().ID, b, lastMsg, texto, markup)
	return err
}

func finishSSHCreation(c tele.Context, b *tele.Bot, chatID int64, lastMsg *tele.Message) error {
	// Bloquear estado inmediatamente para evitar spam/carreras
	mData := GetTempData(chatID)
	DeleteUserStep(chatID)

	user := mData["username"]
	pass := mData["password"]
	days, _ := strconv.Atoi(mData["days"])
	limit, _ := strconv.Atoi(mData["limit"])

	SafeEdit(chatID, b, lastMsg, "⏳ <i>Creando cuenta en el sistema...</i>", nil)

	// Crear usuario en el sistema
	err := sys.CreateSSHUser(user, pass, days)
	if err != nil {
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
		SafeEdit(chatID, b, lastMsg, fmt.Sprintf("❌ <b>ERROR:</b>\n<pre>%v</pre>", err), markup)
		return err
	}

	// Aplicar límite
	sys.SetConnectionLimit(user, limit)

	// Guardar en DB
	db.Update(func(data *db.ConfigData) error {
		if data.SSHOwners == nil {
			data.SSHOwners = make(map[string]string)
		}
		data.SSHOwners[user] = fmt.Sprintf("%d", chatID)

		if data.SSHTimeUsers == nil {
			data.SSHTimeUsers = make(map[string]string)
		}
		// Calcular fecha de vencimiento (YYYY-MM-DD)
		expireDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")
		data.SSHTimeUsers[user] = expireDate

		if c.Sender() != nil && c.Sender().Username != "" {
			if data.SSHHandles == nil {
				data.SSHHandles = make(map[string]string)
			}
			data.SSHHandles[user] = "@" + c.Sender().Username
		}

		// Guardar título del banner
		bannerTitle := mData["banner_title"]
		if bannerTitle == "" {
			bannerTitle = "INTERNET ILIMITADO"
		}
		if data.SSHBannerTitles == nil {
			data.SSHBannerTitles = make(map[string]string)
		}
		data.SSHBannerTitles[user] = bannerTitle

		return nil
	})

	// Generar banner individual
	expireDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")
	bannerTitle := mData["banner_title"]
	if bannerTitle == "" {
		bannerTitle = "INTERNET ILIMITADO"
	}
	dataBanner, _ := db.Load()
	sys.WriteUserBanner(user, bannerTitle, limit, expireDate, dataBanner)
	sys.SyncSSHDBanners()

	// Respuesta final
	ip := sys.GetPublicIP()
	dataFinal, _ := db.Load()
	res := "✅ <b>Usuario SSH Creado</b>\n"
	res += "━━━━━━━━━━━━━━\n"
	res += fmt.Sprintf("👤 <b>User:</b> <code>%s</code>\n", user)
	res += fmt.Sprintf("🔑 <b>Pass:</b> <code>%s</code>\n", html.EscapeString(pass))
	res += fmt.Sprintf("⏳ <b>Días:</b> %d\n", days)
	res += fmt.Sprintf("💻 <b>Límite:</b> %d\n", limit)
	res += "━━━━━━━━━━━━━━\n"
	res += fmt.Sprintf("🌐 <b>IP Principal:</b> <code>%s</code>\n\n", ip)

	res += "🔌 <b>PUERTOS SSH ACTIVOS</b>\n"
	res += "• Directo: <code>22</code>\n"
	if dataFinal.Dropbear != "" {
		res += fmt.Sprintf("• Dropbear: <code>%s</code>\n", dataFinal.Dropbear)
	}
	if dataFinal.SSLTunnel != "" {
		res += fmt.Sprintf("• SSL Tunnel (HAProxy): <code>%s, 80, 8080</code>\n", dataFinal.SSLTunnel)
	}
	if dataFinal.Falcon != "" {
		res += fmt.Sprintf("• Falcon Proxy: <code>%s</code>\n", dataFinal.Falcon)
	}
	res += "\n"

	if dataFinal.CloudflareDomain != "" || dataFinal.CloudfrontDomain != "" {
		res += "🌐 <b>CONEXIONES CDN / SNI</b>\n"
		if dataFinal.CloudflareDomain != "" {
			res += fmt.Sprintf("• Cloudflare: <code>%s</code>\n", dataFinal.CloudflareDomain)
		}
		if dataFinal.CloudfrontDomain != "" {
			res += fmt.Sprintf("• Cloudfront: <code>%s</code>\n", dataFinal.CloudfrontDomain)
		}
		res += "\n"
	}

	if dataFinal.SlowDNS.NS != "" {
		res += "🐢 <b>SLOWDNS</b>\n"
		res += fmt.Sprintf("• NS: <code>%s</code>\n", dataFinal.SlowDNS.NS)
		if dataFinal.SlowDNS.Key != "" {
			res += fmt.Sprintf("• Key: <code>%s</code>\n", dataFinal.SlowDNS.Key)
		}
		res += "\n"
	}

	if dataFinal.SSHWebSocket || dataFinal.SSLTunnel != "" {
		res += "🌐 <b>SSH WEBSOCKET</b>\n"
		res += "• WS:  <code>ws://" + ip + ":80</code>\n"
		res += "• WSS: <code>wss://" + ip + ":443</code>\n"
		if dataFinal.SSLTunnel != "" {
			res += "• WS CDN: <code>ws://" + ip + ":8080</code>\n"
		}
		res += "\n"
	}
	res += "━━━━━━━━━━━━━━\n"

	if dataFinal.CloudflareDomain != "" {
		domain := dataFinal.CloudflareDomain
		res += "🚀 <b>PAYLOADS AVANZADOS CLOUDFLARE</b>\n"
		
		res += "<b>1. Normal WS (Puerto 80)</b>\n"
		res += fmt.Sprintf("<code>GET / HTTP/1.1[crlf]Host: %s[crlf]Upgrade: websocket[crlf]Connection: Keep-Alive[crlf][crlf]</code>\n\n", domain)
		
		res += "<b>2. WSS / TLS (Puerto 443 SNI)</b>\n"
		res += fmt.Sprintf("<code>GET wss://%s/ HTTP/1.1[crlf]Host: %s[crlf]Upgrade: websocket[crlf]Connection: Keep-Alive[crlf][crlf]</code>\n\n", domain, domain)
		
		res += "<b>3. HTTP Injector (Modo SNI / Payload)</b>\n"
		res += fmt.Sprintf("<code>[method] [host_port] HTTP/1.1[crlf]Host: %s[crlf]Upgrade: websocket[crlf]Connection: Keep-Alive[crlf][crlf]</code>\n\n", domain)

		res += "🛠 <b>HTTP CUSTOM (Copiar y Pegar)</b>\n"
		res += fmt.Sprintf("• WS HTTP: <code>%s:80@%s:%s</code>\n", domain, user, pass)
		
		if dataFinal.SSLTunnel != "" {
			res += fmt.Sprintf("• WSS (TLS/SSL): <code>%s:443@%s:%s</code>\n", domain, user, pass)
		}
		
		if dataFinal.Dropbear != "" {
			dpPort := strings.Split(dataFinal.Dropbear, ",")[0]
			dpPort = strings.TrimSpace(dpPort)
			res += fmt.Sprintf("• Dropbear: <code>%s:%s@%s:%s</code>\n", domain, dpPort, user, pass)
		}
		
		res += "━━━━━━━━━━━━━━\n"
	}

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
	_, err = SafeEdit(chatID, b, lastMsg, res, markup)
	return err
}

func handleCancel(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	DeleteUserStep(chatID)
	return handleStart(c, b)
}

func handleRandomPass(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	pass := fmt.Sprintf("%06d", 100000+((time.Now().UnixNano()/1000)%900000))
	SetTempValue(chatID, "password", pass)
	lastMsg := GetLastBotMsg(chatID)

	if !isSuperAdminID(chatID) {
		data, _ := db.Load()
		if isAdmin(chatID) {
			SetTempValue(chatID, "days", strconv.Itoa(data.GetMaxDaysAdmin()))
			SetTempValue(chatID, "limit", strconv.Itoa(data.GetMaxLimitAdmin()))
		} else {
			SetTempValue(chatID, "days", strconv.Itoa(data.GetMaxDaysPublic()))
			SetTempValue(chatID, "limit", strconv.Itoa(data.GetMaxLimitPublic()))
		}
		if data.SSHBanner != "" {
			return finishSSHCreation(c, b, chatID, lastMsg)
		}
		// Pedir título del banner
		SetUserStep(chatID, "awaiting_ssh_banner_title")
		markupTitle := &tele.ReplyMarkup{}
		btnDefault := markupTitle.Data("🏷️ Título Predeterminado", "ssh_default_title")
		btnCancel := markupTitle.Data("❌ Cancelar", "cancelar_accion")
		markupTitle.Inline(markupTitle.Row(btnDefault), markupTitle.Row(btnCancel))
		_, err := SafeEdit(chatID, b, lastMsg, "✅ Pass: "+pass+"\n\n🏷️ <b>Título del Banner SSH</b>\n\n✏️ <i>Escribe un título personalizado (ej: INTERNET ILIMITADO):</i>", markupTitle)
		return err
	}

	SetUserStep(chatID, "awaiting_ssh_days")
	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))

	_, err := SafeEdit(chatID, b, lastMsg, "✅ Pass: "+pass+"\n⏳ Días:", markupCancel)
	return err
}

func handleDefaultTitle(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetTempValue(chatID, "banner_title", "INTERNET ILIMITADO")
	lastMsg := GetLastBotMsg(chatID)
	return finishSSHCreation(c, b, chatID, lastMsg)
}

func HandleEditPass(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	user := GetTempValue(chatID, "edit_target")
	SetUserStep(chatID, "awaiting_edit_pass_val")
	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))
	return c.Edit(fmt.Sprintf("🔑 <b>Cambiando Pass:</b> <code>%s</code>\n✏️ Nueva pass:", user), markupCancel, tele.ModeHTML)
}

func HandleEditRenew(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	user := GetTempValue(chatID, "edit_target")
	SetUserStep(chatID, "awaiting_edit_renew_val")
	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))
	return c.Edit(fmt.Sprintf("📅 <b>Renovando:</b> <code>%s</code>\n✏️ ¿Días extra?", user), markupCancel, tele.ModeHTML)
}

func HandleEditLimit(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	user := GetTempValue(chatID, "edit_target")
	SetUserStep(chatID, "awaiting_edit_limit_val")
	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))
	return c.Edit(fmt.Sprintf("📱 <b>Límite:</b> <code>%s</code>\n✏️ Nuevo límite (0=inf):", user), markupCancel, tele.ModeHTML)
}

func handleEditSelection(c tele.Context, b *tele.Bot) error {
	return handleMenuEditar(c, b)
}

func handleDeleteSelection(c tele.Context, b *tele.Bot) error {
	return handleMenuEliminar(c, b)
}

func handleMenuInfoCuenta(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()
	sa, _ := strconv.ParseInt(superAdmin, 10, 64)
	isSA := chatID == sa

	res := "🔍 <b>CONSULTAR ESTADO DE CUENTA</b>\n━━━━━━━━━━━━━━\n"
	count := 0

	// 1. SSH Users
	for user, ownerID := range data.SSHOwners {
		if isSA || ownerID == fmt.Sprintf("%d", chatID) {
			handle := data.SSHHandles[user]
			if handle != "" {
				res += fmt.Sprintf("👤 SSH: <code>%s</code> (%s)\n", user, handle)
			} else {
				res += fmt.Sprintf("👤 SSH: <code>%s</code>\n", user)
			}
			count++
		}
	}

	// 2. ZiVPN Users
	for pass, ownerID := range data.ZivpnOwners {
		if isSA || ownerID == fmt.Sprintf("%d", chatID) {
			handle := data.ZivpnHandles[pass]
			if handle != "" {
				res += fmt.Sprintf("🛰️ ZiVPN: <code>%s</code> (%s)\n", pass, handle)
			} else {
				res += fmt.Sprintf("🛰️ ZiVPN: <code>%s</code>\n", pass)
			}
			count++
		}
	}

	// 3. Xray Users
	for _, user := range data.XrayUsers {
		if isSA || user.Owner == fmt.Sprintf("%d", chatID) {
			if user.Handle != "" {
				res += fmt.Sprintf("💎 Xray: <code>%s</code> (%s)\n", user.Alias, user.Handle)
			} else {
				res += fmt.Sprintf("💎 Xray: <code>%s</code>\n", user.Alias)
			}
			count++
		}
	}

	res += "━━━━━━━━━━━━━━\n"
	if count == 0 {
		res += "<i>No tienes cuentas activas.</i>\n\n"
	}

	res += "✏️ <i>Escribe el nombre de usuario (SSH), contraseña (ZiVPN) o Alias/UUID (Xray):</i>"

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))

	SetUserStep(chatID, "awaiting_info_cuenta")
	return SafeEditCtx(c, b, res, markup)
}

func processInfoCuenta(target string, chatID int64, c tele.Context, b *tele.Bot) error {
	data, _ := db.Load()
	sa, _ := strconv.ParseInt(superAdmin, 10, 64)
	isSA := chatID == sa
	target = strings.TrimSpace(target)

	res := "ℹ️ <b>RESULTADO DE BÚSQUEDA</b>\n"
	res += "━━━━━━━━━━━━━━\n"

	found := false
	ownerID := ""
	accType := ""
	expire := ""
	details := ""

	// 1. Buscar en SSH
	if exp, ok := data.SSHTimeUsers[target]; ok {
		ownerID = data.SSHOwners[target]
		if isSA || ownerID == fmt.Sprintf("%d", chatID) {
			found = true
			accType = "🔒 SSH / Dropbear"
			expire = exp
			limit := sys.GetUserMaxLogins(target)
			details = fmt.Sprintf("👤 <b>Usuario:</b> <code>%s</code>\n💻 <b>Límite:</b> %d", target, limit)
		}
	}

	// 2. Buscar en ZiVPN
	if !found {
		if exp, ok := data.ZivpnUsers[target]; ok {
			ownerID = data.ZivpnOwners[target]
			if isSA || ownerID == fmt.Sprintf("%d", chatID) {
				found = true
				accType = "🛰️ ZiVPN UDP"
				expire = exp
				details = fmt.Sprintf("🔑 <b>Password:</b> <code>%s</code>", target)
			}
		}
	}

	// 3. Buscar en Xray
	if !found {
		for uid, user := range data.XrayUsers {
			if strings.EqualFold(user.Alias, target) || uid == target {
				ownerID = user.Owner
				if isSA || ownerID == fmt.Sprintf("%d", chatID) {
					found = true
					accType = "💎 VMess (Xray)"
					expire = user.Expire
					details = fmt.Sprintf("👤 <b>Alias:</b> <code>%s</code>\n🆔 <b>UUID:</b> <code>%s</code>", user.Alias, uid)
					break
				}
			}
		}
	}

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))

	if !found {
		return SafeEditCtx(c, b, "❌ <b>Cuenta no encontrada o no tienes acceso.</b>", markup)
	}

	// Calcular días restantes
	daysLeft := 0
	parsedExpire, err := time.Parse("2006-01-02", expire)
	if err == nil {
		daysLeft = int(time.Until(parsedExpire).Hours() / 24)
		if daysLeft < 0 {
			daysLeft = 0
		}
	}

	res += fmt.Sprintf("📌 <b>Tipo:</b> %s\n", accType)
	res += details + "\n"
	res += fmt.Sprintf("📅 <b>Vence:</b> <code>%s</code> (%d días restantes)\n", expire, daysLeft)
	if isSA {
		res += fmt.Sprintf("👤 <b>Dueño ID:</b> <code>%s</code>\n", ownerID)
	}
	res += "━━━━━━━━━━━━━━\n"

	return SafeEditCtx(c, b, res, markup)
}
