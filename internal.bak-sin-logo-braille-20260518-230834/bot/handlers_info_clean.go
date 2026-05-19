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

PORTS=$(ss -tuln | awk 'NR>1 {print $5}' | awk -F: '{print $NF}' | sort -n | uniq | tr '\n' ' ')

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
echo "🛰️ PROTOCOLOS / PUERTOS ACTIVOS"

hay=0

# SSH
if echo "$PORTS" | grep -qw "22"; then
  echo "🔐 SSH: activo (22)"
  hay=1
fi

# DNS / SlowDNS
if echo "$PORTS" | grep -qw "53"; then
  echo "🐢 DNS / SlowDNS: activo (53)"
  hay=1
fi

# WebSocket 200
if systemctl is-active --quiet darkzsaid-ws200 2>/dev/null && echo "$PORTS" | grep -qw "80"; then
  echo "🚀 WebSocket 200 Established: activo (80)"
  hay=1
elif echo "$PORTS" | grep -qw "80"; then
  echo "🌐 Puerto 80: activo"
  hay=1
fi

# SSL 443
if systemctl is-active --quiet stunnel4 2>/dev/null && echo "$PORTS" | grep -qw "443"; then
  echo "📜 SSL/Stunnel: activo (443)"
  hay=1
elif echo "$PORTS" | grep -qw "443"; then
  echo "🔒 Puerto 443: activo"
  hay=1
fi

# Puertos VPN personalizados que usás en el panel
if echo "$PORTS" | grep -qw "5667"; then
  echo "🛰️ ZiVPN: activo (5667)"
  hay=1
fi

if echo "$PORTS" | grep -qw "7100"; then
  echo "📡 VPN: activo (7100)"
  hay=1
fi

if echo "$PORTS" | grep -qw "7200"; then
  echo "📡 VPN: activo (7200)"
  hay=1
fi

# BadVPN
if echo "$PORTS" | grep -qw "7300"; then
  echo "🎮 BadVPN UDPGW: activo (7300)"
  hay=1
fi

# Mostrar cualquier otro puerto activo que no esté identificado arriba
for p in $PORTS; do
  case "$p" in
    22|53|80|443|5667|7100|7200|7300)
      ;;
    *)
      echo "🔹 Puerto activo: $p"
      hay=1
      ;;
  esac
done

if [ "$hay" = "0" ]; then
  echo "Sin protocolos activos detectados."
fi

echo "━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "ℹ️ Extrainfo:"
echo "Puertos activos: $PORTS"
`

out, err := exec.Command("bash", "-lc", cmd).CombinedOutput()
if err != nil {
return SafeEditCtx(c, b, "❌ Error obteniendo información del servidor:\n\n"+string(out), markup)
}

return SafeEditCtx(c, b, string(out), markup)
}
