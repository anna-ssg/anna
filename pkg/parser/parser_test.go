package parser_test

import (
	"html/template"
	"log"
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/acmpesuecc/anna/pkg/parser"
)

const TestDirPath = "../../test/parser/"

func TestAddFileandRender(t *testing.T) {
	got_parser := parser.Parser{
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[template.URL][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	want_layout := parser.LayoutConfig{
		Navbar:    []string{"index", "docs", "tags", "posts"},
		BaseURL:   "example.org",
		SiteTitle: "ssg",
		Author:    "Anna",
	}
	got_parser.LayoutConfig = want_layout
	t.Run("parsing data out of one markdown post", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}
		want_parser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: got_parser.ErrorLogger,
		}
		sample_frontmatter, _, parseSuccess := got_parser.ParseMarkdownContent(string(inputMd))
		sample_body := "sample_body"
		if !parseSuccess {
			return
		}

		filename := "testpost.md"
		fileurl := "posts/testpost.html"
		want_parser.MdFilesName = append(want_parser.MdFilesName, filename)
		want_parser.MdFilesPath = append(want_parser.MdFilesPath, filename)
		want_page := parser.TemplateData{
			CompleteURL: template.URL(fileurl),
			Date:        want_parser.DateParse(sample_frontmatter.Date).Unix(),
			Frontmatter: sample_frontmatter,
			Body:        template.HTML(sample_body),
			// Layout:      want_layout,
		}
		want_parser.LayoutConfig = want_layout

		want_parser.Templates[template.URL("posts/testpost.html")] = want_page
		for _, tag := range sample_frontmatter.Tags {
			want_parser.TagsMap[template.URL(tag)] = append(want_parser.TagsMap[template.URL(tag)], want_page)
		}

		if sample_frontmatter.Type == "post" {
			want_parser.Posts = append(want_parser.Posts, want_page)
		}

		got_parser.AddFileAndRender("", filename, sample_frontmatter, sample_body)

		if !reflect.DeepEqual(got_parser, want_parser) {
			t.Errorf("want %v; \ngot %v", want_parser, got_parser)
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

		_, body_got, parseSuccess := p.ParseMarkdownContent(string(inputMd))

		if parseSuccess {

			body_want, err := os.ReadFile(TestDirPath + "parse_md/html_want_output.html")
			if err = os.WriteFile(TestDirPath+"parse_md/html_got_output.html", []byte(body_got), 0666); err != nil {
				t.Errorf("%v", err)
			}
			if string(body_want) != body_got {
				t.Errorf("%v\nThe expected and generated html can be found in test/", err)
			}
		}
	})

	t.Run("parse frontmatter from markdown files", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		frontmatter_got, _, parseSuccess := p.ParseMarkdownContent(string(inputMd))

		if !parseSuccess {
			frontmatter_want := parser.Frontmatter{
				Title:       "Markdown Test",
				Date:        "2024-03-23",
				Description: "File containing markdown to test the SSG",
				Type:        "post",
			}

			if !reflect.DeepEqual(frontmatter_got, frontmatter_want) {
				t.Errorf("got %v, \nwant %v", frontmatter_got, frontmatter_want)
			}
		}
	})
}

func TestParseConfig(t *testing.T) {
	t.Run("unmarshal `config.yml` to LayoutConfig", func(t *testing.T) {
		got_parser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}

		want_layout := parser.LayoutConfig{
			Navbar:    []string{"index", "docs", "tags", "posts"},
			BaseURL:   "example.org",
			SiteTitle: "ssg",
			Author:    "Anna",
		}

		got_parser.ParseConfig(TestDirPath + "layout/config.yml")

		if !reflect.DeepEqual(got_parser.LayoutConfig, want_layout) {
			t.Errorf("got \n%v want \n%v", got_parser.LayoutConfig, want_layout)
		}
	})
}

func TestParseRobots(t *testing.T) {
	t.Run("parse and render `robots.txt`", func(t *testing.T) {
		parser := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		parser.LayoutConfig.BaseURL = "example.org"

		parser.ParseRobots(TestDirPath+"layout/robots_txt/robots.txt", TestDirPath+"layout/robots_txt/got_robots.txt")

		got_robots_txt, err := os.ReadFile(TestDirPath + "layout/robots_txt/got_robots.txt")
		if err != nil {
			t.Errorf("%v", err)
		}
		want_robots_txt, err := os.ReadFile(TestDirPath + "layout/robots_txt/want_robots.txt")
		if err != nil {
			t.Errorf("%v", err)
		}
		if !slices.Equal(got_robots_txt, want_robots_txt) {
			t.Errorf("The expected and generated robots.txt can be found in test/layout/robots_txt/")
		}
	})
}
