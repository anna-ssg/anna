package parser_test

import (
	"html/template"
	"log"
	"os"
	"testing"

	"github.com/acmpesuecc/anna/pkg/parser"
)

func TestParseMDDir(t *testing.T) {
	t.Run("reading markdown files and rendering without drafts", func(t *testing.T) {
		p := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		p.RenderDrafts = false

		TestDirFS := os.DirFS(TestDirPath + "input")
		p.ParseMDDir(TestDirPath+"input/", TestDirFS)

		got_parsed_files := len(p.MdFilesName)
		want_parsed_files := 1

		if got_parsed_files != want_parsed_files {
			t.Errorf("got %v, want %v", got_parsed_files, want_parsed_files)
		}

	})

	t.Run("reading all markdown files inluding drafts", func(t *testing.T) {
		p := parser.Parser{
			Templates:   make(map[template.URL]parser.TemplateData),
			TagsMap:     make(map[template.URL][]parser.TemplateData),
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		p.RenderDrafts = true

		inpBaseDirFS := os.DirFS(TestDirPath + "input")
		p.ParseMDDir(TestDirPath+"input/", inpBaseDirFS)

		got_parsed_files := len(p.MdFilesName)
		want_parsed_files := 2
		if got_parsed_files != want_parsed_files {
			t.Errorf("got %v, want %v", got_parsed_files, want_parsed_files)
		}
	})
}
