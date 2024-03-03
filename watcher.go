

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var (
	cmd        *exec.Cmd
	cmdMutex   sync.Mutex
	fileTimes  = make(map[string]time.Time)
	fileMutex  sync.Mutex
	rootDirs   = []string{"content", "cmd/ssg", "layout", "."} // Directories to monitor, so add or remove as needed
	extensions = []string{".go", ".md"}                        // File extensions to monitor,
)

func main() {
	fmt.Println("Watcher is running...")
	watch()
}

func watch() {
	for {
		fileMutex.Lock()
		for _, rootDir := range rootDirs {
			go func(rootDir string) {
				err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && hasValidExtension(path) {
						checkFile(path, info.ModTime())
					}
					return nil
				})
				if err != nil {
					fmt.Println("Error walking directory:", err)
				}
			}(rootDir)
		}
		fileMutex.Unlock()
		time.Sleep(5 * time.Second) // this is supposed to be 1s but for now let it be 5s, change it to 1 when u want it instantaneous but it will spam ur console, so remove the printlns before u do so, added those printlns while troubleshooting this
	}
}

func hasValidExtension(path string) bool {
	ext := filepath.Ext(path)
	for _, validExt := range extensions {
		if ext == validExt {
			return true
		}
	}
	return false
} //the only reason this exists is, the program was checking gitignores and shit too, so i added this check

func checkFile(path string, modTime time.Time) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	prevModTime, ok := fileTimes[path]
	if !ok || !modTime.Equal(prevModTime) {
		fileTimes[path] = modTime
		restartServer()
	} else {
		fmt.Println("No changes detected in", path)
	}
}

func restartServer() {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Error killing server:", err)
		}
	}

	fmt.Println("Starting server...")
	cmd = exec.Command("go", "run", "main.go", "--serve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server started")
}
