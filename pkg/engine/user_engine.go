package engine

import (
	"html/template"
	"runtime"
	"sync"
)

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
