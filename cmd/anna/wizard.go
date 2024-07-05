package anna

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anna-ssg/anna/v3/pkg/parser"
)

type WizardServer struct {
	server   *http.Server
	serveMux *http.ServeMux

	// Common logger for all parser functions
	InfoLogger *log.Logger
	// Common logger for all parser functions
	ErrorLogger *log.Logger
}

var FormSubmittedCh = make(chan struct{})

func NewWizardServer(addr string) *WizardServer {
	serveMuxLocal := http.NewServeMux()

	wizardServer := WizardServer{
		serveMux:    serveMuxLocal,
		server:      &http.Server{Addr: addr, Handler: serveMuxLocal},
		InfoLogger:  log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime),
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return &wizardServer
}

func (ws *WizardServer) Start() {
	ws.serveMux.HandleFunc("/submit", ws.handleSubmit)
	fs := http.FileServer(http.Dir("./site/static/wizard"))
	ws.serveMux.Handle("/", fs)
	ws.InfoLogger.Printf("Wizard is running at: http://localhost%s\n", ws.server.Addr)

	if err := ws.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		ws.ErrorLogger.Fatalf("Could not start server: %v", err)
	}
}

func (ws *WizardServer) Stop() error {
	return ws.server.Shutdown(context.Background())
}

func (ws *WizardServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ws.ErrorLogger.Println("Method not allowed")
		return
	}
	// got the form data, now ask theme.go to unzip and place in current dir
	var config parser.LayoutConfig
	err := json.NewDecoder(r.Body).Decode(&config)

	err = ws.writeConfigToFile(config)
	if err != nil {
		ws.ErrorLogger.Println(err)
		return
	}
	// Call DownloadTheme function from theme.go
	err = DownloadTheme(config.ThemeURL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ws.ErrorLogger.Println("Error downloading and extracting theme:", err)
		return
	}
	FormSubmittedCh <- struct{}{}
}

func (ws *WizardServer) writeConfigToFile(config parser.LayoutConfig) error {
	configFilePath := "./site/layout/config.json"
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
		return err
	}

	marshaledJsonConfig, err := json.Marshal(config)
	if err != nil {
		ws.ErrorLogger.Fatal(err)
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err = configFile.Close()
		if err != nil {
			ws.ErrorLogger.Fatal(err)
		}
	}()
	os.WriteFile(configFilePath, marshaledJsonConfig, 0666)

	if err != nil {
		return err
	}

	return nil
}
