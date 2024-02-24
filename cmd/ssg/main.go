package ssg

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

type Page struct {
	Frontmatter Frontmatter
	Body        string
}

type Generator struct {
	ErrorLogger     *log.Logger
	mdFilesName     []string
	mdFilesPath     []string
	mdParsed        []Page
	layoutFilesPath []string
	staticFilesPath []string
}

// Write rendered HTML to disk
func (g *Generator) RenderSite() {
	g.readMarkdownFiles()
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
	for i, page := range g.mdParsed {

		filename, _ := strings.CutSuffix(g.mdFilesName[i], ".md")
		filepath := "rendered/" + filename + ".html"
		var buffer bytes.Buffer

		// Storing the rendered HTML file to a buffer
		err = templ.ExecuteTemplate(&buffer, "layout", page)
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

// Serves the rendered files over the address 'addr'
func (g *Generator) ServeSite(addr string) {
	fmt.Println("Serving content at", addr)
	err := http.ListenAndServe(addr, http.FileServer(http.Dir("./rendered")))
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) readMarkdownFiles() {
	// Listing all files in the content/ directory
	files, err := os.ReadDir("content/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Storing the markdown file names and paths
	for _, filename := range files {
		if !strings.HasSuffix(filename.Name(), ".md") {
			continue
		}

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

		frontmatter, body := g.parseMarkdownContent(string(content))

		page := Page{
			Frontmatter: frontmatter,
			Body:        body,
		}

		g.mdParsed = append(g.mdParsed, page)
	}
}

func (g *Generator) parseMarkdownContent(filecontent string) (Frontmatter, string) {
	var parsedFrontmatter Frontmatter
	var markdown string

	// Find the '---' tags for frontmatter in the markdown file
	re := regexp.MustCompile(`(---[\S\s]*---)`)
	frontmatter := re.FindString(filecontent)

	if frontmatter != "" {
		// Parsing YAML frontmatter
		err := yaml.Unmarshal([]byte(frontmatter), &parsedFrontmatter)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		// Splitting and storing pure markdown content separately
		markdown = strings.Split(filecontent, "---")[2]
	} else {
		markdown = filecontent
	}

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		g.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String()
}

// Copies the contents of the 'static/' directory to 'rendered/'
func (g *Generator) copyStaticContent() {
	g.copyDirectoryContents("static/", "rendered/static/")
}
