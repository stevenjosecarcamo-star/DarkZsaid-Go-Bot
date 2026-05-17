# 💎 DARKZSAID BOT GO EDITION

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/OS-Ubuntu%2024.04-E95420?style=for-the-badge&logo=ubuntu" alt="Ubuntu 24.04">
  <img src="https://img.shields.io/badge/Platform-Linux-FCC624?style=for-the-badge&logo=linux" alt="Linux">
  <img src="https://img.shields.io/badge/Status-Stable-success?style=for-the-badge" alt="Status">
  <img src="https://img.shields.io/badge/Version-7.7-blue?style=for-the-badge" alt="Version">
</p>

---

## 🚀 ¿Qué es DarkZsaid Bot?

**DarkZsaid Bot Go Edition** es una solución integral y de alto rendimiento para la gestión de servidores VPN y cuentas SSH a través de Telegram. Reescrito completamente en **Go** para garantizar velocidad, estabilidad y bajo consumo de recursos, este bot transforma tu VPS en un panel de control profesional y automatizado.

---

## 🚀 Novedades v7.7 — Modo Desarrollador y Forzado de Actualizaciones
- Se añadió un botón "⚠️ Forzar Reinstalación (Dev)" que permite reinstalar completamente el repositorio en cualquier momento, ignorando verificaciones de versión.
- El motor interno del actualizador automático (`systemd-run`) fue reescrito para escapar limpiamente de los cgroups, garantizando un índice de éxito del 100% al auto-compilarse.

---

## 🚀 Novedades v7.5 — Payloads Avanzados e Interfaz Mejorada
- Se ha integrado soporte extendido para **múltiples tipos de payloads Cloudflare** (WS Puerto 80, WSS/TLS Puerto 443, Injector SNI).
- El mensaje de creación de cuenta SSH ahora genera automáticamente las líneas de configuración para **HTTP Custom**, adaptándose inteligentemente a los puertos activos (Dropbear, SSL Tunnel, WS) listas para tocar y copiar.

---

## 🔥 Novedades v7.4 — Ultra Seguridad y Optimización de CPU

### 🛡️ Seguridad SSH Extrema (Anti-Abuso)
- **Bloqueo de Consola:** Las nuevas cuentas SSH se crean con consola `/bin/false`. Esto significa que **solo sirven como VPN** (túnel HTTP Custom/Injector, SlowDNS, WebSocket), pero los usuarios ya no pueden entrar a la terminal para ejecutar comandos.
- **Bloqueo de Archivos:** Ya no se crea la carpeta `/home` para los clientes. Es imposible que usen tu servidor para subir, descargar o alojar archivos basura.
- **Eliminación y Desconexión Agresiva:** Al eliminar un usuario (manual o por expiración), el bot lanza `killall -9` y `pkill` para **desconectar forzosamente** todas sus conexiones activas al instante, y ejecuta un borrado profundo (`rm -rf`) en el SSD para destruir cualquier rastro del usuario.

### ⚡ Optimización Masiva de Rendimiento
- **Reducción del 90% en uso de CPU:** Se reescribió por completo el motor de monitoreo de límites SSH. Antes, el bot consultaba a nivel de sistema (`grep`) y recargaba el servicio SSH cada 60 segundos por cada usuario. Ahora la lectura es directa en la RAM de Go y `sshd` solo se recarga estrictamente cuando se crea o elimina una cuenta. Ideal para servidores con cientos de usuarios conectados simultáneamente sin lag ni congelamientos.

---

## 🆕 Novedades v7.3 — Resiliencia HAProxy + Fix SSL/TLS Directo

### 🔒 HAProxy Auto-Recovery tras Reboot

HAProxy ahora **sobrevive cualquier reinicio del VPS** de forma automática gracias a:

| Mecanismo | Descripción |
|-----------|-------------|
| **Systemd Override** | `Restart=always` + `StartLimitIntervalSec=0` para reinicio ilimitado |
| **Network Dependency** | `After=network-online.target` — espera a que la red esté lista |
| **Socket Directory** | `ExecStartPre=/bin/mkdir -p /run/haproxy` — recrea el directorio tmpfs perdido en reboot |
| **Bot Startup Check** | Al iniciar, el bot verifica HAProxy, mata procesos invasores en 80/443/8080 y reinicia si es necesario |
| **Backend Recovery** | También verifica y reinicia `ssh-ws-internal` (puerto 10015) automáticamente |

### 📱 Fix SSL/TLS Directo (HTTP Injector)

Corregido el modo **SSL/TLS → SSH** en HTTP Injector. El tráfico SSH directo (sin WebSocket) ahora se enruta correctamente a OpenSSH puerto 22 en lugar del antiguo puerto 2222 (Dropbear).

**Flujo corregido:**
```
HTTP Injector (SSL/TLS) → HAProxy(443 TLS) → detecta SSH-2.0 → OpenSSH(22) ✅
```

---

## 📋 Novedades v7.2 — Banner SSH Dinámico por Usuario

Cada cuenta SSH ahora genera automáticamente un **banner HTML personalizado** que se muestra al conectarse, compatible con **HTTP Injector**, **HTTP Custom**, **HA Tunnel** y todas las apps VPN.

### ¿Qué incluye el banner?

| Elemento | Descripción |
|----------|-------------|
| 🎨 **Logo DarkZsaid** | Logo animado en arte braille (idéntico al banner global) |
| 🏷️ **Título Personalizado** | El admin elige el título al crear la cuenta (ej: `INTERNET ILIMITADO`, `SPEED PREMIUM VIP`) |
| 👤 **Datos de la Cuenta** | Usuario, fecha de vencimiento, días restantes, límite de dispositivos |
| 📢 **Promoción** | Canal @DarkZsaid2 y soporte @DarkZsaid |
| ⚠️ **Reglas** | Normas del servidor con advertencia de ban automático |

### Ejemplo del Banner (vista en HTTP Injector/Custom)

<p align="center">
  <img width="350" alt="Banner SSH DarkZsaid" src="https://img.shields.io/badge/Formato-HTML%20VPN%20Apps-29b6f6?style=for-the-badge">
</p>

```text
══════════════════════
⠀⠀⢀⣶⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣶⡀⠀⠀
⠀⠀⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀
⠀⠀⢸⣿⡇⠀⠀⠀⣠⣶⣄⠀⠀⠀⢸⣿⡇⠀⠀
⠀⠀⢸⣿⡇⠀⠀⢰⣿⣿⣿⡆⠀⠀⢸⣿⡇⠀⠀
⠀⠀⠈⣿⣿⡄⢀⣿⣿⠻⣿⣿⡀⢠⣿⣿⠁⠀⠀
⠀⠀⠀⠹⣿⣿⣾⣿⡏⠀⢹⣿⣷⣿⣿⠏⠀⠀⠀
⠀⠀⠀⠀⠙⢿⣿⡿⠀⠀⠀⢿⣿⡿⠋⠀⠀⠀⠀
        DARKZSAID       
══════════════════════
 ⚡ INTERNET ILIMITADO ⚡
══════════════════════
 👤 Usuario: pepito
 📅 Vence: 2026-05-25
 ⏳ Días Restant.: 30
 💻 Límite: 3
══════════════════════
 🔥 ¡SERVIDORES PREMIUM A 8.5 SOLES! 🔥
 📢 Canal: @DarkZsaid2
 👤 Soporte: @DarkZsaid
══════════════════════
 ✅ CREADO EN : @DarkZsaid_bot
══════════════════════
```

> [!NOTE]
> El banner real usa **formato HTML con colores** (verde, cyan, amarillo, magenta). La vista anterior es una representación simplificada. En las apps VPN se ve con los colores del banner predeterminado de DarkZsaid.

### ¿Cómo funciona?

### ¿Cómo funciona?

- **Banners Individuales (Por Defecto):**
  - Al crear una cuenta SSH, el bot pide un **título personalizado** para ese cliente.
  - Se genera un archivo HTML en `/etc/ssh_banners/{usuario}.banner` con sus días restantes.
  - Los **textos promocionales** (Canal, Soporte y Mensaje de Venta) se pueden editar directamente desde el menú del bot en *Ajustes Pro -> Banner -> Editar Textos Promo*.
  - Los **días restantes se actualizan automáticamente** cada 60 segundos.
  - Al renovar o eliminar un usuario, el banner se regenera o limpia automáticamente.
- **Banner Global:**
  - Si prefieres no pedir títulos personalizados para cada usuario, puedes activar el **Banner Global** desde *Ajustes Pro -> Banner*.
  - Al hacer esto, **el bot omitirá la pregunta del título al crear cuentas SSH** y todos los usuarios verán el mismo banner general.

---

## ✨ Características Principales

### 🛠️ Gestión de Protocolos (All-in-One)
- **SSH/Dropbear/SSL Tunnel:** Gestión completa de cuentas con límites de conexión y banners individuales.
- **SlowDNS:** Instalación y configuración de túneles DNS.
- **Xray (VMess):** Protocolo de última generación sobre WebSocket compatible con Cloudflare y HAProxy.
- **ZiVPN & UDP Custom:** Soporte para protocolos de gaming y bypass robusto.
- **ProxyDT:** Integración con ProxyDT Cracked para túneles HTTP estables.
- **Falcon Proxy:** Proxy de alto rendimiento integrado.
- **SSH WebSocket:** Proxy WS/WSS en puertos 10015 y 2082.

### 🔍 Escáner de Red (Módulo Independiente)
- Submenú dedicado en **Protocolos → 🔍 Escaner** con estado de instalación y barra de progreso.
- Instalación/desinstalación on-demand de `assetfinder` y `httpx`.
- No bloquea la instalación del bot — se instala cuando el admin lo necesite.

### 🛡️ Administración Pro y Utilidades
- **Ajustes Pro y Auto-Reboot:** Panel interno con control de estado, configuración de accesos públicos/privados y **reinicio automático diario**.
- **Mensaje Global (Broadcast):** Envía anuncios masivos a todos los usuarios.
- **Monitoreo en Tiempo Real:** Visualización en vivo de métricas de VPS (Uptime, Consumo) y conexiones activas (SSH, ZiVPN, Xray).
- **Gestión Avanzada:** Soporte nativo para dominios Cloudflare (CF) y Cloudfront con autogeneración de payloads.
- **Cuotas de Creación:** Límites configurables de días y dispositivos para usuarios públicos y admins.

### 🧹 Mantenimiento Inteligente
- **Persistencia de Datos Inquebrantable:** Tu tráfico de red y configuraciones de usuario están a salvo ante cualquier reinicio.
- **Resiliencia de Servicios:** Recuperación automática de HAProxy, Xray y ssh-ws-internal mediante systemd overrides (`Restart=always`, `network-online.target`).
- **HAProxy Auto-Recovery:** Al iniciar el bot, verifica que HAProxy esté corriendo, recrea `/run/haproxy`, mata procesos invasores en puertos 80/443/8080 y reinicia si es necesario.
- **Deep System Cleanup:** Liberación automática de memoria y procesos huérfanos.

---

## 📥 Instalación Rápida (Universal)

> [!NOTE]
> **Compatibilidad OS:** Este bot fue desarrollado y probado rigurosamente en **Ubuntu 24.04**. Se recomienda encarecidamente utilizar esta versión (o distribuciones basadas en ella) para garantizar el correcto funcionamiento de todas las dependencias (Go, Systemd, SSH, Xray, SlowDNS, etc).

Ejecuta el siguiente comando en tu terminal como usuario **root**:

```bash
bash <(curl -sL https://raw.githubusercontent.com/DarkZsaidscript/BOT-TELEGRAM-VPN/main/install_go.sh)
```

> [!IMPORTANT]
> Selecciona la opción **1** del menú. El instalador configurará automáticamente Go, compilará el bot y desplegará el servicio Systemd 24/7.

---

## 🔄 Cómo Actualizar

Si ya tienes el bot funcionando y quieres recibir parches y nuevas funciones **sin perder usuarios ni configuraciones**:

```bash
bash <(curl -sL https://raw.githubusercontent.com/DarkZsaidscript/BOT-TELEGRAM-VPN/main/install_go.sh)
```

> [!NOTE]
> Selecciona la opción **1 (Instalar / Actualizar Bot)**. El sistema detectará tus credenciales existentes y solo actualizará el código.

### Método alternativo (manual):

```bash
systemctl stop darkzsaid
cd /tmp && rm -rf BOT-TELEGRAM-VPN
git clone https://github.com/DarkZsaidscript/BOT-TELEGRAM-VPN.git
cd BOT-TELEGRAM-VPN
export PATH=$PATH:/usr/local/go/bin
go build -o /usr/local/bin/darkzsaid-bot cmd/darkzsaid/main.go
systemctl restart darkzsaid
```

---

## ☁️ Copias de Seguridad en Google Drive

El bot incluye un sistema de respaldos integrado con tu Google Drive personal. Permite copias **manuales** desde el panel y **automáticas cada 24 horas**. Solo necesitas configurarlo una vez.

### Requisitos Previos

- Una cuenta de Google (Gmail)
- Acceso a [Google Cloud Console](https://console.cloud.google.com/)

### Paso 1: Crear Proyecto en Google Cloud

1. Ve a [Google Cloud Console](https://console.cloud.google.com/) e inicia sesión con tu Gmail.
2. Crea un **nuevo proyecto** (o usa uno existente). Ponle un nombre como `DarkZsaidBackup`.
3. En el menú lateral, ve a **APIs y Servicios** → **Biblioteca**.
4. Busca **"Google Drive API"** y haz clic en **Habilitar**.

### Paso 2: Configurar Pantalla de Consentimiento OAuth

1. Ve a **APIs y Servicios** → **Pantalla de consentimiento de OAuth**.
2. Selecciona **Externo** como tipo de usuario y haz clic en **Crear**.
3. Rellena solo los campos obligatorios:
   - **Nombre de la aplicación:** `DarkZsaidBot`
   - **Correo de asistencia:** tu Gmail
   - **Correo del desarrollador:** tu Gmail
4. Haz clic en **Guardar y Continuar** hasta llegar a **Usuarios de Prueba**.
5. **⚠️ IMPORTANTE:** Haz clic en **+ Agregar Usuarios** y escribe **tu propio correo de Gmail**. Sin esto, no podrás autorizar el bot.
6. Haz clic en **Guardar y Continuar** → **Volver al Panel**.
7. **⚠️ MUY IMPORTANTE (Evitar error invalid_grant):** En la página principal de la Pantalla de Consentimiento, bajo **"Estado de la publicación"**, haz clic en el botón **"Publicar la aplicación" (Publish App)** o "Poner en producción" y confirma. Si dejas la aplicación "En Pruebas" (Testing), **el token caducará cada 7 días** automáticamente.

### Paso 3: Crear Credenciales OAuth

1. Ve a **APIs y Servicios** → **Credenciales**.
2. Haz clic en **+ Crear Credenciales** → **ID de cliente de OAuth**.
3. En **Tipo de aplicación**, selecciona **Aplicación de escritorio**.
4. Ponle un nombre (ej: `DarkZsaidBot`) y haz clic en **Crear**.
5. Se mostrará un diálogo con tu **Client ID** y **Client Secret**.
6. Haz clic en **⬇ Descargar JSON** para descargar el archivo de credenciales.

### Paso 4: Subir Credenciales al VPS

1. Renombra el archivo descargado a exactamente: `credentials.json`
2. Súbelo a tu VPS en la ruta: `/opt/darkzsaid_bot/credentials.json`

Puedes subirlo con `scp` desde tu PC:
```bash
scp credentials.json root@TU_IP_VPS:/opt/darkzsaid_bot/credentials.json
```

O crearlo directamente en el VPS:
```bash
nano /opt/darkzsaid_bot/credentials.json
# Pega el contenido del JSON y guarda (Ctrl+O, Enter, Ctrl+X)
```

### Paso 5: Vincular con Telegram

1. Abre Telegram y envía al bot: `/authdrive`
2. El bot te enviará un **enlace de Google**. Ábrelo en tu navegador.
3. Inicia sesión con el **mismo Gmail** que agregaste como usuario de prueba.
4. Google te pedirá permisos — haz clic en **Continuar** y **Permitir**.
5. Al final verás una página con un **código de autorización** (o si la página muestra error de `localhost`, copia el código largo que aparece en la barra de direcciones después de `code=` y antes de `&scope`).
6. Envía al bot: `/authdrive PEGA_TU_CODIGO_AQUI`

> [!TIP]
> El código es largo (puede tener más de 50 caracteres). Cópialo completo.

### ¡Listo!

- El bot vinculará tu Drive **de por vida** (el token se renueva automáticamente).
- Se creará la carpeta **`BotVPN_Backups`** en tu Drive.
- El bot subirá tu base de datos (`bot_data.json`) automáticamente cada 24 horas.
- Mantiene solo las **2 copias más recientes** para no ocupar espacio.
- Usa el botón **📥 Restaurar Backup** desde Ajustes Pro cuando instales el bot en un nuevo VPS para recuperar todos tus usuarios automáticamente.

> [!WARNING]
> Si el token se revoca o caduca verás el error `oauth2: "invalid_grant" "Token has been expired"`. Esto pasa si no pusiste la app "En Producción" (Paso 2.7) caducando a los 7 días. Para arreglarlo hoy, solo necesitas repetir el paso 5 enviando `/authdrive` al bot.

---

## 🛠️ Solución de Problemas (Troubleshooting)

| Síntoma | Causa Probable | Solución |
| :--- | :--- | :--- |
| **El bot no responde** | OOM mató el proceso o token inválido | `systemctl status darkzsaid` → `systemctl restart darkzsaid` |
| **Puerto 443 no responde tras reboot** | HAProxy caído o `/run/haproxy` borrado | Ver sección "Fix HAProxy" abajo |
| **HTTP Injector no conecta (SSL/TLS)** | Backend apuntaba a Dropbear (2222) | Actualizar bot — fix incluido en v7.3 |
| **Xray/VMess no conecta** | HAProxy o Xray no iniciaron | `systemctl status haproxy xray` |
| **VPS muy lenta** | RAM saturada | Activar **Auto Reboot** en Ajustes Pro |
| **Banner no aparece** | sshd no recargó | `systemctl reload ssh` o recrear la cuenta |
| **Error en Google Drive (`invalid_grant`)** | Token caducado a los 7 días (App "En Pruebas") | Enviar `/authdrive` al bot. Cambiar app a "Producción" en GCP. |
| **Escáner no funciona** | Herramientas no instaladas | Ir a **Protocolos → 🔍 Escaner → 📥 Instalar Todo** |

### 🔒 Fix HAProxy / Puerto 443 Caído

Si el puerto 443 no responde después de un reinicio del VPS:

```bash
# 1. Diagnosticar estado
systemctl status haproxy --no-pager
ss -tlnp | grep -E ':443|:80|:8080'
systemctl status ssh-ws-internal --no-pager

# 2. Fix rápido manual
mkdir -p /run/haproxy
fuser -k 443/tcp 2>/dev/null
systemctl restart ssh-ws-internal
systemctl restart haproxy

# 3. Verificar que todo funciona
ss -tlnp | grep haproxy
echo | openssl s_client -connect 127.0.0.1:443 2>&1 | grep CONNECTED
```

> [!NOTE]
> En la versión 7.3+, el bot aplica automáticamente un **override de systemd** que hace que HAProxy se auto-reinicie tras cualquier reboot. Si acabas de actualizar, el fix se aplica al primer arranque del bot.

### 📱 Verificar Conexión SSL para HTTP Injector

```bash
# Test completo de la cadena SSL → WebSocket → SSH
# 1. SSL handshake
echo | openssl s_client -connect 127.0.0.1:443 -servername localhost 2>&1 | grep -E "CONNECTED|subject"

# 2. WebSocket upgrade (simula HTTP Injector)
echo -e "GET / HTTP/1.1\r\nHost: localhost\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n" | timeout 3 nc 127.0.0.1 10015

# Resultado esperado:
# HTTP/1.1 101 Switching Protocols
# SSH-2.0-OpenSSH_9.6p1
```

### Comandos Útiles de Diagnóstico

```bash
# Estado del bot
systemctl status darkzsaid

# Logs en tiempo real
journalctl -u darkzsaid -f --no-pager -n 50

# Reiniciar bot
systemctl restart darkzsaid

# Estado completo de servicios VPN
for svc in darkzsaid haproxy ssh-ws-internal xray ssh-ws; do
  echo -n "$svc: "; systemctl is-active $svc 2>/dev/null || echo "not found"
done

# Puertos en uso
ss -tlnp | grep -E ':22|:80|:443|:8080|:10015'

# Verificar banners de usuarios
ls -la /etc/ssh_banners/

# Ver configuración SSH (Match User blocks)
grep -A1 "Match User" /etc/ssh/sshd_config

# Logs de HAProxy
journalctl -u haproxy --no-pager -n 20
```

---

## 💎 Créditos y Soporte

Este proyecto es desarrollado y mantenido con pasión por:

- **👨‍💻 Desarrollador:** [@DarkZsaid](https://t.me/Dan3651)
- **📢 Canal Oficial:** [DarkZsaid Channel](https://t.me/DarkZsaid2)

---

<p align="center">
  <i>"Potenciando tu VPS con la velocidad de Go."</i><br>
  <b>© 2026 DarkZsaid Project</b>
</p>
