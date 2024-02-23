package ssg

import (
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

type Generator struct {
	ErrorLogger     *log.Logger
	mdFilesName     []string
	mdFilesPath     []string
	mdFilesContent  [][]byte
	layoutFilesPath []string
	renderedHTML    []bytes.Buffer
}

func (g *Generator) ParseMarkdown() {
	// Listing all files in the ./content/ directory
	files, err := os.ReadDir("./content/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Storing the markdown file names and paths
	for _, filename := range files {
		g.mdFilesName = append(g.mdFilesName, filename.Name())

		filepath := "content/" + filename.Name()
		g.mdFilesPath = append(g.mdFilesPath, filepath)
	}

	// Reading the markdown files into memory
	for _, filepath := range g.mdFilesPath {
		content, err := os.ReadFile(filepath)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		g.mdFilesContent = append(g.mdFilesContent, content)
	}

	// Parsing markdown to HTML
	for _, filecontent := range g.mdFilesContent {
		var buffer bytes.Buffer
		if err := goldmark.Convert(filecontent, &buffer); err != nil {
			g.ErrorLogger.Fatal(err)
		}

		g.renderedHTML = append(g.renderedHTML, buffer)
	}
}

// Write rendered HTML to disk
func (g *Generator) RenderSite() {
	// Creating the "rendered" directory if not present
	err := os.MkdirAll("rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	templ, err := template.ParseFiles("./layout/layout.html")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Writing each parsed markdown file as a separate HTML file
	for i, page := range g.renderedHTML {

		filename, _ := strings.CutSuffix(g.mdFilesName[i], ".md")
		filepath := "rendered/" + filename + ".html"
		var buffer bytes.Buffer

		// Storing the rendered HTML file to a buffer
		err = templ.ExecuteTemplate(&buffer, "layout", page.String())
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		// Flushing data from the buffer to the disk
		err := os.WriteFile(filepath, buffer.Bytes(), 0666)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
	}
}
