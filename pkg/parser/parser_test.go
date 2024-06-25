package parser_test

import (
	"html/template"
	"log"
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/anna-ssg/anna/v2/pkg/parser"
)

const TestDirPath = "../../test/parser/"

func TestAddFileAndRender(t *testing.T) {
	gotParser := parser.Parser{
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[template.URL][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	wantLayout := parser.LayoutConfig{
		Navbar: []map[string]string{
			{"Index": "index.html"},
			{"Docs": "docs.html"},
			{"Tags": "tags.html"},
			{"Posts": "posts.html"},
		},
		BaseURL:   "example.org",
		SiteTitle: "ssg",
		Author:    "Anna",
	}
	gotParser.LayoutConfig = wantLayout
	t.Run("parsing data out of one markdown post", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}
		wantParser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: gotParser.ErrorLogger,
		}
		sampleFrontmatter, _, markdownContent, parseSuccess := gotParser.ParseMarkdownContent(string(inputMd), "sample_test_path")
		sampleBody := "sample_body"
		if !parseSuccess {
			return
		}

		filename := "testpost.md"
		fileURL := "testpost.html"
		wantParser.MdFilesName = append(wantParser.MdFilesName, filename)
		wantParser.MdFilesPath = append(wantParser.MdFilesPath, filename)
		wantPage := parser.TemplateData{
			CompleteURL: template.URL(fileURL),
			Date:        wantParser.DateParse(sampleFrontmatter.Date).Unix(),
			Frontmatter: sampleFrontmatter,
			Body:        template.HTML(sampleBody),
			// Layout:      want_layout,
		}
		wantParser.LayoutConfig = wantLayout

		wantParser.Templates["testpost.html"] = wantPage
		for _, tag := range sampleFrontmatter.Tags {
			wantParser.TagsMap[template.URL(tag)] = append(wantParser.TagsMap[template.URL(tag)], wantPage)
		}

		if sampleFrontmatter.Type == "post" {
			wantParser.Posts = append(wantParser.Posts, wantPage)
		}

		gotParser.AddFile("", filename, sampleFrontmatter, markdownContent, sampleBody)

		if !reflect.DeepEqual(gotParser, wantParser) {
			t.Errorf("want %v; \ngot %v", wantParser, gotParser)
			// t.Errorf("please see the files yourself")
		}
	})
}

func TestParseMarkdownContent(t *testing.T) {
	p := parser.Parser{
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[template.URL][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	t.Run("render markdown files to html", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		_, bodyGot, _, parseSuccess := p.ParseMarkdownContent(string(inputMd), "sample_test_path")

		if parseSuccess {

			bodyWant, err := os.ReadFile(TestDirPath + "parse_md/html_want_output.html")
			if err = os.WriteFile(TestDirPath+"parse_md/html_got_output.html", []byte(bodyGot), 0666); err != nil {
				t.Errorf("%v", err)
			}
			if string(bodyWant) != bodyGot {
				t.Errorf("%v\nThe expected and generated html can be found in test/", err)
			}
		}
	})

	t.Run("parse frontmatter from markdown files", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		frontmatterGot, _, _, parseSuccess := p.ParseMarkdownContent(string(inputMd), "sample_test_path")

		if !parseSuccess {
			frontmatterWant := parser.Frontmatter{
				Title:       "Markdown Test",
				Date:        "2024-03-23",
				Description: "File containing markdown to test the SSG",
				Type:        "post",
			}

			if !reflect.DeepEqual(frontmatterGot, frontmatterWant) {
				t.Errorf("got %v, \nwant %v", frontmatterGot, frontmatterWant)
			}
		}
	})
}

func TestParseConfig(t *testing.T) {
	t.Run("unmarshal `config.yml` to LayoutConfig", func(t *testing.T) {
		gotParser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}

		wantLayout := parser.LayoutConfig{
			Navbar: []map[string]string{
				{"Index": "index.html"},
				{"Docs": "docs.html"},
				{"Tags": "tags.html"},
				{"Posts": "posts.html"},
			},

			BaseURL:   "example.org",
			SiteTitle: "ssg",
			Author:    "Anna",
		}

		gotParser.ParseConfig(TestDirPath + "layout/config.yml")

		if !reflect.DeepEqual(gotParser.LayoutConfig, wantLayout) {
			t.Errorf("got \n%v want \n%v", gotParser.LayoutConfig, wantLayout)
		}
	})
}

func TestParseRobots(t *testing.T) {
	t.Run("parse and render `robots.txt`", func(t *testing.T) {
		testParser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		testParser.LayoutConfig.BaseURL = "example.org"

		testParser.ParseRobots(TestDirPath+"layout/robots_txt/robots.txt", TestDirPath+"layout/robots_txt/got_robots.txt")

		gotRobotsTxt, err := os.ReadFile(TestDirPath + "layout/robots_txt/got_robots.txt")
		if err != nil {
			t.Errorf("%v", err)
		}
		wantRobotsTxt, err := os.ReadFile(TestDirPath + "layout/robots_txt/want_robots.txt")
		if err != nil {
			t.Errorf("%v", err)
		}
		if !slices.Equal(gotRobotsTxt, wantRobotsTxt) {
			t.Errorf("The expected and generated robots.txt can be found in test/layout/robots_txt/")
		}
	})
}

func TestNotesAndBacklinkParsing(t *testing.T) {

	gotParser := parser.Parser{
		Notes:       make(map[template.URL]parser.Note),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// creating dummy notes for testing

	gotParser.Notes["notes/test/test.md"] = parser.Note{
		CompleteURL: "notes/test/test.md",
		Frontmatter: parser.Frontmatter{
			Title: "head note",
			Type:  "note",
			Head:  true,
		},
		Body: "This is a [[backlink]] here",
	}

	gotParser.Notes["notes/test/backlink.md"] = parser.Note{
		CompleteURL: "notes/test/backlink.md",
		Frontmatter: parser.Frontmatter{
			Title: "backlink",
			Type:  "note",
		},
		Body: "Content of note.",
	}

	t.Run("testing notes and backlink parsing", func(t *testing.T) {

		gotParser.ParseBacklink("notes/test/test.md")

		wantParser := parser.Parser{
			Notes: map[template.URL]parser.Note{
				"notes/test/test.md": {
					CompleteURL: template.URL("notes/test/test.md"),
					Frontmatter: parser.Frontmatter{
						Title: "head note",
						Type:  "note",
						Head:  true,
					},
					Body:           template.HTML("This is a <a id=\"zettel-reference\" href=\"/notes/test/backlink.md\">backlink</a> here"),
					LinkedNoteURLs: []template.URL{"notes/test/backlink.md"},
				},

				"notes/test/backlink.md": {
					CompleteURL: template.URL("notes/test/backlink.md"),
					Frontmatter: parser.Frontmatter{
						Title: "backlink",
						Type:  "note",
					},
					Body: template.HTML("Content of note."),
				},
			},
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}

		if !reflect.DeepEqual(gotParser.Notes, wantParser.Notes) {
			t.Errorf("got %v,\n want %v", gotParser.Notes, wantParser.Notes)
		}
	})
}
