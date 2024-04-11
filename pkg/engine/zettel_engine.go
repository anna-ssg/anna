package engine

import (
	"html/template"
	"runtime"
	"sync"
)

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

// func (z *Zettel) RetrieveNotePointer(noteTitle string) *zettel_parser.Note {
// 	for _, Note := range e.NotesMergedData.Notes {
// 		if Note.Frontmatter.Title == noteTitle {
// 			return &Note
// 		}
// 	}
// 	return nil
// }

// func (e *Engine) GenerateRootNote(fileOutPath string, templ *template.Template) {
// 	// This is the page that acts as the root of all the
// 	// notes part of the site

// 	// Creating a map of all head notes

// 	var buffer bytes.Buffer

// 	fmt.Println(e.NotesMergedData.LinkStore)

// 	/*
// 		t := template.Must(templ.Funcs(template.FuncMap{
// 			"Deref": func(i *zettel_parser.Note) zettel_parser.Note { return *note },
// 		}).Parse(src))
// 	*/

// 	err := templ.ExecuteTemplate(&buffer, "root", e.NotesMergedData.LinkStore)
// 	if err != nil {
// 		e.ErrorLogger.Fatal(err)
// 	}

// 	err = os.WriteFile(fileOutPath+"rendered/notes.html", buffer.Bytes(), 0666)
// 	if err != nil {
// 		e.ErrorLogger.Fatal(err)
// 	}
// }
