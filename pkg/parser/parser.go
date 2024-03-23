package parser

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

const SiteDataPath string = "site/"

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
	Tags         []string `yaml:"tags"`
}

// This struct holds all of the data required to render any page of the site
type TemplateData struct {
	URL         template.URL
	Filename    string
	Date        int64
	Frontmatter Frontmatter
	Body        template.HTML
	Layout      LayoutConfig

	// Do not use these fields to store tags!!
	// These fields are populated by the ssg to store merged tag data
	Tags       []string
	MergedTags []TemplateData
}

type Date int64

type Parser struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates    map[template.URL]TemplateData
	Posts        []TemplateData
	LayoutConfig LayoutConfig
	TagsMap      map[string][]TemplateData

	ErrorLogger  *log.Logger
	mdFilesName  []string
	mdFilesPath  []string
	RenderDrafts bool
}

func (p *Parser) ReadMdDir(baseDirPath string, baseDirFS fs.FS) {
	// Listing all files in the dirPath directory
	dirEntries, err := os.ReadDir(baseDirPath)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}

	// Storing the markdown file names and paths
	for _, entry := range dirEntries {

		if entry.IsDir() {
			// p.ReadMdDir(baseDirPath + entry.Name() + "/")
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := os.ReadFile(strings.Join([]string{baseDirPath, entry.Name()}, "/"))
		if err != nil {
			p.ErrorLogger.Fatal(err)
		}

		frontmatter, body := p.ParseMarkdownContent(string(content))

		if (frontmatter.Draft && p.RenderDrafts) || !frontmatter.Draft {
			p.AddFileAndRender(baseDirPath, entry, frontmatter, body)
		}
	}
}

func (p *Parser) AddFileAndRender(baseDirPath string, dirEntry fs.DirEntry, frontmatter Frontmatter, body string) {
	p.mdFilesName = append(p.mdFilesName, dirEntry.Name())
	filepath := baseDirPath + dirEntry.Name()
	p.mdFilesPath = append(p.mdFilesPath, filepath)

	var date int64
	if frontmatter.Date != "" {
		date = p.dateParse(frontmatter.Date).Unix()
	} else {
		date = 0
	}

	key, _ := strings.CutPrefix(filepath, SiteDataPath+"content/")
	url, _ := strings.CutSuffix(key, ".md")
	url += ".html"
	page := TemplateData{
		URL:         template.URL(url),
		Date:        date,
		Filename:    strings.Split(dirEntry.Name(), ".")[0],
		Frontmatter: frontmatter,
		Body:        template.HTML(body),
		Layout:      p.LayoutConfig,
	}

	// Adding the page to the merged map storing all site pages
	if frontmatter.Type == "post" {
		p.Posts = append(p.Posts, page)
	}

	p.Templates[template.URL(key)] = page

	// Adding the page to the tags map with the corresponding tags
	for _, tag := range page.Frontmatter.Tags {
		p.TagsMap[tag] = append(p.TagsMap[tag], page)
	}
}

func (p *Parser) ParseMarkdownContent(filecontent string) (Frontmatter, string) {
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
	if len(splitContents) > 1 {
		frontmatterSplit = splitContents[1]
	}

	if frontmatterSplit != "" {
		// Parsing YAML frontmatter
		err := yaml.Unmarshal([]byte(frontmatterSplit), &parsedFrontmatter)
		if err != nil {
			p.ErrorLogger.Fatal(err)
		}

		// Making sure that all filecontent is included and
		// not ignoring the horizontal markdown splitter "---"
		markdown = strings.Join(strings.Split(filecontent, "---")[2:], "---")
	} else {
		markdown = filecontent
	}

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	if err := md.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		p.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String()
}

func (p *Parser) dateParse(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	return parsedTime
}
