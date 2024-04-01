package engine

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/acmpesuecc/anna/pkg/parser"
)

type Engine struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates map[template.URL]parser.TemplateData

	// K-V pair storing all templates correspoding to a particular tag in the site
	TagsMap map[string][]parser.TemplateData

	// Stores data parsed from layout/config.yml
	LayoutConfig parser.LayoutConfig

	// Posts contains the template data of files in the posts directory
	Posts []parser.TemplateData

	// Stores the index generated for search functionality
	JSONIndex map[template.URL]JSONIndexTemplate

	// Common logger for all engine functions
	ErrorLogger *log.Logger
}

// This structure is solely used for storing the JSON index
type JSONIndexTemplate struct {
	CompleteURL              template.URL
	FilenameWithoutExtension string
	Frontmatter              parser.Frontmatter
	Tags                     []string
}

// fileOutPath for main.go should be refering to helpers.SiteDataPath
func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, pageTemplateData parser.TemplateData, templ *template.Template, templateStartString string) {
	// Creating subdirectories if the filepath contains '/'
	dirPath := ""
	if strings.Contains(string(pagePath), "/") {
		// Extracting the directory path from the page path
		dirPath, _ := strings.CutSuffix(string(pagePath), pageTemplateData.FilenameWithoutExtension)
		dirPath = fileOutPath + "rendered/" + dirPath

		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}

	filename, _ := strings.CutSuffix(string(pagePath), ".md")
	filepath := fileOutPath + "rendered/" + dirPath + filename + ".html"
	var buffer bytes.Buffer

	// Storing the rendered HTML file to a buffer
	err := templ.ExecuteTemplate(&buffer, templateStartString, pageTemplateData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
