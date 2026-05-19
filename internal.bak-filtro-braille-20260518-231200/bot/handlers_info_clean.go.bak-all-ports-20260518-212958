package bot

import (
"os/exec"

tele "gopkg.in/telebot.v3"
)

func handleMenuInfoClean(c tele.Context, b *tele.Bot) error {
markup := &tele.ReplyMarkup{}
markup.Inline(markup.Row(markup.Data("🔙 Volver", "back_main")))

cmd := `
IP=$(curl -4 -s https://api.ipify.org 2>/dev/null || hostname -I | awk '{print $1}')
CPU=$(awk -F: '/model name/ {print $2; exit}' /proc/cpuinfo | sed 's/^ //')
RAM=$(free -m | awk '/Mem:/ {print $3"MB / "$2"MB"}')
DISCO=$(df -h / | awk 'NR==2 {print $3" / "$2}')
UPTIME=$(uptime -p | sed 's/up //')
CPUUSO=$(top -bn1 | awk -F'id,' '/Cpu/ {split($1,a,","); print 100-a[length(a)] "%"}' 2>/dev/null)

echo "🌐 INFORMACIÓN DEL SERVIDOR"
echo "━━━━━━━━━━━━━━━━━━━━"
echo "🌎 IP: $IP"
echo "🖥️ CPU: $CPU"
echo "🔥 Uso: $CPUUSO"
echo "💾 RAM: $RAM"
echo "💿 Disco: $DISCO"
echo "⏱️ Uptime: $UPTIME"
echo "━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🛰️ PROTOCOLOS ACTIVOS"

hay=0

if systemctl is-active --quiet udpcustom 2>/dev/null || systemctl is-active --quiet udp-custom 2>/dev/null; then
  echo "🛰️ UDP Custom: activo"
  hay=1
fi

if systemctl is-active --quiet badvpn-udpgw 2>/dev/null || ss -ulnp 2>/dev/null | grep -q ':7300'; then
  echo "🎮 BadVPN UDPGW: activo (7300)"
  hay=1
fi

if systemctl is-active --quiet darkzsaid-ws200 2>/dev/null && ss -tulnp 2>/dev/null | grep -q ':80'; then
  echo "🚀 WebSocket 200 Established: activo (80)"
  hay=1
fi

if systemctl is-active --quiet dropbear 2>/dev/null; then
  echo "🐻 Dropbear: activo"
  hay=1
fi

if systemctl is-active --quiet ssh 2>/dev/null || systemctl is-active --quiet sshd 2>/dev/null; then
  echo "🔐 SSH: activo (22)"
  hay=1
fi

if systemctl is-active --quiet darkzsaid-stunnel 2>/dev/null || systemctl is-active --quiet stunnel4 2>/dev/null; then
  echo "📜 SSL/Stunnel: activo"
  hay=1
fi

if systemctl is-active --quiet x-ui 2>/dev/null; then
  echo "💎 Xray/3X-UI: activo"
  hay=1
fi

if [ "$hay" = "0" ]; then
  echo "Sin protocolos activos detectados."
fi

echo "━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "ℹ️ Extrainfo:"
echo -n "Puertos activos: "
ss -tuln | awk 'NR>1 {print $5}' | awk -F: '{print $NF}' | sort -n | uniq | tr '\n' ' '
echo ""
`

out, err := exec.Command("bash", "-lc", cmd).CombinedOutput()
if err != nil {
return SafeEditCtx(c, b, "❌ Error obteniendo información del servidor:\n\n"+string(out), markup)
}

return SafeEditCtx(c, b, string(out), markup)
}
