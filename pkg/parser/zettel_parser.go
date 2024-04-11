package parser

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

type Note struct {
	CompleteURL    template.URL
	Date           int64
	Frontmatter    Frontmatter
	Body           template.HTML
	LinkedNoteURLs []template.URL
}

func (p *Parser) BackLinkParser() {
	/*
		This function is going to validate whether all the
		references in the notes have a valid link to another
		note

		Example: for `[Nunc ullamcorper]]` to be
		a valid reference, the title part of the frontmatter
		of the note `/note/1234.md` must have "Nunc ullamcorper"
	*/

	for noteURL, note := range p.Notes {

		noteBody := string(note.Body) // template.HTML -> string
		// noteParentDir := note.CompleteURL

		// fmt.Println("Finding links for :", noteParentDir)

		backlinkRE := regexp.MustCompile(`\[[^\]]*\]\]`)
		backlinks := backlinkRE.FindAllString(noteBody, -1)

		for _, backlink := range backlinks {
			// Now that we have the backlinks to titles,
			// we need to walk the notes dir to find if there
			// are any matches
			noteTitle := strings.Trim(backlink, "[]")

			referenceCompleteURL, err := p.ValidateBackLink(noteTitle)
			if err != nil{
				p.ErrorLogger.Fatal(err)
			} else {
				// creating anchor tag reference for parsed markdown
				anchorReference := fmt.Sprintf(`<a href="/%s">%s</a>`, referenceCompleteURL, noteTitle)
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
			LinkedNoteURLs: note.LinkedNoteURLs,
		}

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
