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
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[string][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),

		Posts: []parser.TemplateData{
			{
				FilenameWithoutExtension: "file1",
				Frontmatter: parser.Frontmatter{
					Title:       "file1",
					Description: "Description of file 1",
					Date:        "2024-03-28",
				},
			},

			{
				FilenameWithoutExtension: "file2",
				Frontmatter: parser.Frontmatter{
					Title:       "file2",
					Description: "Description of file 2",
					Date:        "2024-03-28",
				},
			},

			{
				FilenameWithoutExtension: "file3",
				Frontmatter: parser.Frontmatter{
					Title:       "file3",
					Description: "Description of file 3",
					Date:        "2024-03-28",
				},
			},
		},
	}

	if err := os.MkdirAll(TestDirPath+"render_engine_generated/rendered", 0750); err != nil {
		t.Errorf("%v", err)
	}

	t.Run("test rendering of post.html", func(t *testing.T) {
		templ, err := template.ParseFiles(TestDirPath + "render_engine_generated/posts_template.layout")
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

	if err := os.RemoveAll(TestDirPath + "render_engine_generated/rendered"); err != nil {
		t.Errorf("%v", err)
	}
}
