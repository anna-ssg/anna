package parser_test

import (
	"html/template"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/acmpesuecc/anna/pkg/parser"
)

func TestReadMdDir(t *testing.T) {
}

func TestAddFileAndRender(t *testing.T) {
	// t.Run("appending rendered markdown files to parser struct", func(t *testing.T) {
	// })
}

func TestParseMarkdownContent(t *testing.T) {
	p := parser.Parser{
		Templates:   make(map[template.URL]parser.TemplateData),
		TagsMap:     make(map[string][]parser.TemplateData),
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	t.Run("render markdown files to html", func(t *testing.T) {
		inputMd, err := os.ReadFile("../../test/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		_, body_got := p.ParseMarkdownContent(string(inputMd))

		body_want, err := os.ReadFile("../../test/html_want_output.html")
		if err = os.WriteFile("../../test/html_got_output.html", []byte(body_got), 0666); err != nil {
			t.Errorf("%v", err)
		}
		if string(body_want) != body_got {
			t.Errorf("%v\nThe expected and generated html can be found in test/", err)
		}
	})

	t.Run("parse frontmatter from markdown files", func(t *testing.T) {
		inputMd, err := os.ReadFile("../../test/md_inp.md")
		if err != nil {
			t.Errorf("%v", err)
		}

		frontmatter_got, _ := p.ParseMarkdownContent(string(inputMd))
		frontmatter_want := parser.Frontmatter{
			Title:       "Markdown Test",
			Date:        "2024-03-23",
			Description: "File containing markdown to test the SSG",
		}

		if !reflect.DeepEqual(frontmatter_got, frontmatter_want) {
			t.Errorf("got %v, want %v", frontmatter_got, frontmatter_want)
		}
	})
}

// Testing paths of a FS
// func TestReadFSPaths(t *testing.T) {
// 	t.Run("get paths of an fs", func(t *testing.T) {
// 		fsMock := fstest.MapFS{
// 			"index.md":      {Data: []byte("test sentence.")},
// 			"docs.md":       {Data: []byte("test sentence.")},
// 			"post/post1.md": {Data: []byte("test sentence.")},
// 		}
// 		fs.WalkDir(fsMock, ".", func(path string, dir fs.DirEntry, err error) error {
// 			fmt.Printf("%s\n", path)
// 			return nil
// 		})
// 	})
// }
