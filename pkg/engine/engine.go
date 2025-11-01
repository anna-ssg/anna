package engine

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/anna-ssg/anna/v3/pkg/parser"
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

var renderPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

/*
RenderPage
fileOutPath - stores the parent directory to store rendered files, usually `site/`

pagePath - stores the path to write the given page without the prefix directory
Eg: site/content/posts/file1.html to be passed as posts/file1.html

template - stores the HTML templates parsed from the layout/ directory

templateStartString - stores the name of the template to be passed to ExecuteTemplate()
*/
// func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, template *template.Template, templateStartString string) {
// 	// Creating subdirectories if the filepath contains '/'
// 	if strings.Contains(string(pagePath), "/") {
// 		// Extracting the directory path from the page path
// 		splitPaths := strings.Split(string(pagePath), "/")
// 		filename := splitPaths[len(splitPaths)-1]
// 		pagePathWithoutFilename, _ := strings.CutSuffix(string(pagePath), filename)

// 		err := os.MkdirAll(fileOutPath+"rendered/"+pagePathWithoutFilename, 0750)
// 		if err != nil {
// 			e.ErrorLogger.Fatal(err)
// 		}
// 	}

// 	filepath := fileOutPath + "rendered/" + string(pagePath)
// 	var buffer bytes.Buffer

// 	pageData := PageData{
// 		DeepDataMerge: e.DeepDataMerge,
// 		PageURL:       pagePath,
// 	}

// 	// Storing the rendered HTML file to a buffer
// 	err := template.ExecuteTemplate(&buffer, templateStartString, pageData)
// 	if err != nil {
// 		e.ErrorLogger.Println("Error at path: ", pagePath)
// 		e.ErrorLogger.Fatal(err)
// 	}

// 	// Flushing data from the buffer to the disk
// 	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
// 	if err != nil {
// 		e.ErrorLogger.Fatal(err)
// 	}
// }

func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, template *template.Template, templateStartString string) {
	outPath := filepath.Join(fileOutPath, "rendered", string(pagePath))
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0750); err != nil {
		e.ErrorLogger.Fatal(err)
	}

	buf := renderPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer renderPool.Put(buf)

	PageData := PageData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       pagePath,
	}

	if err := template.ExecuteTemplate(buf, templateStartString, PageData); err != nil {
		e.ErrorLogger.Println("Error at path: ", pagePath)
		e.ErrorLogger.Fatal(err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	if _, err := buf.WriteTo(f); err != nil {
		_ = f.Close()
		e.ErrorLogger.Fatal(err)
	}
	if err := f.Close(); err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
