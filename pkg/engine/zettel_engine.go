package engine

import (
	"bytes"
	"html/template"
	"os"
	"runtime"
	"sync"

	"github.com/acmpesuecc/anna/pkg/parser"
)

type notesTemplateData struct {
	DeepDataMerge DeepDataMerge
	PageURL       template.URL
	TemplateData  parser.TemplateData
}

func (e *Engine) RenderNotes(fileOutPath string, templ *template.Template) {
	// templ.Funcs(funcMap template.FuncMap)

	numCPU := runtime.NumCPU()
	numTemplates := len(e.DeepDataMerge.Notes)
	concurrency := numCPU * 2 // Adjust the concurrency factor based on system hardware resources

	if numTemplates < concurrency {
		concurrency = numTemplates
	}

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.DeepDataMerge.Notes {
		templateURLs = append(templateURLs, string(templateURL))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, url := range templateURLs {
		if url == ".html" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(templateURL string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			e.RenderPage(fileOutPath, template.URL(url), templ, "note")
		}(url)
	}

	wg.Wait()
}

func (e *Engine) GenerateLinkStore() {
	for url, note := range e.DeepDataMerge.Notes {
		for _, linkURL := range note.LinkedNoteURLs {
			linkNote, ok := e.DeepDataMerge.Notes[linkURL]
			if ok {
				e.DeepDataMerge.LinkStore[url] = append(e.DeepDataMerge.LinkStore[url], &linkNote)
			}
		}
	}
}

func (e *Engine) GenerateNoteRoot(fileOutPath string, templ *template.Template) {
	var buffer bytes.Buffer

	notesTemplateData := notesTemplateData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       "notes.html",
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{
				Title:       "Curated Notes",
				Description: "Currated heads of various zettles part of the page",
			},
		},
	}

	err := templ.ExecuteTemplate(&buffer, "notes-root", notesTemplateData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	err = os.WriteFile(fileOutPath+"rendered/notes.html", buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

}
