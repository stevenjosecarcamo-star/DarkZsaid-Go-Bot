package bot

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/db"
	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/vpn"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

func handleCrearXray(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()

	if !data.Xray.Installed {
		return c.Edit("⚠️ <b>Protocolo Inactivo:</b> Primero debes instalar Xray desde el menú de Protocolos.", tele.ModeHTML)
	}

	if !data.PublicAccess && !isAdmin(chatID) {
		return c.Edit("⛔ <b>ACCESO DENEGADO</b>", tele.ModeHTML)
	}

	// Verificar cuota de cuentas VMess (SuperAdmin sin límite)
	if !isSuperAdminID(chatID) {
		maxAccounts := data.GetMaxXrayPublic()
		if isAdmin(chatID) {
			maxAccounts = data.GetMaxXrayAdmin()
		}

		// Contar cuentas existentes de este usuario
		currentCount := 0
		for _, user := range data.XrayUsers {
			if user.Owner == fmt.Sprintf("%d", chatID) {
				currentCount++
			}
		}

		if currentCount >= maxAccounts {
			markup := &tele.ReplyMarkup{}
			markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
			return SafeEditCtx(c, b, fmt.Sprintf("⚠️ <b>Límite Alcanzado</b>\n\nYa tienes <code>%d/%d</code> cuentas VMess activas.\nNo puedes crear más hasta que se elimine o expire alguna.", currentCount, maxAccounts), markup)
		}
	}

	SetUserStep(chatID, "awaiting_xray_alias")
	SetTempData(chatID, make(map[string]string))
	lastMsg := GetLastBotMsg(chatID)

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "cancelar_accion")))

	_, err := SafeEdit(chatID, b, lastMsg, "💎 <b>Crear Usuario VMess (Xray)</b>\n\n👤 <i>Escribe un alias para identificar al usuario:</i>", markup)
	return err
}

func handleManageXrayUsers(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()
	sa, _ := strconv.ParseInt(superAdmin, 10, 64)
	isSA := chatID == sa

	res := "💎 <b>GESTIÓN DE USUARIOS VMESS (XRAY)</b>\n━━━━━━━━━━━━━━\n"
	count := 0
	
	markup := &tele.ReplyMarkup{}
	var rows []tele.Row

	for uid, user := range data.XrayUsers {
		ownerID, _ := strconv.ParseInt(user.Owner, 10, 64)
		if isSA || ownerID == chatID {
			label := fmt.Sprintf("👤 %s (%s)", user.Alias, user.Expire)
			res += fmt.Sprintf("• %s\n<code>%s</code>\n\n", label, uid)
			
			// Botón de eliminación
			btnDel := markup.Data("🗑️ Borrar "+user.Alias, "del_xray_exec", uid)
			rows = append(rows, markup.Row(btnDel))
			count++
		}
	}

	btnBack := markup.Data("🔙 Volver", "submenu_xray")
	rows = append(rows, markup.Row(btnBack))
	markup.Inline(rows...)

	if count == 0 {
		return SafeEditCtx(c, b, "❌ No tienes usuarios VMess activos.", markup)
	}

	res += "━━━━━━━━━━━━━━\n<i>Pulsa sobre un usuario para eliminarlo.</i>"
	return SafeEditCtx(c, b, res, markup)
}

func handleSubMenuXray(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	data, _ := db.Load()
	status := "❌ Desinstalado"
	if data.Xray.Installed {
		status = "✅ Instalado"
	}

	markup := &tele.ReplyMarkup{}
	btnInst := markup.Data("📥 Instalar", "install_xray")
	btnManage := markup.Data("👥 Gestionar Usuarios", "manage_xray_users")
	btnUninst := markup.Data("🗑️ Desinstalar", "uninstall_xray")
	btnBack := markup.Data("🔙 Volver", "menu_protocols")

	// Solo SuperAdmin puede ver Instalar/Desinstalar
	if isSuperAdminID(chatID) {
		if data.Xray.Installed {
			markup.Inline(markup.Row(btnManage), markup.Row(btnUninst), markup.Row(btnBack))
		} else {
			markup.Inline(markup.Row(btnInst), markup.Row(btnBack))
		}
	} else {
		// Admins y públicos solo ven gestión de usuarios
		if data.Xray.Installed {
			markup.Inline(markup.Row(btnManage), markup.Row(btnBack))
		} else {
			markup.Inline(markup.Row(btnBack))
		}
	}

	texto := fmt.Sprintf("💎 <b>Gestión de Xray (VMess)</b>\n\n📊 <b>Estado:</b> %s\n⚙️ <b>Puerto Interno:</b> <code>10002</code>\n\nEste protocolo utiliza <b>VMess sobre WebSocket</b> y está diseñado para funcionar detrás de Cloudflare y HAProxy.\n\n¿Qué deseas hacer?", status)
	return SafeEditCtx(c, b, texto, markup)
}

func processXraySteps(step string, text string, chatID int64, c tele.Context, b *tele.Bot, lastMsg *tele.Message) error {
	markupCancel := &tele.ReplyMarkup{}
	markupCancel.Inline(markupCancel.Row(markupCancel.Data("❌ Cancelar", "cancelar_accion")))

	switch step {
	case "awaiting_xray_alias":
		alias := strings.TrimSpace(text)
		if len(alias) < 3 {
			_, err := SafeEdit(chatID, b, lastMsg, "⚠️ El alias debe tener al menos 3 caracteres.\n👤 <i>Escribe el alias de nuevo:</i>", markupCancel)
			return err
		}
		SetTempValue(chatID, "xray_alias", alias)

		if isSuperAdminID(chatID) {
			SetUserStep(chatID, "awaiting_xray_days")
			_, err := SafeEdit(chatID, b, lastMsg, fmt.Sprintf("✅ Alias <code>%s</code> guardado.\n\n⏳ <i>¿Cuántos días de duración? (ej: 30):</i>", html.EscapeString(alias)), markupCancel)
			return err
		}

		days := 3
		if isAdmin(chatID) {
			days = 30
		}
		return finishXrayCreation(c, b, chatID, lastMsg, alias, days)

	case "awaiting_xray_days":
		days, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil || days <= 0 {
			_, err := SafeEdit(chatID, b, lastMsg, "⚠️ Valor inválido. Escribe un número mayor a 0:\n⏳ <i>Días:</i>", markupCancel)
			return err
		}
		alias := GetTempValue(chatID, "xray_alias")
		return finishXrayCreation(c, b, chatID, lastMsg, alias, days)
	}
	return nil
}

func finishXrayCreation(c tele.Context, b *tele.Bot, chatID int64, lastMsg *tele.Message, alias string, days int) error {
	DeleteUserStep(chatID)
	SafeEdit(chatID, b, lastMsg, "⏳ <i>Generando UUID y configurando Xray...</i>", nil)

	newUUID := uuid.New().String()
	expireDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	// 1. Agregar al sistema core
	err := vpn.AddXrayUser(newUUID, alias)
	if err != nil {
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Volver", "submenu_xray")))
		SafeEdit(chatID, b, lastMsg, "❌ <b>Error al configurar Xray Core:</b>\n"+err.Error(), markup)
		return err
	}

	// 2. Guardar en DB
	db.Update(func(data *db.ConfigData) error {
		if data.XrayUsers == nil {
			data.XrayUsers = make(map[string]db.XrayUser)
		}
		handle := ""
		if c.Sender() != nil && c.Sender().Username != "" {
			handle = "@" + c.Sender().Username
		}
		data.XrayUsers[newUUID] = db.XrayUser{
			Alias:  alias,
			Expire: expireDate,
			Owner:  fmt.Sprintf("%d", chatID),
			Handle: handle,
		}
		return nil
	})

	data, _ := db.Load()
	vmessLink := vpn.GenerateVmessLink(alias, newUUID, data.CloudflareDomain)

	res := "✅ <b>Usuario VMess Creado</b>\n"
	res += "━━━━━━━━━━━━━━\n"
	res += fmt.Sprintf("👤 <b>Alias:</b> <code>%s</code>\n", alias)
	res += fmt.Sprintf("📅 <b>Expira:</b> <code>%s</code>\n", expireDate)
	res += fmt.Sprintf("🌍 <b>Dominio:</b> <code>%s</code>\n", data.CloudflareDomain)
	res += "⚙️ <b>Protocolo:</b> VMess WebSocket\n"
	res += "⚙️ <b>Puerto:</b> 443 (SSL)\n"
	res += "━━━━━━━━━━━━━━\n"
	res += "🔗 <b>LINK DE CONEXIÓN:</b>\n"
	res += fmt.Sprintf("<code>%s</code>\n\n", vmessLink)
	res += "<i>Copia el link y pégalo en v2rayNG, HTTP Custom o NapsternetV.</i>"

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))
	_, err = SafeEdit(chatID, b, lastMsg, res, markup)
	return err
}

func handleDeleteXrayExec(c tele.Context, b *tele.Bot) error {
	uid := c.Data()
	data, _ := db.Load()
	user, exists := data.XrayUsers[uid]
	if !exists {
		return c.Respond(&tele.CallbackResponse{Text: "Usuario no existe.", ShowAlert: true})
	}

	// Borrar del núcleo
	vpn.RemoveXrayUser(uid)

	// Borrar de DB
	db.Update(func(data *db.ConfigData) error {
		delete(data.XrayUsers, uid)
		return nil
	})

	c.Respond(&tele.CallbackResponse{Text: "Usuario " + user.Alias + " eliminado.", ShowAlert: true})
	return handleManageXrayUsers(c, b)
}
