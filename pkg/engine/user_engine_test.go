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
		}

	engine.Templates["posts/hello.md"] = parser.TemplateData{
		FilenameWithoutExtension: "hello",
		Body:                     template.HTML("<h1>Hello World</h1>"),
	}

	t.Run("render a set of user defined pages", func(t *testing.T) {

		templ, err := template.ParseFiles(TestDirPath + "engine/render_user_defined/template_input.html")
		if err != nil {
			t.Errorf("%v", err)
		}
		engine.RenderUserDefinedPages(TestDirPath+"engine/render_user_defined/", templ)

		want_index_file, err := os.ReadFile(TestDirPath + "engine/render_user_defined/want_index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_index_file, err := os.ReadFile(TestDirPath + "engine/render_user_defined/rendered/index.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_index_file, got_index_file) {
			t.Errorf("The expected and generated index.html can be found in test/engine/render_user_defined/rendered/")
		}

		want_post_hello, err := os.ReadFile(TestDirPath + "engine/render_user_defined/want_post_hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_post_hello, err := os.ReadFile(TestDirPath + "engine/render_user_defined/rendered/posts/hello.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_post_hello, got_post_hello) {
			t.Errorf("The expected and generated post/hello.html can be found in test/engine/render_user_defined/rendered/posts/")
		}
	})
}

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

	t.Run("test rendering of post.html", func(t *testing.T) {
		templ, err := template.ParseFiles(TestDirPath + "engine/render_engine_generated/posts_template.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		engine.RenderEngineGeneratedFiles(TestDirPath+"engine/render_engine_generated/", templ)

		want_posts_file, err := os.ReadFile(TestDirPath + "engine/render_engine_generated/want_posts.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		got_posts_file, err := os.ReadFile(TestDirPath + "engine/render_engine_generated/rendered/posts.html")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(want_posts_file, got_posts_file) {
			t.Errorf("The expected and generated posts.html can be found in test/engine/render_engine_generated/rendered/")
		}
	})
}
