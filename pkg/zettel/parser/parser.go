package zettel_parser

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type NotesMerged struct {
	//Stores all the notes
	Notes map[template.URL]Note

	//Stores the links of each note to other notes
	LinkStore map[template.URL][]*Note
}

type Note struct {
	CompleteURL              template.URL
	FilenameWithoutExtension string
	Date                     int64
	Frontmatter              Frontmatter
	Body                     template.HTML
	LinkedNoteTitles         []string
	Layout                   parser.LayoutConfig
}

type Parser struct {
	// Holds the data of all of the notes
	NotesMergedData NotesMerged

	Layout parser.LayoutConfig

	// Common logger for all parser functions
	ErrorLogger *log.Logger
}

type Frontmatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	JSFiles      []string `yaml:"scripts"`
	Type         string   `yaml:"type"`
	Description  string   `yaml:"description"`
	PreviewImage string   `yaml:"previewimage"`
	Head         bool     `yaml:"head"`
	// Tags         []string `yaml:"tags"`
}

func (p *Parser) ParseNotesDir(baseDirPath string, baseDirFS fs.FS) {
	fs.WalkDir(baseDirFS, ".", func(path string, dir fs.DirEntry, err error) error {
		if path != "." && path != ".obsidian" {
			if dir.IsDir() {
				subDir := os.DirFS(path)
				p.ParseNotesDir(path, subDir)
			} else {
				if filepath.Ext(path) == ".md" {
					fileName := filepath.Base(path)

					content, err := os.ReadFile(baseDirPath + path)
					if err != nil {
						p.ErrorLogger.Fatal(err)
					}

					fronmatter, body, linkedNoteTitles, parseSuccess := p.ParseNoteMarkdownContent(string(content))
					if parseSuccess {
						// ISSUE
						p.AddNote(baseDirPath, fileName, fronmatter, body, linkedNoteTitles)
						// fmt.Println(fileName, linkedNoteTitles)
					}
				}
			}
		}
		return nil
	})
	p.ValidateNoteReferences()
}

func (p *Parser) ParseNoteMarkdownContent(filecontent string) (Frontmatter, string, []string, bool) {
	var parsedFrontmatter Frontmatter
	var markdown string
	/*
	   ---
	   frontmatter_content
	   ---

	   markdown content
	   --- => markdown divider and not to be touched while yaml parsing
	*/
	splitContents := strings.Split(filecontent, "---")
	frontmatterSplit := ""

	if len(splitContents) <= 1 {
		return Frontmatter{}, "", []string{}, false
	}

	regex := regexp.MustCompile(`title(.*): (.*)`)
	match := regex.FindStringSubmatch(splitContents[1])

	if match == nil {
		return Frontmatter{}, "", []string{}, false
	}

	frontmatterSplit = splitContents[1]
	// Parsing YAML frontmatter
	err := yaml.Unmarshal([]byte(frontmatterSplit), &parsedFrontmatter)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	markdown = strings.Join(strings.Split(filecontent, "---")[2:], "---")

	// TODO:
	// This section must replace the callouts with
	// the html url references
	/*
		REGEX:
		using regex we need to identify for the references to other
		notes of the following pattern

		[Note Title Name](/note/somefilename)
	*/

	// DONT DELETE: re := regexp.MustCompile(`\[[^\]]*\]\(/notes/[^\]]*\.html\)`)
	re := regexp.MustCompile(`\[\[[^\]]*\]\]`)

	re_sub := regexp.MustCompile(`\[.*\]`)
	matches := re.FindAllString(markdown, -1)
	fmt.Printf("%s : ", parsedFrontmatter.Title)

	linkedNoteTitles := []string{}
	fmt.Printf("%s\n", matches)

	for _, match := range matches {
		/*
									Extracting the file "Titles" from the first match
									ex: [[Nunc ullamcorper]]
									will extract out "Nunc ullamcorper"

						      We will change the reference to use
						      [[Nunc ullamcorper]] => [Nunc ullamcorper](<Nunc ullamcorper.html>)

			            NOTE: This is temoprary and will have to make it such that
			            it could be present in any file name. Hence this method
			            will have to move to another function
		*/
		sub_match := re_sub.FindString(match)
		sub_match = strings.Trim(sub_match, "[]")
		fmt.Printf("\t%s\n", sub_match)

		linkedNoteTitles = append(linkedNoteTitles, sub_match)

		note_name := strings.Join([]string{sub_match, "html"}, ".")

		// replacing reference with a markdown reference
		new_reference := fmt.Sprintf("[%s](</notes/%s>)", sub_match, note_name)
		markdown = strings.ReplaceAll(markdown, match, new_reference)
		fmt.Printf("%s => %s\n", match, fmt.Sprintf("[%s](</notes/%s>)", sub_match, note_name))
	}

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	if err := md.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		p.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String(), linkedNoteTitles, true
}

func (p *Parser) AddNote(baseDirPath string, dirEntryPath string, frontmatter Frontmatter, body string, linkedNoteTitles []string) {
	filepath := baseDirPath + dirEntryPath

	var date int64
	if frontmatter.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", frontmatter.Date)
		if err != nil {
			p.ErrorLogger.Fatal(err)
		}
		date = parsedDate.Unix()
	}

	key, _ := strings.CutPrefix(filepath, helpers.SiteDataPath+"content/")
	url, _ := strings.CutSuffix(key, ".md")
	url += ".html"
	if frontmatter.Type == "post" {
		url = "posts/" + url
	}

	note := Note{
		CompleteURL:              template.URL(url),
		Date:                     date,
		FilenameWithoutExtension: strings.Split(dirEntryPath, ".")[0],
		Frontmatter:              frontmatter,
		Body:                     template.HTML(body),
		LinkedNoteTitles:         linkedNoteTitles,
		Layout:                   p.Layout,
	}

	p.NotesMergedData.Notes[note.CompleteURL] = note
	//fmt.Println(note.Layout)
}

func (p *Parser) ValidateNoteTitle(ReferenceNoteTitle string) bool {
	for _, Note := range p.NotesMergedData.Notes {
		if Note.Frontmatter.Title == ReferenceNoteTitle {
			return true
		}
	}
	return false
}

func (p *Parser) ValidateNoteReferences() {
	/*
		This function is going to validate whether all the
		references in the notes have a valid link to another
		note

		Example: for `[Nunc ullamcorper](/notes/1234.md)` to be
		a valid reference, the title part of the frontmatter
		of the note `/note/1234.md` must have "Nunc ullamcorper"
	*/
	for _, Note := range p.NotesMergedData.Notes {
		for _, ReferenceNoteTitle := range Note.LinkedNoteTitles {
			if !p.ValidateNoteTitle(ReferenceNoteTitle) {
				p.ErrorLogger.Fatalf("ERR: Referenced note title (%s) doesnt have an existing note", ReferenceNoteTitle)
			}
		}
	}
}
