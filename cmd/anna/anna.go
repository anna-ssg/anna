package anna

import (
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/acmpesuecc/anna/pkg/engine"
	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
)

type Cmd struct {
	RenderDrafts bool
	Addr         string
}

func (cmd *Cmd) VanillaRender() {
	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:    make(map[template.URL]parser.TemplateData),
		TagsMap:      make(map[string][]parser.TemplateData),
		ErrorLogger:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RenderDrafts: cmd.RenderDrafts,
	}

	e := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData)
	e.DeepDataMerge.TagsMap = make(map[string][]parser.TemplateData)

	helper := helpers.Helper{
		ErrorLogger:  e.ErrorLogger,
		SiteDataPath: helpers.SiteDataPath,
	}

	helper.CreateRenderedDir(helper.SiteDataPath)

	// Copies the contents of the 'static/' directory to 'rendered/'

	p.ParseConfig(helpers.SiteDataPath + "layout/config.yml")

	fileSystem := os.DirFS(helpers.SiteDataPath + "content/")
	p.ParseMDDir(helpers.SiteDataPath+"content/", fileSystem)

	p.ParseRobots(helpers.SiteDataPath+"layout/robots.txt", helpers.SiteDataPath+"rendered/robots.txt")
	p.ParseLayoutFiles()

	e.DeepDataMerge.Templates = p.Templates
	e.DeepDataMerge.TagsMap = p.TagsMap
	e.DeepDataMerge.LayoutConfig = p.LayoutConfig
	e.DeepDataMerge.Posts = p.Posts

	e.GenerateSitemap(helpers.SiteDataPath + "rendered/sitemap.xml")
	e.GenerateFeed()
	e.GenerateJSONIndex(helpers.SiteDataPath)
	helper.CopyDirectoryContents(helpers.SiteDataPath+"static/", helpers.SiteDataPath+"rendered/static/")

	sort.Slice(e.DeepDataMerge.Posts, func(i, j int) bool {
		return e.DeepDataMerge.Posts[i].Frontmatter.Date > e.DeepDataMerge.Posts[j].Frontmatter.Date
	})

	templ, err := template.ParseGlob(helpers.SiteDataPath + "layout/*.layout")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	templ, err = templ.ParseGlob(helpers.SiteDataPath + "layout/partials/*.layout")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}
	e.RenderEngineGeneratedFiles(helpers.SiteDataPath, templ)
	e.RenderUserDefinedPages(helpers.SiteDataPath, templ)

	e.RenderTags(helpers.SiteDataPath, templ)
}
