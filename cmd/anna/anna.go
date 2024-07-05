package anna

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/anna-ssg/anna/v2/pkg/engine"
	"github.com/anna-ssg/anna/v2/pkg/helpers"
	"github.com/anna-ssg/anna/v2/pkg/parser"
)

type Cmd struct {
	RenderDrafts       bool
	Addr               string
	LiveReload         bool
	RenderSpecificSite string
	ServeSpecificSite  string

	// Common logger for all cmd functions
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

type AnnaConfig struct {
	SiteDataPaths map[string]string `json:"siteDataPaths"`
}

func (cmd *Cmd) VanillaRenderManager() {

	// Check if the configuration file exists
	// If it does not, render only the site/ directory

	_, err := os.Stat("anna.json")
	if os.IsNotExist(err) {
		cmd.VanillaRender("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.json")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	err = json.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	// Rendering sites
	if cmd.RenderSpecificSite == "" {
		siteRendered := false

		for _, path := range annaConfig.SiteDataPaths {
			if !siteRendered {
				siteRendered = true
			}
			cmd.VanillaRender(path)
		}

		// If no site has been rendered due to empty "anna.yml", render the default "site/" path
		if !siteRendered {
			cmd.VanillaRender("site/")
		}
	} else {
		siteRendered := false

		for _, sitePath := range annaConfig.SiteDataPaths {
			if strings.Compare(cmd.RenderSpecificSite, sitePath) == 0 {
				cmd.VanillaRender(sitePath)
				siteRendered = true
			}
		}

		if !siteRendered {
			cmd.ErrorLogger.Fatal("Invalid site path to render")
		}

	}

}

func (cmd *Cmd) ValidateHTMLManager() {
	// Rendering all sites
	cmd.VanillaRenderManager()

	// Check if the configuration file exists
	// If it does not, validate only the site/ directory

	_, err := os.Stat("anna.json")
	if os.IsNotExist(err) {
		cmd.VanillaRender("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.json")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	err = json.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	// Validating sites
	validatedSites := false

	for _, sitePath := range annaConfig.SiteDataPaths {
		cmd.ValidateHTMLContent(sitePath)
		if !validatedSites {
			validatedSites = true
		}
	}

	// If no site has been validated due to empty "anna.yml", validate the default "site/" path
	if !validatedSites {
		cmd.ValidateHTMLContent("site/")
	}

}

func (cmd *Cmd) LiveReloadManager() {

	// Check if the configuration file exists
	// If it does not, serve only the site/ directory

	_, err := os.Stat("anna.json")
	if os.IsNotExist(err) {
		cmd.StartLiveReload("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.json")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	err = json.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	// Serving site
	if cmd.ServeSpecificSite == "" {
		cmd.StartLiveReload("site/")
	} else {
		for _, sitePath := range annaConfig.SiteDataPaths {
			if strings.Compare(cmd.ServeSpecificSite, sitePath) == 0 {
				cmd.StartLiveReload(sitePath)
				return
			}
		}

		cmd.ErrorLogger.Fatal("Invalid site path to serve")

	}

}

func (cmd *Cmd) VanillaRender(siteDirPath string) {

	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:                 make(map[template.URL]parser.TemplateData, 10),
		TagsMap:                   make(map[template.URL][]parser.TemplateData, 10),
		CollectionsMap:            make(map[template.URL][]parser.TemplateData, 10),
		CollectionsSubPageLayouts: make(map[template.URL]string, 10),
		SiteDataPath:              siteDirPath,
		ErrorLogger:               log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RenderDrafts:              cmd.RenderDrafts,
		LiveReload:                cmd.LiveReload,
	}

	e := engine.Engine{
		SiteDataPath: siteDirPath,
		ErrorLogger:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
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
