package anna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func handleInputValues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	themeURL := r.Form.Get("themeURL")

	switch themeURL {
	case "/static/styles/sudhir.css":
		executeCurlCommand("curl https://sudhir.live/static/style.css https://cdn.jsdelivr.net/npm/water.css@2/out/water.css --output-dir site/static -o style.css")
	case "/static/styles/hegde.css":
		executeCurlCommand("curl https://cdn.jsdelivr.net/npm/highlightjs-themes@1.0.0/tomorrow-night.css https://hegde.live/static/style.css --output-dir site/static -o style.css")
	case "/static/styles/nathan.css":
		executeCurlCommand("curl https://polarhive.net/assets/main.css https://polarhive.net/assets/style.css --output-dir site/static -o style.css")
	case "/static/styles/new.css":
		executeCurlCommand("curl https://cdn.jsdelivr.net/npm/@exampledev/new.css@1.1.2/new.min.css -o style.css--output-dir site/static")
	case "/static/styles/pico.css":
		executeCurlCommand("curl https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.7.0/build/styles/tokyo-night-dark.min.css https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css --output-dir site/static -o style.css")
	case "/static/styles/simple.css":
		executeCurlCommand("curl https://cdn.simplecss.org/simple.css --output-dir site/static -o style.css")
	default:
		// Handle default case
		executeCurlCommand("curl default.css")
	}

	w.WriteHeader(http.StatusOK)
}

func executeCurlCommand(curlCommand string) {
	cmd := exec.Command(curlCommand)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute curl command: %v", err)
	}
}
