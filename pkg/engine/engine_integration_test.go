package engine_test

import (
	"html/template"
	"os"
	"slices"
	"testing"

	"github.com/anna-ssg/anna/pkg/engine"
	"github.com/anna-ssg/anna/pkg/logger"
	"github.com/anna-ssg/anna/pkg/parser"
)

func TestRenderUserDefinedPages(t *testing.T) {
	testEngine := engine.Engine{
		ErrorLogger: logger.New(os.Stderr),
	}
	testEngine.DeepDataMerge.Pages = make(map[template.URL]parser.PageData)
	testEngine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.PageData)

	testEngine.DeepDataMerge.Pages["index.html"] =
		parser.PageData{
			Body:        "<h1>Index Page</h1>",
			CompleteURL: "index.html",
			Frontmatter: parser.Frontmatter{
				Layout: "page",
			},
		}

	testEngine.DeepDataMerge.Pages["posts/hello.html"] = parser.PageData{
		Body:        "<h1>Hello World</h1>",
		CompleteURL: "posts/hello.html",
		Frontmatter: parser.Frontmatter{
			Layout: "page",
		},
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
