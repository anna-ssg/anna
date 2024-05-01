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
	testEngine := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	testEngine.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
	testEngine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)

	testEngine.DeepDataMerge.Templates["index.html"] =
		parser.TemplateData{
			Body:        "<h1>Index Page</h1>",
			CompleteURL: "index.html",
		}

	testEngine.DeepDataMerge.Templates["posts/hello.html"] = parser.TemplateData{
		Body:        "<h1>Hello World</h1>",
		CompleteURL: "posts/hello.html",
	}

	if err := os.MkdirAll(TestDirPath+"render_user_defined/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("render a set of user defined pages", func(t *testing.T) {

		templ, err := template.ParseFiles(TestDirPath + "render_user_defined/template_input.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		testEngine.RenderUserDefinedPages(TestDirPath+"render_user_defined/", templ)

		wantIndexFile, err := os.ReadFile(TestDirPath + "render_user_defined/want_index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		gotIndexFile, err := os.ReadFile(TestDirPath + "render_user_defined/rendered/index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(wantIndexFile, gotIndexFile) {
			t.Errorf("The expected and generated index.html can be found in test/testEngine/render_user_defined/rendered/")
		}

		wantPostHello, err := os.ReadFile(TestDirPath + "render_user_defined/want_post_hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		gotPostHello, err := os.ReadFile(TestDirPath + "render_user_defined/rendered/posts/hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(wantPostHello, gotPostHello) {
			t.Errorf("The expected and generated post/hello.html can be found in test/testEngine/render_user_defined/rendered/posts/")
		}
	})

}
