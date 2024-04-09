package zettel_parser

import (
	"bytes"
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
	LayoutConfig             parser.LayoutConfig
}

type Parser struct {
	// Holds the data of all of the notes
	NotesMergedData NotesMerged

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

					fronmatter, body, parseSuccess := p.ParseNoteMarkdownContent(string(content))
					if parseSuccess {
						// ISSUE
						p.AddNote(baseDirPath, fileName, fronmatter, body)
					}
				}
			}
		}
		return nil
	})
}

func (p *Parser) ParseNoteMarkdownContent(filecontent string) (Frontmatter, string, bool) {
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
		return Frontmatter{}, "", false
	}

	regex := regexp.MustCompile(`title(.*): (.*)`)
	match := regex.FindStringSubmatch(splitContents[1])

	if match == nil {
		return Frontmatter{}, "", false
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

	return parsedFrontmatter, parsedMarkdown.String(), true
}

func (p *Parser) AddNote(baseDirPath string, dirEntryPath string, frontmatter Frontmatter, body string) {
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
		LinkedNoteTitles:         []string{},
		// Layout:                   p.LayoutConfig,
	}

	p.NotesMergedData.Notes[template.URL(key)] = note

	/*
		REGEX:
		using regex we need to identify for the references to other
		notes of the following pattern

		[Note Title Name](/note/somefilename)
	*/

	re := regexp.MustCompile(`\[.*\]\((\/notes\/).*(.md)\)`)
	re_sub := regexp.MustCompile(`\[.*\]`)
	matches := re.FindAllString(string(note.Body), -1)

	for _, match := range matches {
		/*
			Extracting the file "Titles" from the first match
			ex: [Nunc ullamcorper](/notes/2021-09-01-nunc-ullamcorper.md)
			will extract out "[Nunc ullamcorper]"
		*/
		sub_match := re_sub.FindString(match)
		sub_match = strings.Trim(sub_match, "[]")

		note.LinkedNoteTitles = append(note.LinkedNoteTitles, sub_match)
	}
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



}
