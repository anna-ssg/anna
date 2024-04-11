package engine_test

import (
	"html/template"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/acmpesuecc/anna/pkg/engine"
	"github.com/acmpesuecc/anna/pkg/parser"
)

const TestDirPath = "../../test/engine/"

func TestRenderPage(t *testing.T) {
	if err := os.MkdirAll(TestDirPath+"render_page/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("render a single page while creating a new directory", func(t *testing.T) {
		engine := engine.Engine{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		engine.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
		engine.DeepDataMerge.TagsMap = make(map[string][]parser.TemplateData)

		page := parser.TemplateData{
			CompleteURL:              template.URL("got"),
			FilenameWithoutExtension: "got",
			Frontmatter: parser.Frontmatter{
				Title:       "Hello",
				Date:        "2024-03-28",
				Draft:       false,
				Type:        "post",
				Description: "Index page of site",
				Tags:        []string{"blog", "thoughts"},
			},
			Body: template.HTML("<h1>Hello World</h1>"),
			Layout: parser.LayoutConfig{
				Navbar:    []string{"index", "posts"},
				BaseURL:   "https://example.org",
				SiteTitle: "Anna",
				Author:    "anna",
			},
		}

		templ, err := template.ParseFiles(TestDirPath + "render_page/template_input.layout")
		if err != nil {
			t.Errorf("%v", err)
		}

		engine.RenderPage(TestDirPath+"render_page/", "posts/got.html", page, templ, "page")

		got_file, err := os.ReadFile(TestDirPath + "render_page/rendered/posts/got.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_file, err := os.ReadFile(TestDirPath + "render_page/want.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(got_file, want_file) {
			t.Errorf("The expected and generated page.html can be found in test/engine/render_page/rendered/")
		}
	})

	if err := os.RemoveAll(TestDirPath + "render_page/rendered"); err != nil {
		t.Errorf("%v", err)
	}
}
