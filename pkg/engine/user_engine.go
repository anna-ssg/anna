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

	if concurrency < 1 {
		concurrency = 1
	}
	// templateURLs := make([]string, 0, numTemplates)
	// for templateURL := range e.DeepDataMerge.Templates {
	// 	templateURLs = append(templateURLs, string(templateURL))
	// }

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.DeepDataMerge.Templates {
		s := string(templateURL)
		if s == ".html" {
			continue
		}
		templateURLs = append(templateURLs, s)
	}

	// var wg sync.WaitGroup
	// semaphore := make(chan struct{}, concurrency)

	tasks := make(chan string, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for templateURL := range tasks {
				URL := template.URL(templateURL)
				layout := e.DeepDataMerge.Templates[URL].Frontmatter.Layout
				e.RenderPage(fileOutPath, URL, templates, layout)
			}
		}()
	}
	for _, v := range templateURLs {
		tasks <- v
	}
	close(tasks)

	// for _, templateURL := range templateURLs {
	// 	if templateURL == ".html" {
	// 		continue
	// 	}

	// 	wg.Add(1)
	// 	semaphore <- struct{}{}

	// 	go func(templateURL string) {
	// 		defer func() {
	// 			<-semaphore
	// 			wg.Done()
	// 		}()

	// 		e.RenderPage(fileOutPath, template.URL(templateURL), templates, e.DeepDataMerge.Templates[template.URL(templateURL)].Frontmatter.Layout)
	// 	}(templateURL)
	// }

	wg.Wait()
}
