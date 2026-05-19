package sys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ExecCmdRun es una función auxiliar para ejecutar comandos del sistema (bash) de manera limpia
func ExecCmdRun(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd error: %v, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

// CreateSSHUser crea un usuario en el sistema con expiración y contraseña.
func CreateSSHUser(username string, password string, days int) error {
	// 1. Calcular Fecha Vencimiento
	expireDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	// 2. Ejecutar useradd -M -s /bin/false -e "fecha" "usuario"
	// Usamos -M para NO crear carpeta home (ahorra espacio y evita subida de archivos)
	// Usamos /bin/false para denegar el acceso a la consola/comandos, pero permite túnel VPN
	_, err := ExecCmdRun("useradd", "-M", "-s", "/bin/false", "-e", expireDate, username)
	if err != nil {
		return fmt.Errorf("fallo al crear usuario: %v", err)
	}

	// 3. chpasswd
	// En Go podemos usar la entrada estándar del comando para chpasswd
	cmd := exec.Command("chpasswd")
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%s:%s", username, password))
	if err := cmd.Run(); err != nil {
		// Rollback (borramos usuario si chpasswd falla)
		_ = DeleteSSHUser(username)
		return fmt.Errorf("fallo al asignar contraseña: %v", err)
	}

	return nil
}

// DeleteSSHUser borra el usuario, home y reglas asociadas de iptables
func DeleteSSHUser(username string) error {
	// 1. Matar TODOS los procesos del usuario para desconexión forzada instantánea
	exec.Command("killall", "-9", "-u", username).Run()
	exec.Command("pkill", "-9", "-u", username).Run()

	// También limpiar usando GetUserProcesses por si acaso (para demonios sshd)
	pids, _ := GetUserProcesses(username)
	for _, pid := range pids {
		exec.Command("kill", "-9", pid).Run()
	}

	// 2. Limpiar Iptables de inmediato (Módulo Quotas robusto)
	CleanUserRules(username)

	// 3. Borrar usuario con sed y userdel (con timeout manual para no congelar)
	exec.Command("sed", "-i", fmt.Sprintf("/^%s hard maxlogins/d", username), "/etc/security/limits.conf").Run()

	cmd := exec.Command("userdel", "-f", "-r", username)
	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case <-time.After(10 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		fmt.Printf("Aviso: userdel para %s tardó demasiado y fue terminado.\n", username)
	case err := <-done:
		if err != nil {
			fmt.Printf("Aviso: error en userdel para %s: %v\n", username, err)
		}
	}

	// 4. Limpieza forzada de rastros en disco SSD
	exec.Command("rm", "-rf", fmt.Sprintf("/home/%s", username)).Run()
	exec.Command("rm", "-rf", fmt.Sprintf("/var/spool/mail/%s", username)).Run()

	// 5. Archivo limit y banner
	os.Remove(fmt.Sprintf("/etc/ssh_limits/%s.limit", username))
	RemoveUserBanner(username)

	return nil
}

// UpdateSSHUserPassword cambia la contraseña de un usuario SSH
func UpdateSSHUserPassword(username, newPassword string) error {
	cmd := exec.Command("chpasswd")
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%s:%s", username, newPassword))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fallo al actualizar contraseña: %v", err)
	}
	return nil
}

// RenewSSHUser renueva un usuario sumando dias a la fecha actual y lo desbloquea.
func RenewSSHUser(username string, days int) error {
	expireDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	// Cambiar expiracion
	_, err := ExecCmdRun("usermod", "-e", expireDate, username)
	if err != nil {
		return err
	}

	// Desbloquear por si estaba vencido
	ExecCmdRun("passwd", "-u", username)
	return nil
}

// SetSSHBanner configura el banner de bienvenida de SSH
func SetSSHBanner(text string) error {
	// Guardar en /etc/sshd_banner de forma segura usando Go nativo
	err := os.WriteFile("/etc/sshd_banner", []byte(text), 0644)
	if err != nil {
		return err
	}

	// 2. Asegurar que sshd_config tiene el banner activado
	_, _ = ExecCmdRun("sed", "-i", "/^Banner/d", "/etc/ssh/sshd_config")
	_, _ = ExecCmdRun("sh", "-c", "echo 'Banner /etc/sshd_banner' >> /etc/ssh/sshd_config")

	// 3. Reiniciar SSH
	ExecCmdRun("systemctl", "reload", "ssh")

	return nil
}
