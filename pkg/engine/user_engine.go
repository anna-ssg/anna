package engine

import (
	"html/template"
	"runtime"
	"strings"
	"sync"
)

func (e *Engine) RenderUserDefinedPages(fileOutPath string, templates *template.Template) {
	numCPU := runtime.NumCPU()
	numTemplates := len(e.DeepDataMerge.Pages)
	concurrency := min(numCPU*2, numTemplates)

	templateURLs := make([]string, 0, numTemplates)
	for templateURL := range e.DeepDataMerge.Pages {
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

			templateType := e.DeepDataMerge.Pages[template.URL(templateURL)].Frontmatter.Layout
			if templateType == "" {
				if strings.HasSuffix(templateURL, "index.html") {
					var ok bool
					if templateType, ok = e.DeepDataMerge.CollectionsSubPageLayouts[template.URL(templateURL)]; !ok {
						templateType = e.DeepDataMerge.Config.DefaultCollectionIndexLayout
					}
				} else {
					// Fallback to the default template used for rendering all user defined content pages
					templateType = "page"
				}
			}

			e.RenderPage(fileOutPath, template.URL(templateURL), templates, templateType)
		}(templateURL)
	}

	wg.Wait()
}
