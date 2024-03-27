package anna

import (
	"html/template"
	"os"
	"sort"

	"github.com/acmpesuecc/anna/pkg/engine"
	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
)

func VanillaRender() {

	// Defining Engine and Parser Structures
	e := engine.Engine{}
	p := parser.Parser{}

	e.Posts = []parser.TemplateData{}
	e.Templates = make(map[template.URL]parser.TemplateData)
	e.TagsMap = make(map[string][]parser.TemplateData)

	p.ParseConfig(helpers.SiteDataPath + "layout/config.yml")
	fileSystem := os.DirFS(helpers.SiteDataPath + "content/")
	p.ParseMDDir(helpers.SiteDataPath+"content/", fileSystem)
	p.ParseRobots(helpers.SiteDataPath+"layout/robots.txt", helpers.SiteDataPath+"rendered/robots.txt")

	e.GenerateSitemap(helpers.SiteDataPath + "layout/sitemap.xml")
	e.GenerateFeed()

	sort.Slice(e.Posts, func(i, j int) bool {
		return e.Posts[i].Frontmatter.Date > e.Posts[j].Frontmatter.Date
	})

	helper := helpers.Helper{
		ErrorLogger:  e.ErrorLogger,
		SiteDataPath: helpers.SiteDataPath,
	}

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(helpers.SiteDataPath+"static/", helpers.SiteDataPath+"rendered/static/")

	// templ := helper.ParseLayoutFiles()

}
