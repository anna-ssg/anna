package engine

import (
	"bytes"
	"html/template"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/acmpesuecc/anna/pkg/parser"
)

type postsTemplateData struct {
	Posts []parser.TemplateData
	parser.TemplateData
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

func (e *Engine) RenderUserDefinedPages(fileOutPath string, templ *template.Template) {
	numCPU := runtime.NumCPU()
	numTemplates := len(e.Templates)
	concurrency := numCPU * 2 // Adjust the concurrency factor based on system hardware resources

	if numTemplates < concurrency {
		concurrency = numTemplates
	}

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.Templates {
		templateURLs = append(templateURLs, string(templateURL))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, templateURL := range templateURLs {
		templData := e.Templates[template.URL(templateURL)]
		fileInPath := strings.TrimSuffix(string(templData.CompleteURL), ".html")
		if fileInPath == "" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(templateURL string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			e.RenderPage(fileOutPath, template.URL(fileInPath), templData, templ, "page")
		}(templateURL)
	}

	wg.Wait()
}
