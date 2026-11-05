package engine_test

import (
	"html/template"
	"io"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/anna-ssg/anna/pkg/engine"
	"github.com/anna-ssg/anna/pkg/parser"

	"github.com/rs/zerolog"
)

const TestDirPath = "../../test/engine/"

func InitLogger(out io.Writer) *zerolog.Logger {
	console := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	}
	logger := zerolog.New(console).With().Timestamp().Caller().Stack().Logger()
	return &logger
}

func TestRenderPage(t *testing.T) {
	if err := os.MkdirAll(TestDirPath+"render_page/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("render a single page while creating a new directory", func(t *testing.T) {
		testEngine := engine.Engine{
			Logger: InitLogger(os.Stdout),
		}
		testEngine.DeepDataMerge.Pages = make(map[template.URL]parser.PageData)
		testEngine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.PageData)

		testEngine.DeepDataMerge.Pages["posts/got.html"] = parser.PageData{
			CompleteURL: "got.html",
			Frontmatter: parser.Frontmatter{
				Title:       "Hello",
				Date:        "2024-03-28",
				Draft:       false,
				Description: "Index page of site",
				Tags:        []string{"blog", "thoughts"},
				Layout:      "page",
			},
			Body: "<h1>Hello World</h1>",
		}

		templ, err := template.ParseFiles(TestDirPath + "render_page/template_input.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		testEngine.RenderPage(TestDirPath+"render_page/", "posts/got.html", templ, "page")

		gotFile, err := os.ReadFile(TestDirPath + "render_page/rendered/posts/got.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantFile, err := os.ReadFile(TestDirPath + "render_page/want.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(gotFile, wantFile) {
			t.Errorf("The expected and generated page.html can be found in test/testEngine/render_page/rendered/")
		}
	})

}
