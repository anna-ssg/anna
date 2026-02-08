package anna

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var reloadPageBool atomic.Bool

type liveReload struct {
	errorLogger *log.Logger
	fileTimes   map[string]time.Time

	// Directories to monitor, so add or remove as needed
	rootDirs []string

	// File extensions to monitor
	extensions []string

	serverRunning bool

	siteDataPath string
}

func newLiveReload(siteDataPath string) *liveReload {
	lr := liveReload{
		errorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		fileTimes:   make(map[string]time.Time),
		rootDirs:    []string{siteDataPath},
		// empty/extensions==nil means watch all files
		extensions:   nil,
		siteDataPath: siteDataPath,
	}
	return &lr
}

func (cmd *Cmd) StartLiveReload(siteDataPath string) {
	fmt.Println("Live Reload is active")
	lr := newLiveReload(siteDataPath)
	go lr.startServer(cmd.Addr)

	for {
		for _, rootDir := range lr.rootDirs {
			if lr.traverseDirectory(rootDir) {
				cmd.VanillaRender(lr.siteDataPath)
				reloadPageBool.CompareAndSwap(false, true)
			}
		}
		if !lr.serverRunning {
			lr.serverRunning = true
		}
		time.Sleep(time.Second)
	}
}

func (lr *liveReload) traverseDirectory(rootDir string) bool {
	filesChanged := false
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip generated and VCS directories to avoid render loops and noisy events
		if info.IsDir() {
			name := info.Name()
			if name == "rendered" || name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Consider all files if extensions list is nil/empty, otherwise filter by extension
		if lr.hasValidExtension(path) {
			if lr.checkFile(path, info.ModTime()) {
				filesChanged = true
				return nil
			}
		}
		return nil
	})
	if err != nil {
		lr.errorLogger.Fatal(err)
	}
	return filesChanged
}

func (lr *liveReload) hasValidExtension(path string) bool {
	// if no extensions are configured, accept all files
	if len(lr.extensions) == 0 {
		return true
	}
	ext := filepath.Ext(path)
	for _, validExt := range lr.extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func (lr *liveReload) checkFile(path string, modTime time.Time) bool {
	prevModTime, ok := lr.fileTimes[path]
	if !ok || !modTime.Equal(prevModTime) {
		lr.fileTimes[path] = modTime
		if lr.serverRunning {
			fmt.Println("The following file has changed: ", path)
			print("-----------------------------\n")
		}
		return true
	}
	return false
}

func (lr *liveReload) startServer(addr string) {
	fmt.Printf("Serving content at address: http://%s\n", addr)
	fmt.Printf("Profile data can be viewed at: http://%s\n", addr+"/debug/pprof")
	http.Handle("/", http.FileServer(http.Dir(lr.siteDataPath+"./rendered")))
	http.HandleFunc("/events", eventsHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		lr.errorLogger.Fatal(err)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if !reloadPageBool.Load() {
		return
	}

	event := "event:\ndata:\n\n"
	_, err := w.Write([]byte(event))
	if err != nil {
		log.Fatal(err)
	}
	w.(http.Flusher).Flush()

	reloadPageBool.CompareAndSwap(true, false)
}
