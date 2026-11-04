package parser_test

import (
	"html/template"
	"os"
	"testing"

	"github.com/anna-ssg/anna/pkg/logger"
	"github.com/anna-ssg/anna/pkg/parser"
)

func TestParseMDDir(t *testing.T) {
	t.Run("reading markdown files and rendering without drafts", func(t *testing.T) {
		p := parser.Parser{
			Templates:            make(map[template.URL]parser.TemplateData),
			TagsMap:              make(map[template.URL][]parser.TemplateData),
			CollectionsMap:       make(map[template.URL][]parser.TemplateData),
			CollectionsHierarchy: make(map[template.URL]*parser.CollectionHierarchyData),
			ErrorLogger:          logger.New(os.Stderr),
		}
		p.RenderDrafts = false

		TestDirFS := os.DirFS(TestDirPath + "input")
		p.ParseMDDir(TestDirPath+"input/", TestDirFS)

		gotParsedFiles := len(p.MdFilesName)
		wantParsedFiles := 1

		if gotParsedFiles != wantParsedFiles {
			t.Errorf("got %v, want %v", gotParsedFiles, wantParsedFiles)
		}

	})

	t.Run("reading all markdown files inluding drafts", func(t *testing.T) {
		p := parser.Parser{
			Templates:            make(map[template.URL]parser.TemplateData),
			TagsMap:              make(map[template.URL][]parser.TemplateData),
			CollectionsMap:       make(map[template.URL][]parser.TemplateData),
			CollectionsHierarchy: make(map[template.URL]*parser.CollectionHierarchyData),
			ErrorLogger:          logger.New(os.Stderr),
		}
		p.RenderDrafts = true

		inpBaseDirFS := os.DirFS(TestDirPath + "input")
		p.ParseMDDir(TestDirPath+"input/", inpBaseDirFS)

		gotParsedFiles := len(p.MdFilesName)
		wantParsedFiles := 2
		if gotParsedFiles != wantParsedFiles {
			t.Errorf("got %v, want %v", gotParsedFiles, wantParsedFiles)
		}
	})
}
