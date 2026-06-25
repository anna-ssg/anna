package anna

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/anna-ssg/anna/v4/pkg/engine"
	"github.com/anna-ssg/anna/v4/pkg/helpers"
	"github.com/anna-ssg/anna/v4/pkg/logger"
	"github.com/anna-ssg/anna/v4/pkg/parser"
)

type Cmd struct {
	RenderDrafts bool
	Addr         string
	LiveReload   bool
	SiteDirPath  string

	// Common logger for all cmd functions
	ErrorLogger *logger.Logger
	InfoLogger  *logger.Logger
}

func (cmd *Cmd) VanillaRender(siteDirPath string) int {
	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:                 make(map[template.URL]parser.TemplateData, 10),
		SourcePaths:               make(map[template.URL]string, 10),
		TagsMap:                   make(map[template.URL][]parser.TemplateData, 10),
		CollectionsMap:            make(map[template.URL][]parser.TemplateData, 10),
		CollectionsSubPageLayouts: make(map[template.URL]string, 10),
		SiteDataPath:              siteDirPath,
		ErrorLogger:               logger.New(os.Stderr),
		RenderDrafts:              cmd.RenderDrafts,
		LiveReload:                cmd.LiveReload,
	}

	e := engine.Engine{
		SiteDataPath: siteDirPath,
		ErrorLogger:  logger.New(os.Stderr),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData, 10)
	e.DeepDataMerge.SourcePaths = make(map[template.URL]string, 10)
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
	e.DeepDataMerge.SourcePaths = p.SourcePaths
	e.DeepDataMerge.TagsMap = p.TagsMap
	e.DeepDataMerge.CollectionsMap = p.CollectionsMap
	e.DeepDataMerge.CollectionsSubPageLayouts = p.CollectionsSubPageLayouts
	e.DeepDataMerge.LayoutConfig = p.LayoutConfig
	e.BuildInputsModTime = p.BuildInputsModTime

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
	removeStaleRenderedHTML(siteDirPath, expectedRenderedHTML(p, e))

	// Return number of templates/pages rendered for reporting
	return len(e.DeepDataMerge.Templates)
}

func expectedRenderedHTML(p parser.Parser, e engine.Engine) map[string]struct{} {
	expected := make(map[string]struct{}, len(p.Templates)+len(e.DeepDataMerge.TagsMap)+len(e.DeepDataMerge.CollectionsMap)+4)

	for pagePath := range p.Templates {
		expected[filepath.ToSlash(string(pagePath))] = struct{}{}
	}

	expected["tags.html"] = struct{}{}
	for tagPath := range e.DeepDataMerge.TagsMap {
		expected[filepath.ToSlash(string(tagPath))] = struct{}{}
	}

	expected["collections.html"] = struct{}{}
	for collectionPath := range e.DeepDataMerge.CollectionsMap {
		expected[filepath.ToSlash(string(collectionPath))] = struct{}{}
	}

	expected["robots.txt"] = struct{}{}
	expected["sitemap.xml"] = struct{}{}
	expected[filepath.ToSlash(filepath.Join("static", "index.json"))] = struct{}{}

	return expected
}

func removeStaleRenderedHTML(siteDirPath string, expected map[string]struct{}) {
	renderedRoot := filepath.Join(siteDirPath, "rendered")
	_ = filepath.WalkDir(renderedRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		relPath, err := filepath.Rel(renderedRoot, path)
		if err != nil {
			return err
		}

		if _, ok := expected[filepath.ToSlash(relPath)]; ok {
			return nil
		}

		return os.Remove(path)
	})
}
