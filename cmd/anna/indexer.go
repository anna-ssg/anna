package anna

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FileIndex struct {
	Files []string `json:"files"`
}

func CreateIndex() {
	rootDir := "site/rendered" // Replace with the root directory of your site

	// Recursively walk through the root directory and collect all HTML files
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Create a file index
	index := FileIndex{
		Files: files,
	}

	// Convert the index to JSON
	jsonData, err := json.Marshal(index)
	if err != nil {
		panic(err)
	}

	// Open the file for writing
	file, err := os.Create("site\\static\\scripts\\index.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}
