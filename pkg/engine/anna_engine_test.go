package engine_test

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/acmpesuecc/anna/pkg/engine"
	"github.com/acmpesuecc/anna/pkg/parser"
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
		got_tags_file, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_tags_file, err := os.ReadFile(TestDirPath + "render_tags/want_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(got_tags_file, want_tags_file) {
			t.Errorf("The expected and generated tags.html can be found in test/engine/render_tags/rendered/")
		}
	})

	t.Run("render tag-subpage.html", func(t *testing.T) {
		got_blogs_file, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags/blogs.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_blogs_file, err := os.ReadFile(TestDirPath + "render_tags/want_blogs_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(got_blogs_file, want_blogs_file) {
			t.Errorf("The expected and generated blogs.html tag-subpage can be found in test/engine/render_tags/rendered/tags/")
		}

		got_tech_file, err := os.ReadFile(TestDirPath + "render_tags/rendered/tags/tech.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_tech_file, err := os.ReadFile(TestDirPath + "render_tags/want_tech_tags.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(got_tech_file, want_tech_file) {
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

		e.GenerateJSONIndex(TestDirPath + "json_index_test")

		got_json, err := os.ReadFile(TestDirPath + "/json_index_test/rendered/static/index.json")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_json, err := os.ReadFile(TestDirPath + "/json_index_test/want_index.json")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_json = bytes.TrimSpace(got_json)
		want_json = bytes.TrimSpace(want_json)

		if !slices.Equal(got_json, want_json) {
			t.Errorf("The expected and generated json can be found in test/engine/json_index_test")
		}
	})
}

func TestGenerateSitemap(t *testing.T) {
	t.Run("render sitemap.xml", func(t *testing.T) {
		engine := engine.Engine{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		engine.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
		engine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)

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

		engine.DeepDataMerge.LayoutConfig.BaseURL = "example.org"
		// setting up engine
		engine.DeepDataMerge.Templates["index"] = t1
		engine.DeepDataMerge.Templates["about"] = t2
		engine.DeepDataMerge.Templates["research"] = t3

		engine.GenerateSitemap(TestDirPath + "sitemap/got_sitemap.xml")

		got_sitemap, err := os.ReadFile(TestDirPath + "sitemap/got_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of got_sitemap.xml")
		}

		want_sitemap, err := os.ReadFile(TestDirPath + "sitemap/want_sitemap.xml")
		if err != nil {
			t.Errorf("Error in reading the contents of _sitemap.xml")
		}

		got_sitemap_string := string(got_sitemap)
		want_sitemap_string := string(want_sitemap)
		got_sitemap_string = strings.TrimFunc(got_sitemap_string, func(r rune) bool {
			return r == '\n' || r == '\t' || r == ' '
		})
		want_sitemap_string = strings.TrimFunc(want_sitemap_string, func(r rune) bool {
			return r == '\n' || r == '\t' || r == ' '
		})

		if strings.Compare(got_sitemap_string, want_sitemap_string) == 0 {
			t.Errorf("The expected and generated sitemap can be found in test/engine/sitemap/")
		}
	})
}
