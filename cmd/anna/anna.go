package anna

import (
	"html/template"
	"log"
	"os"

	"github.com/anna-ssg/anna/v3/pkg/engine"
	"github.com/anna-ssg/anna/v3/pkg/helpers"
	"github.com/anna-ssg/anna/v3/pkg/parser"
)

type Cmd struct {
	RenderDrafts bool
	Addr         string
	LiveReload   bool
	SiteDirPath  string

	// Common logger for all cmd functions
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

func (cmd *Cmd) VanillaRender(siteDirPath string) {
	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:                 make(map[template.URL]parser.TemplateData, 10),
		TagsMap:                   make(map[template.URL][]parser.TemplateData, 10),
		CollectionsMap:            make(map[template.URL][]parser.TemplateData, 10),
		CollectionsSubPageLayouts: make(map[template.URL]string, 10),
		SiteDataPath:              siteDirPath,
		ErrorLogger:               log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		RenderDrafts:              cmd.RenderDrafts,
		LiveReload:                cmd.LiveReload,
	}

	e := engine.Engine{
		SiteDataPath: siteDirPath,
		ErrorLogger:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData, 10)
	e.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData, 10)
	e.DeepDataMerge.CollectionsMap = make(map[template.URL][]parser.TemplateData, 10)

	helper := helpers.Helper{
		ErrorLogger: e.ErrorLogger,
	}

	helper.CreateRenderedDir(siteDirPath)

	p.ParseConfig(siteDirPath + "layout/config.json")
	p.ParseRobots(siteDirPath+"layout/robots.txt", siteDirPath+"rendered/robots.txt")

	fileSystem := os.DirFS(siteDirPath + "content/")
	p.ParseMDDir(siteDirPath+"content/", fileSystem)

	templ := p.ParseLayoutFiles()

	e.DeepDataMerge.Templates = p.Templates
	e.DeepDataMerge.TagsMap = p.TagsMap
	e.DeepDataMerge.CollectionsMap = p.CollectionsMap
	e.DeepDataMerge.CollectionsSubPageLayouts = p.CollectionsSubPageLayouts
	e.DeepDataMerge.LayoutConfig = p.LayoutConfig

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(siteDirPath+"static/", siteDirPath+"rendered/static/")

	// Check if the public folder exists ands copy contents

	_, err := os.Stat(siteDirPath + "public/")
	if os.IsNotExist(err) {
	} else {
		// Check if the public folder exists ands copy contents

		_, err := os.Stat(siteDirPath + "public/")
		if os.IsNotExist(err) {
		} else {
			// Copies the contents of the 'static/' directory to 'rendered/'
			helper.CopyDirectoryContents(siteDirPath+"public/", siteDirPath+"rendered/")
		}
	}

	e.GenerateSitemap(siteDirPath + "rendered/sitemap.xml")
	e.GenerateFeed()
	e.GenerateJSONIndex(siteDirPath)

	e.RenderUserDefinedPages(siteDirPath, templ)
	e.RenderTags(siteDirPath, templ)
	e.RenderCollections(siteDirPath, templ)
}
