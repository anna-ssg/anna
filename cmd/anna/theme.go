package anna

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	baseURL = "https://github.com/anna-ssg/themes/releases/download"
	tagVer  = "v3.0"
	destDir = "site"
)

// DownloadTheme downloads and extracts a theme from a URL based on themeName.
func DownloadTheme(themeName string) error {
	zipURL := fmt.Sprintf("%s/%s/%s.zip", baseURL, tagVer, themeName)
	zipFile := fmt.Sprintf("%s.zip", themeName)

	fmt.Printf("Downloading %s...\n", zipURL)
	if err := downloadFile(zipURL, zipFile); err != nil {
		return fmt.Errorf("error downloading theme '%s': %v", themeName, err)
	}

	fmt.Println("Extracting files...")
	if err := unzip(zipFile, destDir); err != nil {
		return fmt.Errorf("error extracting files: %v", err)
	}

	fmt.Printf("Theme '%s' extracted to '%s' directory successfully.\n", themeName, destDir)

	if err := os.Remove(zipFile); err != nil {
		fmt.Printf("Warning: error deleting zip file '%s': %v\n", zipFile, err)
	}

	return nil
}

// downloadFile downloads a file from a URL and saves it to filepath.
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzip extracts a zip file (src) to a destination directory (dest).
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)

		// Ensure extraction path is safe
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
