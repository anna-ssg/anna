package helpers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var version = "2.1.0" // use variable

type Helper struct {
	ErrorLogger *log.Logger
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
	defer func() {
		err = source.Close()
		if err != nil {
			h.ErrorLogger.Fatal(err)
		}
	}()

	// Creating subdirectories if the filepath contains '/'
	if strings.Contains(destPath, "/") {
		// Extracting the directory path from the page path
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
	defer func() {
		err = destination.Close()
		if err != nil {
			h.ErrorLogger.Fatal(err)
		}
	}()

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

// BootstrapFromURL downloads a zip archive from `url` and extracts the contained
// `.../site/` directory into the current working directory's `site/` folder.
func (h *Helper) BootstrapFromURL(url string) error {
    h.ErrorLogger.Printf("Downloading %s\n", url)

    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to download %s: %w", url, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("download failed: %s", resp.Status)
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
    if err != nil {
        return fmt.Errorf("failed to open zip: %w", err)
    }

    for _, f := range zr.File {
        // We only extract files under the repo's `site/` folder
        if idx := strings.Index(f.Name, "/site/"); idx != -1 {
            rel := f.Name[idx+len("/site/"):] // path inside site/
            destPath := filepath.Join("site", rel)

            if f.FileInfo().IsDir() {
                if err := os.MkdirAll(destPath, 0755); err != nil {
                    return fmt.Errorf("failed to create dir %s: %w", destPath, err)
                }
                continue
            }

            if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
                return fmt.Errorf("failed to create parent dir for %s: %w", destPath, err)
            }

            rc, err := f.Open()
            if err != nil {
                return fmt.Errorf("failed to open zipped file %s: %w", f.Name, err)
            }

            out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
            if err != nil {
                rc.Close()
                return fmt.Errorf("failed to create file %s: %w", destPath, err)
            }

            if _, err := io.Copy(out, rc); err != nil {
                out.Close()
                rc.Close()
                return fmt.Errorf("failed to copy file %s: %w", destPath, err)
            }
            out.Close()
            rc.Close()
        }
    }

    h.ErrorLogger.Println("Bootstrapped site/ from archive")
    return nil
}
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
