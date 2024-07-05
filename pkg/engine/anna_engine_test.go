package engine_test

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/anna-ssg/anna/v3/pkg/engine"
	"github.com/anna-ssg/anna/v3/pkg/parser"
)

func TestRenderTags(t *testing.T) {
	e := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
	e.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)
	e.DeepDataMerge.LayoutConfig.BaseURL = "example.org"

	fileOutPath := "../../test/engine/render_tags/"

	e.DeepDataMerge.TagsMap["tags/blogs.html"] = []parser.TemplateData{
		{
			CompleteURL: "posts/file1.html",
			Frontmatter: parser.Frontmatter{
				Title: "file1",
				Tags:  []string{"blogs"},
			},
		},
		{
			CompleteURL: "posts/file2.html",
			Frontmatter: parser.Frontmatter{
				Title: "file2",
				Tags:  []string{"blogs", "tech"},
			},
		},
	}

	e.DeepDataMerge.TagsMap["tags/tech.html"] = []parser.TemplateData{
		{
			CompleteURL: "posts/file2.html",
			Frontmatter: parser.Frontmatter{
				Title: "file2",
				Tags:  []string{"blogs", "tech"},
			},
		},
		{
			CompleteURL: "posts/file3.html",
			Frontmatter: parser.Frontmatter{
				Title: "file3",
				Tags:  []string{"tech"},
			},
		},
	}

	templ, err := template.ParseFiles(TestDirPath+"render_tags/tags_template.html", TestDirPath+"render_tags/tags_subpage_template.html")
	if err != nil {
		t.Errorf("%v", err)
	}
	if err := os.MkdirAll(TestDirPath+"render_tags/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}
	e.RenderTags(fileOutPath, templ)

	t.Run("render tag.html", func(t *testing.T) {
		gotTagsFile, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantTagsFile, err := os.ReadFile(TestDirPath + "render_tags/want_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(gotTagsFile, wantTagsFile) {
			t.Errorf("The expected and generated tags.html can be found in test/engine/render_tags/rendered/")
		}
	})

	t.Run("render tag-subpage.html", func(t *testing.T) {
		gotBlogsFile, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags/blogs.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantBlogsFile, err := os.ReadFile(TestDirPath + "render_tags/want_blogs_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(gotBlogsFile, wantBlogsFile) {
			t.Errorf("The expected and generated blogs.html tag-subpage can be found in test/engine/render_tags/rendered/tags/")
		}

		gotTechFile, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags/tech.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantTechFile, err := os.ReadFile(TestDirPath + "render_tags/want_tech_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(gotTechFile, wantTechFile) {
			t.Errorf("The expected and generated tech.html tag-subpage can be found in test/engine/render_tags/rendered/tags/")
		}
	})
}

func TestGenerateMergedJson(t *testing.T) {
	if err := os.MkdirAll(TestDirPath+"json_index_test/rendered/static", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("test json creation for the search index", func(t *testing.T) {
		e := engine.Engine{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
		e.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)

		e.DeepDataMerge.Templates["docs.md"] = parser.TemplateData{
			CompleteURL: "docs.html",
			Frontmatter: parser.Frontmatter{
				Title: "Anna Documentation",
			},
		}

		e.GenerateJSONIndex(TestDirPath + "json_index_test/")

		gotJson, err := os.ReadFile(TestDirPath + "/json_index_test/rendered/static/index.json")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantJson, err := os.ReadFile(TestDirPath + "/json_index_test/want_index.json")
		if err != nil {
			t.Errorf("%v", err)
		}

		gotJson = bytes.TrimSpace(gotJson)
		wantJson = bytes.TrimSpace(wantJson)

		if !slices.Equal(gotJson, wantJson) {
			t.Errorf("The expected and generated json can be found in test/engine/json_index_test")
		}
	})
}

func TestGenerateSitemap(t *testing.T) {
	t.Run("render sitemap.xml", func(t *testing.T) {
		testEngine := engine.Engine{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		testEngine.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
		testEngine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)

		t1 := parser.TemplateData{
			CompleteURL: "index.html",
			Frontmatter: parser.Frontmatter{
				Date: "2024-02-23",
			},
		}

		t2 := parser.TemplateData{
			CompleteURL: "about.html",
			Frontmatter: parser.Frontmatter{
				Date: "2024-02-23",
			},
		}

		t3 := parser.TemplateData{
			CompleteURL: "research.html",
			Frontmatter: parser.Frontmatter{
				Date: "2024-02-23",
			},
		}

		testEngine.DeepDataMerge.LayoutConfig.BaseURL = "example.org"
		// setting up testEngine
		testEngine.DeepDataMerge.Templates["index"] = t1
		testEngine.DeepDataMerge.Templates["about"] = t2
		testEngine.DeepDataMerge.Templates["research"] = t3

		testEngine.GenerateSitemap(TestDirPath + "sitemap/got_sitemap.xml")

		gotSitemap, err := os.ReadFile(TestDirPath + "sitemap/got_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of got_sitemap.xml")
		}

		wantSitemap, err := os.ReadFile(TestDirPath + "sitemap/want_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of _sitemap.xml")
		}

		gotSitemapString := string(gotSitemap)
		wantSitemapString := string(wantSitemap)
		gotSitemapString = strings.TrimFunc(gotSitemapString, func(r rune) bool {
			return r == '\n' || r == '\t' || r == ' '
		})
		wantSitemapString = strings.TrimFunc(wantSitemapString, func(r rune) bool {
			return r == '\n' || r == '\t' || r == ' '
		})

		if strings.Compare(gotSitemapString, wantSitemapString) == 0 {
			t.Errorf("The expected and generated sitemap can be found in test/testEngine/sitemap/")
		}
	})
}
