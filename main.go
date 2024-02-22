package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

func main() {
	parseMarkdown()
}

func parseMarkdown() {
	// Listing all files in the ./content/ directory
	files, err := os.ReadDir("./content/")
	if err != nil {
		log.Fatal(err)
	}

	var mdFilesName []string
	var mdFilesPath []string
	var mdFilesContent [][]byte
	var renderedHTML []bytes.Buffer

	// Storing the markdown file pointers
	for _, filename := range files {
		mdFilesName = append(mdFilesName, filename.Name())

		filepath := "content/" + filename.Name()
		mdFilesPath = append(mdFilesPath, filepath)
	}

	// Reading the markdown files into memory
	for _, filepath := range mdFilesPath {
		content, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}

		mdFilesContent = append(mdFilesContent, content)
	}

	// Parsing markdown to HTML
	for _, filecontent := range mdFilesContent {
		var buffer bytes.Buffer
		if err := goldmark.Convert(filecontent, &buffer); err != nil {
			log.Fatal(err)
		}

		renderedHTML = append(renderedHTML, buffer)
	}

	// Creating the "rendered" directory if not present
	err = os.MkdirAll("rendered/", 0750)
	if err != nil {
		log.Fatal(err)
	}

	// Writing parsed HTML to disk
	for i, page := range renderedHTML {
		filename, _ := strings.CutSuffix(mdFilesName[i], ".md")
		filepath := "rendered/" + filename + ".html"

		err := os.WriteFile(filepath, page.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
