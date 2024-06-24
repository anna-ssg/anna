package anna

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/anna-ssg/anna/v2/pkg/engine"
	"github.com/anna-ssg/anna/v2/pkg/helpers"
	"github.com/anna-ssg/anna/v2/pkg/parser"
	"gopkg.in/yaml.v3"
)

type Cmd struct {
	RenderDrafts      bool
	Addr              string
	LiveReload        bool
	RenderAll         bool
	RenderSpecific    string
	ServeSpecificSite string

	// Common logger for all cmd functions
	ErrorLogger *log.Logger
}

type AnnaConfig struct {
	SiteDataPaths []map[string]string `yaml:"siteDataPaths"`
}

func (cmd *Cmd) VanillaRenderManager() {

	// Check if the configuration file exists
	// If it does not, render only the site/ directory

	_, err := os.Stat("anna.yml")
	if os.IsNotExist(err) {
		cmd.VanillaRender("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.yml")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	// There are two sites and user wants to render both

	// There are two sites and user wants to serve one

	err = yaml.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	if cmd.RenderAll {

	} else if cmd.RenderSpecific != "" {
		cmd.ErrorLogger.Fatalln("Did not provide a valid site name present in anna.yml")
	} else {
		if entry := annaConfig.SiteDataPaths[cmd.RenderSpecific]; entry {
			cmd.VanillaRender(cmd.RenderSpecific)
		}
	}

}

// todo
func (cmd *Cmd) ValidateHTMLManager() {

	// Check if the configuration file exists
	// If it does not, render only the site/ directory

	_, err := os.Stat("anna.yml")
	if os.IsNotExist(err) {
		cmd.ValidateHTMLContent("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.yml")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	err = yaml.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	fmt.Println(annaConfig)
}

func (cmd *Cmd) LiveReloadManager() {

	// Check if the configuration file exists
	// If it does not, render only the site/ directory

	_, err := os.Stat("anna.yml")
	if os.IsNotExist(err) {
		cmd.StartLiveReload("site/")
		return
	}

	// Read and parse the configuration file
	annaConfigFile, err := os.ReadFile("anna.yml")
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	var annaConfig AnnaConfig

	err = yaml.Unmarshal(annaConfigFile, &annaConfig)
	if err != nil {
		cmd.ErrorLogger.Fatal(err)
	}

	if cmd.ServeSpecificSite != "" {

	}
}

func (cmd *Cmd) VanillaRender(siteDirPath string) {

	// Defining Engine and Parser Structures
	p := parser.Parser{
		Templates:      make(map[template.URL]parser.TemplateData, 10),
		TagsMap:        make(map[template.URL][]parser.TemplateData, 10),
		CollectionsMap: make(map[template.URL][]parser.TemplateData, 10),
		Notes:          make(map[template.URL]parser.Note, 10),
		ErrorLogger:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RenderDrafts:   cmd.RenderDrafts,
		LiveReload:     cmd.LiveReload,
	}

	e := engine.Engine{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	e.DeepDataMerge.Templates = make(map[template.URL]parser.TemplateData, 10)
	e.DeepDataMerge.TagsMap = make(map[template.URL][]parser.TemplateData, 10)
	e.DeepDataMerge.CollectionsMap = make(map[template.URL][]parser.TemplateData, 10)
	e.DeepDataMerge.Notes = make(map[template.URL]parser.Note, 10)
	e.DeepDataMerge.LinkStore = make(map[template.URL][]*parser.Note, 10)

	helper := helpers.Helper{
		ErrorLogger: e.ErrorLogger,
	}

	helper.CreateRenderedDir(siteDirPath)

	p.ParseConfig(siteDirPath + "layout/config.yml")
	p.ParseRobots(siteDirPath+"layout/robots.txt", siteDirPath+"rendered/robots.txt")

	fileSystem := os.DirFS(siteDirPath + "content/")
	p.ParseMDDir(siteDirPath+"content/", fileSystem)

	p.ParseLayoutFiles()

	// Generate backlinks and validations for notes
	p.BackLinkParser()

	e.DeepDataMerge.Templates = p.Templates
	e.DeepDataMerge.TagsMap = p.TagsMap
	e.DeepDataMerge.CollectionsMap = p.CollectionsMap
	e.DeepDataMerge.LayoutConfig = p.LayoutConfig
	e.DeepDataMerge.Posts = p.Posts
	e.DeepDataMerge.Notes = p.Notes

	sort.Slice(e.DeepDataMerge.Posts, func(i, j int) bool {
		return e.DeepDataMerge.Posts[i].Frontmatter.Date > e.DeepDataMerge.Posts[j].Frontmatter.Date
	})

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(siteDirPath+"static/", siteDirPath+"rendered/static/")

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents(siteDirPath+"public/", siteDirPath+"rendered/")

	e.GenerateSitemap(siteDirPath + "rendered/sitemap.xml")
	e.GenerateFeed()
	e.GenerateJSONIndex(siteDirPath)

	e.GenerateLinkStore()
	e.GenerateNoteJSONIdex(siteDirPath)

	templ, err := template.ParseGlob(siteDirPath + "layout/*.html")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	templ, err = templ.ParseGlob(siteDirPath + "layout/partials/*.html")
	if err != nil {
		e.ErrorLogger.Fatalf("%v", err)
	}

	e.RenderNotes(siteDirPath, templ)
	e.GenerateNoteRoot(siteDirPath, templ)
	e.RenderEngineGeneratedFiles(siteDirPath, templ)
	e.RenderUserDefinedPages(siteDirPath, templ)
	e.RenderTags(siteDirPath, templ)
	e.RenderCollections(siteDirPath, templ)
}
