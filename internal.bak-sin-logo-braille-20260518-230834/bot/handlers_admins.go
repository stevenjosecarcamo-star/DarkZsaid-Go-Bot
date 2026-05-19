<font face="monospace" color="#00ff00"><b>DARKZSAID</b></font>
</h1>
<h5 style="text-align:center;">
<font color='#29b6f6'>==============================</font>
<font color='#29b6f6'><b>✈ TELEGRAM ✈</b></font>
<font color='#29b6f6'>==============================</font>
</h5>
<h5 style="text-align:center;">
<font color='#ffffff'>Dev: </font><a href="https://t.me/Dan3651"><font color='#f1c40f'>@DarkZsaid</font></a>
<font color='#ffffff'>Canal: </font><a href="https://t.me/DarkZsaid2"><font color='#f1c40f'>@DarkZsaid</font></a>
</h5>
<h4 style="text-align:center;">
<font color='#FF00FF'><b>🔥 ¡SERVIDORES PREMIUM 200 CÓRDOBAS! 🔥</b></font>
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
<font color='#00e676'><b>CREADO EN : @DarkZsaid</b></font>
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

	promoText := "🔥 ¡SERVIDORES PREMIUM 200 CÓRDOBAS! 🔥"
	if data.BannerPromoText != "" {
		promoText = data.BannerPromoText
	}

	promoChannel := "@DarkZsaid"
	if data.BannerPromoChannel != "" {
		promoChannel = data.BannerPromoChannel
	}

	promoSupport := "@DarkZsaid"
	if data.BannerPromoSupport != "" {
		promoSupport = data.BannerPromoSupport
	}

	promoBotName := "@DarkZsaid"
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

