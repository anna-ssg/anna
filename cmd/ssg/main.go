package ssg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
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
	staticFilesPath []string
	renderedHTML    []bytes.Buffer
}

func (g *Generator) parseMarkdown() {
	// Listing all files in the content/ directory
	files, err := os.ReadDir("content/")
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
	g.parseMarkdown()
	g.copyStaticContent()

	// Creating the "rendered" directory if not present
	err := os.MkdirAll("rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	templ, err := template.ParseFiles("layout/layout.html")
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

// Copies the 'static/' directory and its contents to 'rendered/'
func (g *Generator) copyStaticContent() {
	files, err := os.ReadDir("static/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Storing the static file names and paths
	for _, filename := range files {
		filepath := "static/" + filename.Name()
		g.staticFilesPath = append(g.staticFilesPath, filepath)
	}

	err = os.MkdirAll("rendered/static", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Copying the contents of the 'static/' directory
	for _, filepath := range g.staticFilesPath {
		source, err := os.Open(filepath)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
		defer source.Close()

		new_filepath := "rendered/" + filepath
		destination, err := os.Create(new_filepath)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

	}
}

// Serves the rendered files over the address 'addr'
func (g *Generator) ServeSite(addr string) {
	fmt.Println("Serving content at", addr)
	err := http.ListenAndServe(addr, http.FileServer(http.Dir("./rendered")))
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}
