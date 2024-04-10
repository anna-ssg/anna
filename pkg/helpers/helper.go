package helpers

import (
	"archive/zip"
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
}

func (h *Helper) Bootstrap() {
	fmt.Println("Are you sure you want to proceed with the bootstrap process? (y/n)")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Println("Bootstrap process cancelled.")
		return
	}
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
	ext()
}

func ext() {
	zipFilePath := "anna-repo.zip"

	// Open the zip file
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return
	}
	defer zipFile.Close()

	// Extract each file from the zip archive
	for _, file := range zipFile.File {
		// Open file from zip archive
		zippedFile, err := file.Open()
		if err != nil {
			fmt.Println("Error opening zipped file:", err)
			return
		}
		defer zippedFile.Close()

		// Create the file in the extraction directory
		extractedFilePath := file.Name
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			extractedFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				fmt.Println("Error creating extracted file:", err)
				return
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, zippedFile)
			if err != nil {
				fmt.Println("Error extracting file:", err)
				return
			}
		}
	}

	helper := &Helper{
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
	helper.CopyDirectoryContents("anna-main/site/", "site/")

	if err := os.RemoveAll("anna-main"); err != nil {
		log.Println("Error deleting directory:", err)
	}

	if err := os.RemoveAll("anna-repo.zip"); err != nil {
		log.Println("Error deleting zip file:", err)
	}
	log.Println("Bootstrapped base theme")

}
