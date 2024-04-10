package anna

import (
	"html/template"
	"log"
	"os"

	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
	zettel_engine "github.com/acmpesuecc/anna/pkg/zettel/engine"
	zettel_parser "github.com/acmpesuecc/anna/pkg/zettel/parser"
)

func (cmd *Cmd) VanillaNoteRender(LayoutConfig parser.LayoutConfig) {
	p := zettel_parser.Parser{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	p.NotesMergedData.Notes = make(map[template.URL]zettel_parser.Note)
	p.NotesMergedData.LinkStore = make(map[template.URL][]*zettel_parser.Note)

	fileSystem := os.DirFS(helpers.SiteDataPath + "content/notes")
	p.Layout = LayoutConfig
	p.ParseNotesDir(helpers.SiteDataPath+"content/notes/", fileSystem)

	e := zettel_engine.Engine{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	e.NotesMergedData.Notes = make(map[template.URL]zettel_parser.Note)
	e.NotesMergedData.LinkStore = make(map[template.URL][]*zettel_parser.Note)
	e.NotesMergedData = p.NotesMergedData

	templ, err := template.ParseGlob(helpers.SiteDataPath + "layout/notes/*.layout")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	templ, err = templ.ParseGlob(helpers.SiteDataPath + "layout/partials/*.layout")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	e.GenerateLinkStore()
	e.RenderUserNotes(helpers.SiteDataPath, templ)
	e.GenerateRootNote(helpers.SiteDataPath, templ)
}
