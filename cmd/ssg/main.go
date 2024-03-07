package ssg

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/acmpesuecc/anna/pkg/helpers"
)

type LayoutConfig struct {
	Navbar      []string `yaml:"navbar"`
	BaseURL     string   `yaml:"baseURL"`
	SiteTitle   string   `yaml:"siteTitle"`
	SiteScripts []string `yaml:"siteScripts"`
	Author      string   `yaml:"author"`
}

type Frontmatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Draft        bool     `yaml:"draft"`
	JSFiles      []string `yaml:"scripts"`
	Type         string   `yaml:"type"`
	Description  string   `yaml:"description"`
	PreviewImage string   `yaml:"previewimage"`
}

type Date int64

type Generator struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates    map[template.URL]TemplateData
	Posts        []TemplateData
	LayoutConfig LayoutConfig

	ErrorLogger  *log.Logger
	mdFilesName  []string
	mdFilesPath  []string
	RenderDrafts bool
}

// This struct holds all of the data required to render any page of the site
// Pass this struct without modification to ExecuteTemplate()
type TemplateData struct {
	Filename    string
	Date        int64
	Frontmatter Frontmatter
	Body        template.HTML
	Layout      LayoutConfig
}

// This struct holds the data required to render posts.html
type postsTemplateData struct {
	Posts []TemplateData
	TemplateData
}

func (g *Generator) RenderSite(addr string) {
	// Creating the "rendered" directory if not present
	err := os.RemoveAll("rendered/")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll("rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	g.Posts = []TemplateData{}
	g.parseConfig()
	g.readMdDir("content/")
	g.parseRobots()
	g.generateSitemap()
	g.generateFeed()

	sort.Slice(g.Posts, func(i, j int) bool {
		return g.Posts[i].Frontmatter.Date > g.Posts[j].Frontmatter.Date
	})

	helper := helpers.Helper{
		ErrorLogger: g.ErrorLogger,
	}

	// Copies the contents of the 'static/' directory to 'rendered/'
	helper.CopyDirectoryContents("static/", "rendered/static/")
	helper.CopyDirectoryContents("script/", "rendered/script/")

	template := helper.ParseLayoutFiles()

	for pagePath, templateData := range g.Templates {
		g.RenderPage(pagePath, templateData, template)
	}

	var buffer bytes.Buffer

	postsData := postsTemplateData{
		Posts: g.Posts,
		TemplateData: TemplateData{
			Frontmatter: Frontmatter{Title: "Posts"},
			Layout:      g.LayoutConfig,
		},
	}

	err = template.ExecuteTemplate(&buffer, "posts", postsData)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile("rendered/posts.html", buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) RenderPage(pagePath template.URL, templateData TemplateData, template *template.Template) {

	// Creating subdirectories if the filepath contains '/'
	if strings.Contains(string(pagePath), "/") {
		// Extracting the directory path from the filepath
		dirPath, _ := strings.CutSuffix(string(pagePath), templateData.Filename+".md")
		dirPath = "rendered/" + dirPath

		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
	}

	filename, _ := strings.CutSuffix(string(pagePath), ".md")
	filepath := "rendered/" + filename + ".html"
	var buffer bytes.Buffer

	// Storing the rendered HTML file to a buffer
	err := template.ExecuteTemplate(&buffer, "page", templateData)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(filepath, buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}
