package engine

//	type postsTemplateData struct {
//		Posts []parser.TemplateData
//		parser.TemplateData
//	}

/*
func (e *Engine) RenderSite(addr string) {
	// Creating the "rendered" directory if not present
	err := os.RemoveAll(helpers.SiteDataPath + "rendered/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(helpers.SiteDataPath+"rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	var wg sync.WaitGroup
	concurrency := 3
	semaphore := make(chan struct{}, concurrency) // Each goroutine handles 3 files at a time

	files := make([]string, 0, len(e.Templates))
	for pagePath := range g.Templates {
		files = append(files, string(pagePath))
	}

	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(file string) {
			defer func() {
				<-semaphore // Release semaphore
				wg.Done()
			}()

			pagePath := template.URL(file)
			templateData := e.Templates[pagePath]
			g.RenderPage(pagePath, templateData, templ, "page")
		}(file)
	}

	wg.Wait()

	var postsBuffer bytes.Buffer

	postsData := postsTemplateData{
		Posts: g.Posts,
		TemplateData: parser.TemplateData{
			Frontmatter: parser.Frontmatter{Title: "Posts"},
			Layout:      g.LayoutConfig,
		},
	}

	err = templ.ExecuteTemplate(&postsBuffer, "posts", postsData)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile(SiteDataPath+"rendered/posts.html", postsBuffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	g.RenderTags(templ)
}

*/

// FUNC: RenderUserDefinedPages

// FUNC: RenderSsgFiles
