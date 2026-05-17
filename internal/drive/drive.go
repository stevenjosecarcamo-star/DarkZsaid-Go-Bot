package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	CredentialsFile = "/opt/darkzsaid_bot/credentials.json"
	TokenFile       = "/opt/darkzsaid_bot/token.json"
	FolderName      = "BotVPN_Backups"
	MaxBackups      = 2
)

func getClient(ctx context.Context, config *oauth2.Config, tok *oauth2.Token) *http.Client {
	return config.Client(ctx, tok)
}

func getConfig() (*oauth2.Config, error) {
	b, err := os.ReadFile(CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("imposible leer credentials.json: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("imposible parsear la configuracion: %v", err)
	}
	return config, nil
}

func GetAuthURL() (string, error) {
	config, err := getConfig()
	if err != nil {
		return "", err
	}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

func SaveToken(authCode string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return fmt.Errorf("no se pudo recuperar el token: %v", err)
	}
	f, err := os.OpenFile(TokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("no se pudo crear token.json: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

func getService() (*drive.Service, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	
	f, err := os.Open(TokenFile)
	if err != nil {
		return nil, fmt.Errorf("no estas logueado")
	}
	defer f.Close()
	
	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	
	client := getClient(context.TODO(), config, tok)
	srv, err := drive.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func IsAuthenticated() bool {
	_, err := getService()
	return err == nil
}

func UploadBackup(filePath string) error {
	srv, err := getService()
	if err != nil {
		return err
	}
	
	folderId, err := ensureFolder(srv)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("bot_data_%s.json", time.Now().Format("20060102_150405"))
	
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("imposible leer archivo a respaldar: %v", err)
	}
	defer f.Close()

	driveFile := &drive.File{
		Name:     fileName,
		Parents:  []string{folderId},
	}
	
	_, err = srv.Files.Create(driveFile).Media(f).Do()
	if err != nil {
		return fmt.Errorf("error subiendo archivo a drive: %v", err)
	}

	cleanupOldBackups(srv, folderId)
	return nil
}

func RestoreBackup(filePath string) error {
	srv, err := getService()
	if err != nil {
		return err
	}
	folderId, err := ensureFolder(srv)
	if err != nil {
		return err
	}
	
	q := fmt.Sprintf("'%s' in parents and trashed=false", folderId)
	r, err := srv.Files.List().Q(q).OrderBy("createdTime desc").Fields("files(id, name, createdTime)").Do()
	if err != nil {
		return err
	}
	
	if len(r.Files) == 0 {
		return fmt.Errorf("no se encontraron copias de seguridad en Drive")
	}
	
	latest := r.Files[0]
	resp, err := srv.Files.Get(latest.Id).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, filePath)
}

func ensureFolder(srv *drive.Service) (string, error) {
	q := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", FolderName)
	r, err := srv.Files.List().Q(q).Fields("files(id, name)").Do()
	if err != nil {
		return "", err
	}
	if len(r.Files) > 0 {
		return r.Files[0].Id, nil
	}
	
	f := &drive.File{
		Name:     FolderName,
		MimeType: "application/vnd.google-apps.folder",
	}
	// Al tratarse de la cuenta de un humano normal, no hay bloqueo de quota y crea la carpeta
	folder, err := srv.Files.Create(f).Do()
	if err != nil {
		return "", err
	}
	return folder.Id, nil
}

func cleanupOldBackups(srv *drive.Service, folderId string) {
	q := fmt.Sprintf("'%s' in parents and trashed=false", folderId)
	r, err := srv.Files.List().Q(q).OrderBy("createdTime desc").Fields("files(id, name, createdTime)").Do()
	if err != nil || len(r.Files) <= MaxBackups {
		return
	}
	
	for i := MaxBackups; i < len(r.Files); i++ {
		srv.Files.Delete(r.Files[i].Id).Do()
	}
}
