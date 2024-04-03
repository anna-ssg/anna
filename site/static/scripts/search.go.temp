package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"syscall/js"
)

func searchFiles(this js.Value, inputs []js.Value) interface{} {
	// Retrieve the search query from JavaScript
	query := inputs[0].Get("value").String()

	// Read index.json
	indexFile := "index.json"
	indexData, err := ioutil.ReadFile(indexFile)
	if err != nil {
		return err.Error()
	}

	// Parse index.json
	var indexMap map[string][]string
	if err := json.Unmarshal(indexData, &indexMap); err != nil {
		return err.Error()
	}

	// Search through files
	var results []string
	for _, files := range indexMap {
		for _, file := range files {
			fileContent, err := ioutil.ReadFile(file)
			if err != nil {
				return err.Error()
			}
			if strings.Contains(string(fileContent), query) {
				results = append(results, file)
			}
		}
	}

	// Return results to JavaScript
	return results
}

func main() {
	// Register the Go function as a JavaScript function
	js.Global().Set("searchFiles", js.FuncOf(searchFiles))

	// Keep the program running
	select {}
}
