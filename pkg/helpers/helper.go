package helpers

import (
	"io"
	"log"
	"os"
	"strings"
)

const SiteDataPath string = "site/"

var version string = "2.0.0" // use variable

type Helper struct {
	ErrorLogger *log.Logger
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

	// Creating subdirectories if the filepath contains '/'
	if strings.Contains(string(destPath), "/") {
		// Extracting the directory path from the page path
		splitPaths := strings.Split(string(destPath), "/")
		filename := splitPaths[len(splitPaths)-1]
		pagePathWithoutFilename, _ := strings.CutSuffix(string(destPath), filename)

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

// func (h *Helper) Bootstrap() {
// 	fmt.Println("Are you sure you want to proceed with the bootstrap process? (y/n)")
// 	var confirm string
// 	fmt.Scanln(&confirm)
// 	if confirm != "y" {
// 		fmt.Println("Bootstrap process cancelled.")
// 		return
// 	}
// 	url := fmt.Sprintf("https://github.com/acmpesuecc/anna/archive/refs/tags/v%s.zip", version)

// 	output, err := os.Create("anna-dl.zip")
// 	if err != nil {
// 		fmt.Println("Error creating output file:", err)
// 		return
// 	}
// 	defer output.Close()
// 	response, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("Error downloading:", err)
// 		return
// 	}
// 	defer response.Body.Close()
// 	if response.StatusCode != http.StatusOK {
// 		fmt.Println("Error: Server returned status", response.Status)
// 		return
// 	}

// 	bar := pb.Full.Start64(response.ContentLength)
// 	defer bar.Finish()
// 	reader := bar.NewProxyReader(response.Body)

// 	_, err = io.Copy(output, reader)
// 	if err != nil {
// 		fmt.Println("Error copying:", err)
// 		return
// 	}
// 	ext(version)
// }

// func ext(version string) {
// 	zipFilePath := "anna-dl.zip"
// 	zipFile, err := zip.OpenReader(zipFilePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer zipFile.Close()
// 	for _, file := range zipFile.File {
// 		zippedFile, err := file.Open()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer zippedFile.Close()

// 		extractedFilePath := filepath.Join(".", file.Name) // Extract to the current directory
// 		if file.FileInfo().IsDir() {
// 			err := os.MkdirAll(extractedFilePath, os.ModePerm)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		} else {
// 			extractedFile, err := os.Create(extractedFilePath)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer extractedFile.Close()

// 			_, err = io.Copy(extractedFile, zippedFile)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}

// 	// Bootstrap the base theme
// 	helper := &Helper{
// 		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
// 	}
// 	helper.CopyDirectoryContents(fmt.Sprintf("anna-%s/site/", version), "site/")
// 	err = os.Remove("anna-dl.zip")
// 	if err != nil {
// 		log.Fatal("Error cleaning-up zip file:", err)
// 	}
// 	err = os.RemoveAll(fmt.Sprintf("anna-%s", version))
// 	if err != nil {
// 		log.Fatal("Error deleting directory:", err)
// 	}
// 	log.Println("Bootstrapped base theme")
// }
