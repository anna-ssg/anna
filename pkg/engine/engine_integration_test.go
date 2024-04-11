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

func TestRenderUserDefinedPages(t *testing.T) {
	engine := engine.Engine{
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[string][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	engine.Templates["index.md"] =
		parser.TemplateData{
			FilenameWithoutExtension: "index",
			Body:                     template.HTML("<h1>Index Page</h1>"),
			CompleteURL:              "index.html",
		}

	engine.Templates["posts/hello.md"] = parser.TemplateData{
		FilenameWithoutExtension: "hello",
		Body:                     template.HTML("<h1>Hello World</h1>"),
		CompleteURL:              "posts/hello.html",
	}

	if err := os.MkdirAll(TestDirPath+"render_user_defined/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("render a set of user defined pages", func(t *testing.T) {

		templ, err := template.ParseFiles(TestDirPath + "render_user_defined/template_input.layout")
		if err != nil {
			t.Errorf("%v", err)
		}

		engine.RenderUserDefinedPages(TestDirPath+"render_user_defined/", templ)

		want_index_file, err := os.ReadFile(TestDirPath + "render_user_defined/want_index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_index_file, err := os.ReadFile(TestDirPath + "render_user_defined/rendered/index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_index_file, got_index_file) {
			t.Errorf("The expected and generated index.html can be found in test/engine/render_user_defined/rendered/")
		}

		want_post_hello, err := os.ReadFile(TestDirPath + "render_user_defined/want_post_hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_post_hello, err := os.ReadFile(TestDirPath + "render_user_defined/rendered/posts/hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_post_hello, got_post_hello) {
			t.Errorf("The expected and generated post/hello.html can be found in test/engine/render_user_defined/rendered/posts/")
		}
	})

	if err := os.RemoveAll(TestDirPath + "render_user_defined/rendered"); err != nil {
		t.Errorf("%v", err)
	}
}
