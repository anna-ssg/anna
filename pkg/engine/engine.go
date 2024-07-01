package engine

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/anna-ssg/anna/v2/pkg/parser"
)

// DeepDataMerge This struct holds all the ssg data
type DeepDataMerge struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates map[template.URL]parser.TemplateData

	// Templates stores the template data of all tag sub-pages of the site
	Tags map[template.URL]parser.TemplateData

	// K-V pair storing all templates corresponding to a particular tag in the site
	TagsMap map[template.URL][]parser.TemplateData

	// Stores data parsed from layout/config.yml
	LayoutConfig parser.LayoutConfig

	// Templates stores the template data of all collection sub-pages of the site
	Collections map[template.URL]parser.TemplateData

	// K-V pair storing all templates corresponding to a particular collection in the site
	CollectionsMap map[template.URL][]parser.TemplateData

	// K-V pair storing the template layout name for a particular collection in the site
	CollectionsSubPageLayouts map[template.URL]string

	//Stores all the notes
	Notes map[template.URL]parser.Note

	//Stores the links of each note to other notes
	LinkStore map[template.URL][]*parser.Note

	// Stores the index generated for search functionality
	JSONIndex map[template.URL]JSONIndexTemplate
}

type Engine struct {
	// Stores the merged ssg data
	DeepDataMerge DeepDataMerge

	// Common logger for all engine functions
	ErrorLogger *log.Logger

	// The path to the directory being rendered
	SiteDataPath string
}

type PageData struct {
	DeepDataMerge DeepDataMerge

	PageURL template.URL
}

// JSONIndexTemplate This structure is solely used for storing the JSON index
type JSONIndexTemplate struct {
	CompleteURL template.URL
	Frontmatter parser.Frontmatter
	Tags        []string
}

/*
RenderPage
fileOutPath - stores the parent directory to store rendered files, usually `site/`

pagePath - stores the path to write the given page without the prefix directory
Eg: site/content/posts/file1.html to be passed as posts/file1.html

template - stores the HTML templates parsed from the layout/ directory

templateStartString - stores the name of the template to be passed to ExecuteTemplate()
*/
func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, template *template.Template, templateStartString string) {
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

	pageData := PageData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       pagePath,
	}

	// Storing the rendered HTML file to a buffer
	err := template.ExecuteTemplate(&buffer, templateStartString, pageData)
	if err != nil {
		e.ErrorLogger.Println("Error at path: ", pagePath)
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
