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

func TestRenderEngineGeneratedFiles(t *testing.T) {

	engine := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	engine.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
	engine.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData)

	engine.DeepDataMerge.Posts = []parser.TemplateData{
		{
			CompleteURL: "posts/file1.html",
			Frontmatter: parser.Frontmatter{
				Title:       "file1",
				Description: "Description of file 1",
				Date:        "2024-03-28",
			},
		},

		{
			CompleteURL: "posts/file2.html",
			Frontmatter: parser.Frontmatter{
				Title:       "file2",
				Description: "Description of file 2",
				Date:        "2024-03-28",
			},
		},

		{
			CompleteURL: "posts/file3.html",
			Frontmatter: parser.Frontmatter{
				Title:       "file3",
				Description: "Description of file 3",
				Date:        "2024-03-28",
			},
		},
	}

	if err := os.MkdirAll(TestDirPath+"render_engine_generated/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("test rendering of post.html", func(t *testing.T) {
		templ, err := template.ParseFiles(TestDirPath + "render_engine_generated/posts_template.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		engine.RenderEngineGeneratedFiles(TestDirPath+"render_engine_generated/", templ)

		want_posts_file, err := os.ReadFile(TestDirPath + "render_engine_generated/want_posts.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_posts_file, err := os.ReadFile(TestDirPath + "render_engine_generated/rendered/posts.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_posts_file, got_posts_file) {
			t.Errorf("The expected and generated posts.html can be found in test/engine/render_engine_generated/rendered/")
		}
	})
}
