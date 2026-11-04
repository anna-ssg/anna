package parser_test

import (
	"html/template"
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/anna-ssg/anna/pkg/logger"
	"github.com/anna-ssg/anna/pkg/parser"
)

const TestDirPath = "../../test/parser/"

func TestAddFileAndRender(t *testing.T) {
	gotParser := parser.Parser{
		Pages:                     make(map[template.URL]parser.PageData),
		TagsMap:                   make(map[template.URL][]parser.PageData),
		CollectionsMap:            make(map[template.URL][]parser.PageData),
		CollectionsHierarchy:      make(map[template.URL]*parser.CollectionHierarchyData),
		CollectionsSubPageLayouts: make(map[template.URL]string),
		ErrorLogger:               logger.New(os.Stderr),
	}

	wantLayout := parser.Config{
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
	gotParser.Config = wantLayout
	t.Run("parsing frontmatter and body from a markdown file", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}
		wantParser := parser.Parser{
			Pages:                     make(map[template.URL]parser.PageData),
			TagsMap:                   make(map[template.URL][]parser.PageData),
			CollectionsMap:            make(map[template.URL][]parser.PageData),
			CollectionsHierarchy:      make(map[template.URL]*parser.CollectionHierarchyData),
			CollectionsSubPageLayouts: make(map[template.URL]string),
			ErrorLogger:               gotParser.ErrorLogger,
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
		wantPage := parser.PageData{
			CompleteURL: template.URL(fileURL),
			Date:        wantParser.DateParse(sampleFrontmatter.Date).Unix(),
			Frontmatter: sampleFrontmatter,
			Body:        template.HTML(sampleBody),
		}
		wantParser.Config = wantLayout
		wantParser.CollectionsMap[template.URL("index.html")] = []parser.PageData{
			wantPage,
		}

		wantParser.Pages["testpost.html"] = wantPage
		for _, tag := range sampleFrontmatter.Tags {
			wantParser.TagsMap[template.URL(tag)] = append(wantParser.TagsMap[template.URL(tag)], wantPage)
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
		Pages:                make(map[template.URL]parser.PageData),
		TagsMap:              make(map[template.URL][]parser.PageData),
		CollectionsMap:       make(map[template.URL][]parser.PageData),
		CollectionsHierarchy: make(map[template.URL]*parser.CollectionHierarchyData),
		ErrorLogger:          logger.New(os.Stderr),
	}
	t.Run("render markdown files to html", func(t *testing.T) {
		inputMd, err := os.ReadFile(TestDirPath + "parse_md/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		_, bodyGot, _, parseSuccess := p.ParseMarkdownContent(string(inputMd), "sample_test_path")

		if parseSuccess {

			bodyWant, err := os.ReadFile(TestDirPath + "parse_md/html_want_output.html")
			if err != nil {
				t.Errorf("%v", err)
			}
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
			}

			if !reflect.DeepEqual(frontmatterGot, frontmatterWant) {
				t.Errorf("got %v, \nwant %v", frontmatterGot, frontmatterWant)
			}
		}
	})
}

func TestParseConfig(t *testing.T) {
	t.Run("unmarshal `config.json` to Config", func(t *testing.T) {
		gotParser := parser.Parser{
			Pages:                     make(map[template.URL]parser.PageData),
			TagsMap:                   make(map[template.URL][]parser.PageData),
			CollectionsMap:            make(map[template.URL][]parser.PageData),
			CollectionsHierarchy:      make(map[template.URL]*parser.CollectionHierarchyData),
			CollectionsSubPageLayouts: make(map[template.URL]string),
			ErrorLogger:               logger.New(os.Stderr),
		}

		wantLayout := parser.Config{
			Navbar: []map[string]string{
				{"Index": "index.html"},
				{"Docs": "docs.html"},
				{"Tags": "tags.html"},
				{"Posts": "posts.html"},
			},

			BaseURL:                      "example.org",
			SiteTitle:                    "ssg",
			Author:                       "Anna",
			DefaultPageLayout:            "page-test",
			DefaultCollectionIndexLayout: "collection-index-test",
			CollectionLayouts: map[string]string{
				"posts/index-test.html": "all-posts-test",
			},
		}

		gotParser.ParseConfig(TestDirPath + "layout/config.json")

		if !reflect.DeepEqual(gotParser.Config, wantLayout) {
			t.Errorf("got \n%v want \n%v", gotParser.Config, wantLayout)
		}
	})
}

func TestParseRobots(t *testing.T) {
	t.Run("parse and render `robots.txt`", func(t *testing.T) {
		testParser := parser.Parser{
			Pages:                make(map[template.URL]parser.PageData),
			TagsMap:              make(map[template.URL][]parser.PageData),
			CollectionsMap:       make(map[template.URL][]parser.PageData),
			CollectionsHierarchy: make(map[template.URL]*parser.CollectionHierarchyData),
			ErrorLogger:          logger.New(os.Stderr),
		}
		testParser.Config.BaseURL = "example.org"

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
