package anna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL     string   `yaml:"baseURL"`
	SiteTitle   string   `yaml:"siteTitle"`
	SiteScripts []string `yaml:"siteScripts"`
	Author      string   `yaml:"author"`
	ThemeURL    string   `yaml:"themeURL"`
	Navbar      []string `yaml:"navbar"`
}

type WizardServer struct {
	server *http.Server
}

func NewWizardServer(addr string) *WizardServer {
	return &WizardServer{
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (ws *WizardServer) Start() {
	http.HandleFunc("/submit", ws.handleSubmit)
	fs := http.FileServer(http.Dir("./site/static/wizard"))
	http.Handle("/", fs)
	fmt.Printf("Wizard is running at: http://localhost%s\n", ws.server.Addr)
	if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}
}

func (ws *WizardServer) Stop() error {
	return ws.server.Shutdown(context.Background())
}

func (ws *WizardServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = writeConfigToFile(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	FormSubmittedCh <- struct{}{}
}

func writeConfigToFile(config Config) error {
	configFilePath := "./site/layout/config.yml"
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the config into YAML format and write it to the file.
	if err := yaml.NewEncoder(file).Encode(&config); err != nil {
		return err
	}
	return nil
}

var FormSubmittedCh = make(chan struct{})
