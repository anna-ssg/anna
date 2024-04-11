package engine

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/acmpesuecc/anna/pkg/parser"
)

// This struct holds all of the ssg data
type MergedSiteData struct {
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
}

type Engine struct {
	// Stores the merged ssg data
	DeepDataMerge MergedSiteData

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

/*
fileOutPath - stores the parent directory to store rendered files, usually `site/`

pagePath - stores the path to write the given page without the prefix directory
Eg: site/content/posts/file1.html to be passed as posts/file1.html

pageTemplateData - an interface that accepts any type of data to be passed to ExecuteTemplate()

template - stores the HTML templates parsed from the layout/ directory

templateStartString - stores the name of the template to be passed to ExecuteTemplate()
*/
func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, pageTemplateData interface{}, template *template.Template, templateStartString string) {
	// Creating subdirectories if the filepath contains '/'
	if strings.Contains(string(pagePath), "/") {
		// Extracting the directory path from the page path
		splitPaths := strings.Split(string(pagePath), "/")
		filename := splitPaths[len(splitPaths)-1]
		pagePathWithoutFilename, _ := strings.CutSuffix(string(pagePath), filename)

		err := os.MkdirAll(fileOutPath+"rendered/"+pagePathWithoutFilename, 0750)
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}

	filepath := fileOutPath + "rendered/" + string(pagePath)
	var buffer bytes.Buffer

	// Storing the rendered HTML file to a buffer
	err := template.ExecuteTemplate(&buffer, templateStartString, pageTemplateData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
