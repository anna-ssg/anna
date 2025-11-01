package engine

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// func (e *Engine) RenderUserDefinedPages(fileOutPath string, templates *template.Template) {
// 	numCPU := runtime.NumCPU()
// 	numTemplates := len(e.DeepDataMerge.Templates)
// 	concurrency := numCPU * 2 // Adjust the concurrency factor based on system hardware resources

// 	if numTemplates < concurrency {
// 		concurrency = numTemplates
// 	}

// 	if concurrency < 1 {
// 		concurrency = 1
// 	}
// 	// templateURLs := make([]string, 0, numTemplates)
// 	// for templateURL := range e.DeepDataMerge.Templates {
// 	// 	templateURLs = append(templateURLs, string(templateURL))
// 	// }

// 	templateURLs := make([]string, 0, numTemplates)
// 	for templateURL := range e.DeepDataMerge.Templates {
// 		s := string(templateURL)
// 		if s == ".html" {
// 			continue
// 		}
// 		templateURLs = append(templateURLs, s)
// 	}

// 	// var wg sync.WaitGroup
// 	// semaphore := make(chan struct{}, concurrency)

// 	tasks := make(chan string, concurrency)
// 	var wg sync.WaitGroup

// 	for i := 0; i < concurrency; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for templateURL := range tasks {
// 				URL := template.URL(templateURL)
// 				layout := e.DeepDataMerge.Templates[URL].Frontmatter.Layout
// 				e.RenderPage(fileOutPath, URL, templates, layout)
// 			}
// 		}()
// 	}
// 	for _, v := range templateURLs {
// 		tasks <- v
// 	}
// 	close(tasks)

// 	// for _, templateURL := range templateURLs {
// 	// 	if templateURL == ".html" {
// 	// 		continue
// 	// 	}

// 	// 	wg.Add(1)
// 	// 	semaphore <- struct{}{}

// 	// 	go func(templateURL string) {
// 	// 		defer func() {
// 	// 			<-semaphore
// 	// 			wg.Done()
// 	// 		}()

// 	// 		e.RenderPage(fileOutPath, template.URL(templateURL), templates, e.DeepDataMerge.Templates[template.URL(templateURL)].Frontmatter.Layout)
// 	// 	}(templateURL)
// 	// }

// 	wg.Wait()
// }

func (e *Engine) RenderUserDefinedPages(fileOutPath string, templates *template.Template) {
	numCPU := runtime.NumCPU()
	numTemplates := len(e.DeepDataMerge.Templates)
	if numTemplates == 0 {
		return
	}

	templateURLs := make([]template.URL, 0, numTemplates)
	for u := range e.DeepDataMerge.Templates {
		if string(u) == ".html" {
			continue
		}
		templateURLs = append(templateURLs, u)
	}

	renderers := max(1, numCPU)
	writers := max(1, numCPU)

	type workItem struct {
		url template.URL
		buf *bytes.Buffer
	}

	tasks := make(chan template.URL, renderers*2)
	writeCh := make(chan workItem, writers*2)

	var wgRender sync.WaitGroup
	var wgWrite sync.WaitGroup

	//writer pool
	for i := 0; i < writers; i++ {
		wgWrite.Add(1)
		go func() {
			defer wgWrite.Done()
			for w := range writeCh {
				outPath := filepath.Join(fileOutPath, "rendered", string(w.url))
				if err := os.MkdirAll(filepath.Dir(outPath), 0750); err != nil {
					e.ErrorLogger.Fatal(err)
				}
				f, err := os.Create(outPath)
				if err != nil {
					e.ErrorLogger.Fatal(err)
				}
				if _, err := w.buf.WriteTo(f); err != nil {
					_ = f.Close()
					e.ErrorLogger.Fatal(err)
				}
				if err := f.Close(); err != nil {
					e.ErrorLogger.Fatal(err)
				}
				renderPool.Put(w.buf)
			}
		}()
	}

	//render pool
	for i := 0; i < renderers; i++ {
		wgRender.Add(1)
		go func() {
			defer wgRender.Done()
			for url := range tasks {
				buf := renderPool.Get().(*bytes.Buffer)
				buf.Reset()

				PageData := PageData{
					DeepDataMerge: e.DeepDataMerge,
					PageURL:       url,
				}

				layout := e.DeepDataMerge.Templates[url].Frontmatter.Layout
				if err := templates.ExecuteTemplate(buf, layout, PageData); err != nil {
					e.ErrorLogger.Println("Error at path: ", url)
					e.ErrorLogger.Fatal(err)
				}

				writeCh <- workItem{url: url, buf: buf}
			}
		}()
	}
	for _, url := range templateURLs {
		tasks <- url
	}
	close(tasks)

	wgRender.Wait()
	close(writeCh)
	wgWrite.Wait()
}
