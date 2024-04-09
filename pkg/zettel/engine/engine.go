package zettel_engine

import (
	"bytes"
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
	err := templ.ExecuteTemplate(&buffer, "note", e.NotesMergedData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) RenderUserNotes(){
	// Loop and render user notes
}

func (e *Engine) GenerateLinkStore(){
	// Populate the LinkStore map
}
