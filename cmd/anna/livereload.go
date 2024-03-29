package anna

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/acmpesuecc/anna/pkg/helpers"
)

type liveReload struct {
	errorLogger *log.Logger
	fileTimes   map[string]time.Time

	// Directories to monitor, so add or remove as needed
	rootDirs []string

	// File extensions to monitor
	extensions []string

	serverRunning bool
}

func newLiveReload() *liveReload {
	lr := liveReload{
		errorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		fileTimes:   make(map[string]time.Time),
		rootDirs:    []string{helpers.SiteDataPath},
		extensions:  []string{".go", ".md"},
	}
	return &lr
}

func (cmd *Cmd) StartLiveReload() {
	fmt.Println("Live Reload is active")
	lr := newLiveReload()
	go lr.startServer(cmd.Addr)

	for {
		for _, rootDir := range lr.rootDirs {
			if lr.traverseDirectory(rootDir) {
				cmd.VanillaRender()
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
		if !info.IsDir() && lr.hasValidExtension(path) {
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
	fmt.Print("Serving content at: http://localhost:", addr, "\n")
	err := http.ListenAndServe(":"+addr, http.FileServer(http.Dir(helpers.SiteDataPath+"./rendered")))
	if err != nil {
		lr.errorLogger.Fatal(err)
	}
}
