package zettel_engine

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	zettel_parser "github.com/acmpesuecc/anna/pkg/zettel/parser"
)

type Engine struct {

	// Holds the data of all of the notes
	NotesMergedData zettel_parser.NotesMerged

	// Common logger for all engine functions
	ErrorLogger *log.Logger
}

func (e *Engine) RenderNote(fileOutPath string, pagePath template.URL, templ *template.Template, noteURL template.URL) {
	// Creating subdirectories if the filepath contains '/'
	dirPath := ""
	if strings.Contains(string(pagePath), "/") {
		// Extracting the directory path from the page path
		dirPath, _ := strings.CutSuffix(string(pagePath), e.NotesMergedData.Notes[noteURL].FilenameWithoutExtension)
		dirPath = fileOutPath + "rendered/" + dirPath

		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}

	filename, _ := strings.CutSuffix(string(pagePath), ".md")
	filepath := fileOutPath + "rendered/" + dirPath + filename + ".html"
	var buffer bytes.Buffer

	// Storing the rendered HTML file to a buffer
	err := templ.ExecuteTemplate(&buffer, "note", e.NotesMergedData.Notes[noteURL])
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) RenderUserNotes(fileOutPath string, templ *template.Template) {
	// Loop and render user notes
	for _, Note := range e.NotesMergedData.Notes {

		//htmlFilePath, _ := strings.CutSuffix(string(noteURL), ".md")
		//destinationPath := fileOutPath + "render/" + htmlFilePath + ".html"
		fileInPath := strings.TrimSuffix(string(Note.CompleteURL), ".html")

		e.RenderNote(fileOutPath, template.URL(fileInPath), templ, Note.CompleteURL)

	}
}

func (e *Engine) RetrieveNotePointer(noteTitle string) *zettel_parser.Note {
	for _, Note := range e.NotesMergedData.Notes {
		if Note.Frontmatter.Title == noteTitle {
			return &Note
		}
	}
	return nil
}

func (e *Engine) GenerateLinkStore() {
	// Populate the LinkStore map
	for _, Note := range e.NotesMergedData.Notes {
		for _, referencedNoteTitle := range Note.LinkedNoteTitles {
			referencedNotePointer := e.RetrieveNotePointer(referencedNoteTitle)
			if referencedNotePointer == nil {
				e.ErrorLogger.Fatalf("ERR: Failed to get pointer to note %s\n", referencedNoteTitle)
			}
			e.NotesMergedData.LinkStore[Note.CompleteURL] = append(
				e.NotesMergedData.LinkStore[Note.CompleteURL],
				referencedNotePointer,
			)
		}
	}
}

func (e *Engine) GenerateRootNote(fileOutPath string, templ *template.Template) {
	// This is the page that acts as the root of all the
	// notes part of the site

	// Creating a map of all head notes

	var buffer bytes.Buffer

	fmt.Println(e.NotesMergedData.LinkStore)

	/*
		t := template.Must(templ.Funcs(template.FuncMap{
			"Deref": func(i *zettel_parser.Note) zettel_parser.Note { return *note },
		}).Parse(src))
	*/

	err := templ.ExecuteTemplate(&buffer, "root", e.NotesMergedData.LinkStore)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	err = os.WriteFile(fileOutPath+"rendered/notes.html", buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
