package engine

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/acmpesuecc/anna/pkg/parser"
)

type postsTemplateData struct {
	Posts []parser.TemplateData
	parser.TemplateData
}

func (e *Engine) RenderUserDefinedPages(fileOutPath string, templ *template.Template) {
	for _, templData := range e.Templates {
		fileInPath, _ := strings.CutSuffix(string(templData.CompleteURL), ".html")
		e.RenderPage(fileOutPath, template.URL(fileInPath), templData, templ, "page")
	}
}

func (e *Engine) RenderEngineGeneratedFiles(fileOutPath string, templ *template.Template) {
	// Rendering "posts.html"
	var postsBuffer bytes.Buffer

	postsData := postsTemplateData{
		Posts: e.Posts,
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{Title: "Posts"},
			Layout:      e.LayoutConfig,
		},
	}

	err := templ.ExecuteTemplate(&postsBuffer, "posts", postsData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile(fileOutPath+"rendered/posts.html", postsBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

/*
func ParallelCode() {
	// Adhesh's code
	var wg sync.WaitGroup
	concurrency := 3
	// Each goroutine handles 3 files at a time
	semaphore := make(chan struct{}, concurrency)

	files := make([]string, 0, len(e.Templates))
	for pagePath := range e.Templates {
		files = append(files, string(pagePath))
	}

	for _, file := range files {
		wg.Add(1)
		// Acquire semaphore
		semaphore <- struct{}{}

		go func(file string) {
			defer func() {
				// Release semaphore
				<-semaphore
				wg.Done()
			}()

			pageURL := template.URL(file)
			templateData := e.Templates[pageURL]
			e.RenderPage(fileOutPath, pageURL, templateData, templ, "page")
		}(file)
	}
	wg.Wait()
}
*/
