package parser

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/anna-ssg/anna/v2/pkg/helpers"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/toc"
	"gopkg.in/yaml.v3"
)

type LayoutConfig struct {
	Navbar            []map[string]string `json:"navbar"`
	BaseURL           string              `json:"baseURL"`
	SiteTitle         string              `json:"siteTitle"`
	SiteScripts       []string            `json:"siteScripts"`
	Author            string              `json:"author"`
	Copyright         string              `json:"copyright"`
	ThemeURL          string              `json:"themeURL"`
	Socials           map[string]string   `json:"socials"`
	CollectionLayouts map[string]string   `json:"collectionLayouts"`
}

type Frontmatter struct {
	Title        string              `yaml:"title"`
	Date         string              `yaml:"date"`
	Draft        bool                `yaml:"draft"`
	JSFiles      []string            `yaml:"scripts"`
	Description  string              `yaml:"description"`
	PreviewImage string              `yaml:"previewimage"`
	Tags         []string            `yaml:"tags"`
	TOC          bool                `yaml:"toc"`
	Authors      []string            `yaml:"authors"`
	Collections  []string            `yaml:"collections"`
	Layout       string              `yaml:"layout"`
	CustomFields []map[string]string `yaml:"customFields"`
}

// TemplateData This struct holds all of the data required to render any page of the site
type TemplateData struct {
	CompleteURL template.URL
	Date        int64
	Frontmatter Frontmatter
	Body        template.HTML
	LiveReload  bool
}

type Date int64

type Parser struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates map[template.URL]TemplateData

	// K-V pair storing all templates correspoding to a particular tag in the site
	TagsMap map[template.URL][]TemplateData

	// Collections stores template data of files in collections
	CollectionsMap map[template.URL][]TemplateData

	// K-V pair storing the template layout name for a particular collection in the site
	CollectionsSubPageLayouts map[template.URL]string

	// Stores data parsed from layout/config.yml
	LayoutConfig LayoutConfig

	MdFilesName []string
	MdFilesPath []string

	// Stores flag value to render draft posts
	RenderDrafts bool

	// Common logger for all parser functions
	ErrorLogger *log.Logger

	Helper *helpers.Helper

	// Determines the injection of Live Reload JS in HTML
	LiveReload bool

	// The path to the directory being rendered
	SiteDataPath string
}

func (p *Parser) ParseMDDir(baseDirPath string, baseDirFS fs.FS) {
	helper := helpers.Helper{
		ErrorLogger: p.ErrorLogger,
	}
	err := fs.WalkDir(baseDirFS, ".", func(path string, dir fs.DirEntry, err error) error {
		if path != "." && path != ".obsidian" {
			if dir.IsDir() {
				subDir := os.DirFS(path)
				p.ParseMDDir(path, subDir)
			} else {
				fileName := strings.TrimPrefix(path, baseDirPath)
				if filepath.Ext(path) == ".md" {
					content, err := os.ReadFile(baseDirPath + path)
					if err != nil {
						p.ErrorLogger.Fatal(err)
					}

					frontmatter, body, markdownContent, parseSuccess := p.ParseMarkdownContent(string(content), path)
					if parseSuccess && (p.RenderDrafts || !frontmatter.Draft) {
						p.AddFile(baseDirPath, fileName, frontmatter, markdownContent, body)
					}
				} else {
					helper.CopyFiles(p.SiteDataPath+"content/"+fileName, p.SiteDataPath+"rendered/"+fileName)
				}
			}
		}
		return nil
	})
	if err != nil {
		helper.ErrorLogger.Fatal(err)
	}
}

func (p *Parser) AddFile(baseDirPath string, dirEntryPath string, frontmatter Frontmatter, markdownContent string, body string) {
	p.MdFilesName = append(p.MdFilesName, dirEntryPath)
	testFilepath := baseDirPath + dirEntryPath
	p.MdFilesPath = append(p.MdFilesPath, testFilepath)

	var date int64
	if frontmatter.Date != "" {
		date = p.DateParse(frontmatter.Date).Unix()
	} else {
		date = 0
	}

	key, _ := strings.CutPrefix(testFilepath, p.SiteDataPath+"content/")
	url, _ := strings.CutSuffix(key, ".md")
	url += ".html"

	page := TemplateData{
		CompleteURL: template.URL(url),
		Date:        date,
		Frontmatter: frontmatter,
		Body:        template.HTML(body),
		LiveReload:  p.LiveReload,
	}

	p.Templates[template.URL(url)] = page

	// Adding the page to the tags map with the corresponding tags
	for _, tag := range page.Frontmatter.Tags {
		tagsMapKey := "tags/" + tag + ".html"
		p.TagsMap[template.URL(tagsMapKey)] = append(p.TagsMap[template.URL(tagsMapKey)], page)

	}

	p.collectionsParser(page)
}

func (p *Parser) ParseMarkdownContent(filecontent string, path string) (Frontmatter, string, string, bool) {
	var parsedFrontmatter Frontmatter
	var markdown string
	/*
	   ---
	   frontmatter_content
	   ---

	   markdown content
	   --- => markdown divider and not to be touched while yaml parsing
	*/
	splitContents := strings.Split(filecontent, "---")
	frontmatterSplit := ""

	if len(splitContents) <= 1 {
		p.ErrorLogger.Fatal("Frontmatter missing on path: ", path)
	}

	// If the first section of the page contains a title field, continue parsing
	// Else, prevent parsing of the current file
	regex := regexp.MustCompile(`title(.*): (.*)`)
	match := regex.FindStringSubmatch(splitContents[1])

	if match == nil {
		p.ErrorLogger.Fatal("Title field missing from frontmatter missing on path: ", path)
	}

	frontmatterSplit = splitContents[1]
	// Parsing YAML frontmatter
	err := yaml.Unmarshal([]byte(frontmatterSplit), &parsedFrontmatter)
	if err != nil {
		p.ErrorLogger.Println("Error at path: ", path)
		p.ErrorLogger.Fatal(err)
	}

	if parsedFrontmatter.Layout == "" {
		parsedFrontmatter.Layout = "page"
	}

	markdown = strings.Join(strings.Split(filecontent, "---")[2:], "---")

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer
	var md goldmark.Markdown

	if parsedFrontmatter.TOC {
		md = goldmark.New(
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithExtensions(
				extension.TaskList,
				figure.Figure,
				&toc.Extender{
					Compact: true,
				},
				&mermaid.Extender{
					RenderMode: mermaid.RenderModeClient, // or RenderModeClient
				},
				&anchor.Extender{
					Texter: anchor.Text("#"),
				},
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
	} else {
		md = goldmark.New(
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithExtensions(
				extension.TaskList,
				figure.Figure,
				&mermaid.Extender{
					RenderMode: mermaid.RenderModeClient, // or RenderModeClient
				},
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
	}

	if err := md.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		p.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String(), markdown, true
}

func (p *Parser) DateParse(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	return parsedTime
}

func (p *Parser) ParseConfig(inFilePath string) {
	// // Check if the configuration file exists
	// _, err := os.Stat(inFilePath)
	// if os.IsNotExist(err) {
	// 	p.Helper.Bootstrap()
	// 	return
	// }

	// Read and parse the configuration file
	configFile, err := os.ReadFile(inFilePath)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	err = json.Unmarshal(configFile, &p.LayoutConfig)
	if err != nil {
		p.ErrorLogger.Println("Error at: ", inFilePath)
		p.ErrorLogger.Fatal(err)
	}

	p.parseCollectionLayoutEntries()
}

func (p *Parser) ParseRobots(inFilePath string, outFilePath string) {
	tmpl, err := template.ParseFiles(inFilePath)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, p.LayoutConfig)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	outputFile, err := os.Create(outFilePath)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			p.ErrorLogger.Fatal(err)
		}
	}()

	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
}

// ParseLayoutFiles Parse all the ".html" layout files in the layout/ directory
func (p *Parser) ParseLayoutFiles() *template.Template {

	// Function to check if an element is present in a slice
	templ := template.New("templates").Funcs(template.FuncMap{
		"strSliceContains": func(items []string, search string) bool {
			for _, item := range items {
				if search == item {
					return true
				}
			}
			return false
		},
	})

	// Parsing all files in the layout/ dir hich match the "*.html" pattern
	templ, err := templ.ParseGlob(p.SiteDataPath + "layout/*.html")
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	// Parsing all files in the partials/ dir which match the "*.html" pattern
	templ, err = templ.ParseGlob(p.SiteDataPath + "layout/partials/*.html")
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	return templ
}

// Adding the page to the collections map with the corresponding collections and sub-collections
func (p *Parser) collectionsParser(page TemplateData) {
	// Iterating over all sets of collections defined in the frontmatter
	for _, collectionSet := range page.Frontmatter.Collections {

		var collections []string
		// Collections will be nested using > as the separator - "posts>tech>Go"
		for _, item := range strings.Split(collectionSet, ">") {
			collections = append(collections, strings.TrimSpace(item))
		}

		for i := range len(collections) {
			collectionKey := "collections/"
			for j := range i + 1 {
				collectionKey += collections[j]
				if j != i {
					collectionKey += "/"
				}

			}
			collectionKey += ".html"

			var found bool
			for _, map_page := range p.CollectionsMap[template.URL(collectionKey)] {
				if map_page.CompleteURL == page.CompleteURL {
					found = true
				}
			}
			if !found {
				p.CollectionsMap[template.URL(collectionKey)] = append(p.CollectionsMap[template.URL(collectionKey)], page)
			}

		}

	}
}

func (p *Parser) parseCollectionLayoutEntries() {
	for collectionURL, layoutName := range p.LayoutConfig.CollectionLayouts {
		p.CollectionsSubPageLayouts[template.URL(collectionURL)] = layoutName
	}
}
