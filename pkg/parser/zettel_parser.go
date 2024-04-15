package parser

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type Note struct {
	CompleteURL    template.URL
	Date           int64
	Frontmatter    Frontmatter
	Body           template.HTML
	MarkdownBody   string
	LinkedNoteURLs []template.URL
	LiveReload     bool
}

var backlinkRE = regexp.MustCompile(`\[[^\]]*\]\]`)

// TODO: The current regex will search for all types of callouts of
// [[]] in the body of the markdown. The disadvantage is that it will
// not be able to ignore the callouts mentioned inside code blocks
// on inline code blocks.
//
// Change the current method with which these are parsed such that
// this edge cases is handled correctly.

func (p *Parser) BackLinkParser() {
	/*
		This function is going to validate whether all the
		references in the notes have a valid link to another
		note

		Example: for `[Nunc ullamcorper]]` to be
		a valid reference, the title part of the frontmatter
		of the note `/note/1234.md` must have "Nunc ullamcorper"
	*/

	numCPU := runtime.NumCPU()
	numNotes := len(p.Notes)
	concurrency := numCPU * 2

	if numNotes < concurrency {
		concurrency = numNotes
	}

	noteURLS := make([]string, 0, numNotes)
	for noteURL := range p.Notes {
		noteURLS = append(noteURLS, string(noteURL))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, url := range noteURLS {
		if url == ".html" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(noteURL string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			p.ParseBacklink(template.URL(noteURL))
		}(url)

		wg.Wait()
	}
}

func (p *Parser) ParseBacklink(noteURL template.URL) {
	note := p.Notes[noteURL]
	noteBody := string(note.Body) // template.HTML -> string
	// noteParentDir := note.CompleteURL

	// fmt.Println("Finding links for :", noteParentDir)

	backlinks := backlinkRE.FindAllString(noteBody, -1)

	for _, backlink := range backlinks {
		// Now that we have the backlinks to titles,
		// we need to walk the notes dir to find if there
		// are any matches
		noteTitle := strings.Trim(backlink, "[]")

		referenceCompleteURL, err := p.ValidateBackLink(noteTitle)
		if err != nil {
			p.ErrorLogger.Fatal(err)
		} else {
			// creating anchor tag reference for parsed markdown
			anchorReference := fmt.Sprintf(`<a id="zettel-reference" href="/%s">%s</a>`, referenceCompleteURL, noteTitle)
			noteBody = strings.ReplaceAll(noteBody, backlink, anchorReference)

			// fmt.Println(note.LinkedNoteURLs)
			note.LinkedNoteURLs = append(note.LinkedNoteURLs, referenceCompleteURL)
		}
	}

	p.Notes[noteURL] = Note{
		CompleteURL:    note.CompleteURL,
		Date:           note.Date,
		Frontmatter:    note.Frontmatter,
		Body:           template.HTML(noteBody),
		MarkdownBody:   note.MarkdownBody,
		LinkedNoteURLs: note.LinkedNoteURLs,
	}
}

func (p *Parser) ValidateBackLink(noteTitle string) (template.URL, error) {
	for _, note := range p.Notes {
		if note.Frontmatter.Title == noteTitle {
			return note.CompleteURL, nil
		}
	}

	errorMessage := fmt.Sprintf("ERR: Failed to find a note for backlink %s\n", noteTitle)
	return "", errors.New(errorMessage)
}
