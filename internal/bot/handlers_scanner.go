package bot

import (
	"fmt"
	"os"
	"strings"

	"github.com/DarkZsaidscript/BOT-TELEGRAM-VPN/internal/sys"
	tele "gopkg.in/telebot.v3"
)

func handleMenuScanner(c tele.Context, b *tele.Bot) error {
	assetOK, httpxOK := sys.GetScannerStatus()

	markup := &tele.ReplyMarkup{}
	btnBack := markup.Data("🔙 Volver", "back_main")

	if !assetOK || !httpxOK {
		// Herramientas no instaladas
		markup.Inline(markup.Row(btnBack))
		texto := "🔍 <b>DarkZsaid Scanner 🌐</b>\n\n"
		texto += "⚠️ <b>Herramientas no instaladas</b>\n\n"
		if !assetOK {
			texto += "❌ <b>assetfinder:</b> No instalado\n"
		} else {
			texto += "✅ <b>assetfinder:</b> Instalado\n"
		}
		if !httpxOK {
			texto += "❌ <b>httpx:</b> No instalado\n"
		} else {
			texto += "✅ <b>httpx:</b> Instalado\n"
		}
		texto += "\n<i>Instala las herramientas desde</i> ⚙️ <b>Protocolos</b> → 🔍 <b>Escaner</b>"
		return SafeEditCtx(c, b, texto, markup)
	}

	btnStart := markup.Data("🔍 Iniciar Escaneo", "start_scanner_prompt")
	markup.Inline(markup.Row(btnStart), markup.Row(btnBack))

	texto := "🔍 <b>DarkZsaid Scanner 🌐</b>\n\n"
	texto += "✅ <b>Herramientas listas</b>\n"
	texto += "• assetfinder: ✅\n"
	texto += "• httpx: ✅\n\n"
	texto += "🚀 <b>Funciones:</b>\n"
	texto += "- Enumeración pasiva (Assetfinder)\n"
	texto += "- Detección de Tech/CDN (httpx)\n"
	texto += "- Reporte en tiempo real\n\n"
	texto += "<i>Disponible para todos los usuarios.</i>"

	return SafeEditCtx(c, b, texto, markup)
}

func handleStartScanPrompt(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	SetUserStep(chatID, "awaiting_scanner_domain")
	SetLastBotMsg(chatID, c.Message())

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("❌ Cancelar", "menu_scanner")))

	return c.Edit("🌐 <b>Escaner:</b>\n\n✏️ <i>Escribe el dominio que deseas escanear (ej: google.com):</i>", markup, tele.ModeHTML)
}

func processScannerSteps(step string, text string, chatID int64, c tele.Context, b *tele.Bot, lastMsg *tele.Message) error {
	if step != "awaiting_scanner_domain" {
		return nil
	}

	domain := strings.TrimSpace(text)
	DeleteUserStep(chatID)

	// Verificar herramientas antes de escanear
	assetOK, httpxOK := sys.GetScannerStatus()
	if !assetOK || !httpxOK {
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Volver", "menu_scanner")))
		b.Edit(lastMsg, "❌ <b>Herramientas no instaladas.</b>\n\nInstálalas desde ⚙️ <b>Protocolos</b> → 🔍 <b>Escaner</b>", markup, tele.ModeHTML)
		return nil
	}

	b.Edit(lastMsg, fmt.Sprintf("⏳ <b>Escaneando:</b> <code>%s</code>\n\n<i>Esto puede tardar unos segundos, por favor espera...</i>", domain), tele.ModeHTML)

	go func() {
		result, err := sys.RunScanner(domain)
		markup := &tele.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("🔙 Volver al Escaner", "menu_scanner")))

		if err != nil {
			b.Edit(lastMsg, fmt.Sprintf("❌ <b>Error en el Escaneo:</b>\n<pre>%v</pre>", err), markup, tele.ModeHTML)
			return
		}

		// Limitar el resultado si es muy largo para Telegram (4096 chars)
		if len(result) > 3500 {
			// Enviar primero un adelanto
			preview := result[:3000] + "\n\n... (Reporte completo en el archivo adjunto)"
			header := fmt.Sprintf("✅ <b>Resultados (Previsualización):</b> <code>%s</code>\n", domain)
			header += "━━━━━━━━━━━━━━\n"
			b.Edit(lastMsg, header+preview, markup, tele.ModeHTML)

			// Crear archivo temporal para el reporte completo
			tmpFile := fmt.Sprintf("/tmp/scan_%s.txt", domain)
			_ = os.WriteFile(tmpFile, []byte(result), 0644)
			defer os.Remove(tmpFile)

			doc := &tele.Document{
				File:     tele.FromDisk(tmpFile),
				FileName: fmt.Sprintf("reporte_escaneo_%s.txt", domain),
				Caption:  fmt.Sprintf("📄 Reporte completo de escaneo para %s", domain),
			}
			b.Send(c.Sender(), doc)
			return
		}

		header := fmt.Sprintf("✅ <b>Resultados para:</b> <code>%s</code>\n", domain)
		header += "━━━━━━━━━━━━━━\n"
		b.Edit(lastMsg, header+result, markup, tele.ModeHTML)
	}()

	return nil
}

// handleSubMenuScanner muestra el submenú de gestión del escáner en Protocolos
func handleSubMenuScanner(c tele.Context, b *tele.Bot) error {
	assetOK, httpxOK := sys.GetScannerStatus()

	installed := 0
	total := 2
	if assetOK {
		installed++
	}
	if httpxOK {
		installed++
	}

	pct := (installed * 100) / total

	statusAsset := "❌ No instalado"
	if assetOK {
		statusAsset = "✅ Instalado"
	}
	statusHttpx := "❌ No instalado"
	if httpxOK {
		statusHttpx = "✅ Instalado"
	}

	globalStatus := "❌ No instalado"
	if installed == total {
		globalStatus = "✅ Completo"
	} else if installed > 0 {
		globalStatus = "⚠️ Parcial"
	}

	markup := &tele.ReplyMarkup{}
	btnInstall := markup.Data("📥 Instalar Todo", "install_scanner_all")
	btnUninstall := markup.Data("🗑️ Desinstalar Todo", "uninstall_scanner_all")
	btnBack := markup.Data("🔙 Volver", "menu_protocols")

	markup.Inline(
		markup.Row(btnInstall),
		markup.Row(btnUninstall),
		markup.Row(btnBack),
	)

	// Barra de progreso visual
	barFull := installed
	barEmpty := total - installed
	bar := strings.Repeat("█", barFull*5) + strings.Repeat("░", barEmpty*5)

	texto := "🔍 <b>Gestión de Herramientas de Escaneo</b>\n\n"
	texto += fmt.Sprintf("📊 <b>Estado:</b> %s\n", globalStatus)
	texto += fmt.Sprintf("📈 <b>Progreso:</b> [%s] <code>%d%%</code> (%d/%d)\n\n", bar, pct, installed, total)
	texto += fmt.Sprintf("🔧 <b>assetfinder:</b> %s\n", statusAsset)
	texto += fmt.Sprintf("🔧 <b>httpx:</b> %s\n\n", statusHttpx)
	texto += "<i>Estas herramientas son necesarias para el módulo 🔍 Escaner del menú principal.</i>\n\n"
	texto += "⚠️ <i>La instalación puede tardar 1-3 minutos por herramienta según los recursos del servidor.</i>"

	return SafeEditCtx(c, b, texto, markup)
}

// handleInstallScannerAll instala todas las herramientas del escáner con progreso
func handleInstallScannerAll(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	if !isSuperAdminID(chatID) {
		return c.Respond(&tele.CallbackResponse{Text: "⛔ Solo el SuperAdmin puede instalar herramientas.", ShowAlert: true})
	}

	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "submenu_scanner")))

	go func() {
		lastMsg := c.Message()
		assetOK, httpxOK := sys.GetScannerStatus()
		totalSteps := 0
		if !assetOK {
			totalSteps++
		}
		if !httpxOK {
			totalSteps++
		}

		if totalSteps == 0 {
			b.Edit(lastMsg, "✅ <b>Todas las herramientas ya están instaladas.</b>", markup, tele.ModeHTML)
			return
		}

		currentStep := 0

		// Instalar assetfinder
		if !assetOK {
			currentStep++
			pct := (currentStep * 100) / (totalSteps + 1)
			bar := strings.Repeat("█", pct/10) + strings.Repeat("░", 10-pct/10)
			b.Edit(lastMsg, fmt.Sprintf("⏳ <b>Instalando herramientas de Escaner...</b>\n\n📈 [%s] <code>%d%%</code>\n\n🔧 Instalando <b>assetfinder</b>...\n<i>Esto puede tardar 1-2 minutos...</i>", bar, pct), tele.ModeHTML)

			if err := sys.InstallScannerTool("assetfinder"); err != nil {
				b.Edit(lastMsg, fmt.Sprintf("❌ <b>Error instalando assetfinder:</b>\n<pre>%v</pre>", err), markup, tele.ModeHTML)
				return
			}
		}

		// Instalar httpx
		if !httpxOK {
			currentStep++
			pct := (currentStep * 100) / (totalSteps + 1)
			bar := strings.Repeat("█", pct/10) + strings.Repeat("░", 10-pct/10)
			b.Edit(lastMsg, fmt.Sprintf("⏳ <b>Instalando herramientas de Escaner...</b>\n\n📈 [%s] <code>%d%%</code>\n\n✅ assetfinder instalado\n🔧 Instalando <b>httpx</b>...\n<i>Esto puede tardar 2-3 minutos...</i>", bar, pct), tele.ModeHTML)

			if err := sys.InstallScannerTool("httpx"); err != nil {
				b.Edit(lastMsg, fmt.Sprintf("❌ <b>Error instalando httpx:</b>\n<pre>%v</pre>", err), markup, tele.ModeHTML)
				return
			}
		}

		// Finalizado
		b.Edit(lastMsg, "✅ <b>Herramientas de Escaner Instaladas Correctamente</b>\n\n📈 [██████████] <code>100%</code>\n\n✅ assetfinder: Instalado\n✅ httpx: Instalado\n\n<i>Ya puedes usar el módulo 🔍 Escaner desde el menú principal.</i>", markup, tele.ModeHTML)
	}()

	return nil
}

// handleUninstallScannerAll desinstala todas las herramientas
func handleUninstallScannerAll(c tele.Context, b *tele.Bot) error {
	chatID := c.Chat().ID
	if !isSuperAdminID(chatID) {
		return c.Respond(&tele.CallbackResponse{Text: "⛔ Solo el SuperAdmin puede desinstalar herramientas.", ShowAlert: true})
	}

	SafeEditCtx(c, b, "⏳ <i>Desinstalando herramientas de Escaner...</i>", nil)

	err := sys.UninstallAllScannerTools()
	markup := &tele.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("🔙 Volver", "submenu_scanner")))

	if err != nil {
		return SafeEditCtx(c, b, fmt.Sprintf("⚠️ <b>Error parcial:</b>\n%v", err), markup)
	}
	return SafeEditCtx(c, b, "✅ <b>Herramientas de Escaner desinstaladas.</b>\n\n<i>assetfinder y httpx fueron eliminados del sistema.</i>", markup)
}
