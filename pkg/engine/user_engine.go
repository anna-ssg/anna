package engine

import (
	"bytes"
	"html/template"
	"os"
	"runtime"
	"sync"

	"github.com/anna-ssg/anna/v2/pkg/parser"
)

type postsTemplateData struct {
	DeepDataMerge DeepDataMerge
	PageURL       template.URL
	TemplateData  parser.TemplateData
}

func (e *Engine) RenderEngineGeneratedFiles(fileOutPath string, templates *template.Template) {
	// Rendering "posts.html"
	var postsBuffer bytes.Buffer

	postsData := postsTemplateData{
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{Title: "Posts"},
		},
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       "posts.html",
	}

	err := templates.ExecuteTemplate(&postsBuffer, "posts", postsData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile(fileOutPath+"rendered/posts.html", postsBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) RenderUserDefinedPages(fileOutPath string, templates *template.Template) {
	numCPU := runtime.NumCPU()
	numTemplates := len(e.DeepDataMerge.Templates)
	concurrency := numCPU * 2 // Adjust the concurrency factor based on system hardware resources

	if numTemplates < concurrency {
		concurrency = numTemplates
	}

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.DeepDataMerge.Templates {
		templateURLs = append(templateURLs, string(templateURL))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, templateURL := range templateURLs {
		if templateURL == ".html" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(templateURL string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			e.RenderPage(fileOutPath, template.URL(templateURL), templates, e.DeepDataMerge.Templates[template.URL(templateURL)].Frontmatter.Layout)
		}(templateURL)
	}

	wg.Wait()
}
