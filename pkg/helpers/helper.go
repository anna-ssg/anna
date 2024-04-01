package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const SiteDataPath string = "site/"

type Helper struct {
	ErrorLogger  *log.Logger
	SiteDataPath string
}

// Copies the contents of the dirPath directory to outDirPath
func (h *Helper) CopyDirectoryContents(dirPath string, outDirPath string) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(outDirPath, 0750)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}

	// Copying the contents of the dirPath directory
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

	err = os.MkdirAll(fileOutPath+"rendered/layout/", 0750)
	if err != nil {
		h.ErrorLogger.Fatal(err)
	}
}

func (h *Helper) Bootstrap() {
	log.Println("Downloading base theme")
	url := "https://github.com/acmpesuecc/anna/archive/refs/heads/main.zip"
	output, err := os.Create("anna-repo.zip")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading:", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned status", response.Status)
		return
	}

	bar := pb.Full.Start64(response.ContentLength)
	defer bar.Finish()
	reader := bar.NewProxyReader(response.Body)

	_, err = io.Copy(output, reader)
	if err != nil {
		fmt.Println("Error copying:", err)
		return
	}
	log.Println("Downloaded successfully")
}
