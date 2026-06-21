package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/anna-ssg/anna/v4/pkg/logger"
	embeddedsite "github.com/anna-ssg/anna/v4/site"
)

type Helper struct {
	ErrorLogger *logger.Logger
}

/*
CopyDirectoryContents
Copies the contents of the dirPath directory to outDirPath
*/
func (h *Helper) CopyDirectoryContents(dirPath string, outDirPath string) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(outDirPath, 0750)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			h.CopyDirectoryContents(dirPath+entry.Name()+"/", outDirPath+entry.Name()+"/")
		} else {
			h.CopyFiles(dirPath+entry.Name(), outDirPath+entry.Name())
		}
	}
}

func (h *Helper) CopyFiles(srcPath string, destPath string) {
	source, err := os.Open(srcPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
	defer source.Close()

	if strings.Contains(destPath, "/") {
		splitPaths := strings.Split(destPath, "/")
		filename := splitPaths[len(splitPaths)-1]
		pagePathWithoutFilename, _ := strings.CutSuffix(destPath, filename)

		err := os.MkdirAll(pagePathWithoutFilename, 0750)
		if err != nil {
			h.ErrorLogger.Fatal(err)
		}
	}

	destination, err := os.Create(destPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
}

func (h *Helper) CreateRenderedDir(fileOutPath string) {
	err := os.RemoveAll(fileOutPath + "rendered/")
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(fileOutPath+"rendered/", 0750)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
}

// BootstrapEmbedded extracts the embedded starter site into ./site.
func (h *Helper) BootstrapEmbedded(overwrite bool) error {
	if overwrite {
		if err := os.RemoveAll("site"); err != nil {
			return err
		}
	}

	if err := fs.WalkDir(embeddedsite.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the virtual root.
		if path == "." {
			return nil
		}

		dest := filepath.Join("site", path)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		src, err := embeddedsite.FS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		out, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			return err
		}

		return out.Close()
	}); err != nil {
		return err
	}

	h.ErrorLogger.Println("Bootstrapped embedded site.")
	return nil
}

// EnsureSiteExists bootstraps automatically if ./site is missing.
func (h *Helper) EnsureSiteExists() error {
	if _, err := os.Stat("site"); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	fmt.Println("No site/ directory found.")
	fmt.Print("Create a new starter site here? (Y/n): ")

	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp != "" && resp != "y" && resp != "yes" {
		return errors.New("aborted")
	}

	return h.BootstrapEmbedded(false)
}
