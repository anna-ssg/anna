package parser

import (
	"bytes"
	"regexp"

	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	URL                      template.URL
	FilenameWithoutExtension string
	Date                     int64
	Frontmatter              Frontmatter
	Body                     template.HTML
	Layout                   LayoutConfig

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
	MdFilesName  []string
	MdFilesPath  []string
	RenderDrafts bool
}

func (p *Parser) ReadMdDir(baseDirPath string, baseDirFS fs.FS) {
	fs.WalkDir(baseDirFS, ".", func(path string, dir fs.DirEntry, err error) error {
		if dir.IsDir() && path != "." {
			subDir := os.DirFS(path)
			p.ReadMdDir(path, subDir)
		} else {
			if strings.HasSuffix(path, ".md") {
				fileName := filepath.Base(path)

				content, err := os.ReadFile(baseDirPath + "/" + path)
				if err != nil {
					p.ErrorLogger.Fatal(err)
				}

				fronmatter, body, parseSuccess := p.ParseMarkdownContent(string(content))
				if parseSuccess {
					if (fronmatter.Draft && p.RenderDrafts) || !fronmatter.Draft {
						p.AddFileAndRender(baseDirPath, fileName, fronmatter, body)
					}
				}
			}
		}
		return nil
	})
}

func (p *Parser) AddFileAndRender(baseDirPath string, dirEntryPath string, frontmatter Frontmatter, body string) {
	p.MdFilesName = append(p.MdFilesName, dirEntryPath)
	filepath := baseDirPath + dirEntryPath
	p.MdFilesPath = append(p.MdFilesPath, filepath)

	var date int64
	if frontmatter.Date != "" {
		date = p.DateParse(frontmatter.Date).Unix()
	} else {
		date = 0
	}

	key, _ := strings.CutPrefix(filepath, SiteDataPath+"content/")
	url, _ := strings.CutSuffix(key, ".md")
	url += ".html"
	page := TemplateData{
		URL:                      template.URL(url),
		Date:                     date,
		FilenameWithoutExtension: strings.Split(dirEntryPath, ".")[0],
		Frontmatter:              frontmatter,
		Body:                     template.HTML(body),
		Layout:                   p.LayoutConfig,
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

func (p *Parser) ParseMarkdownContent(filecontent string) (Frontmatter, string, bool) {
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

	regex := regexp.MustCompile(`title: (.*)`)
	match := regex.FindStringSubmatch(splitContents[1])

	if match == nil {
		return Frontmatter{}, "", false
	}

	frontmatterSplit = splitContents[1]
	// Parsing YAML frontmatter
	err := yaml.Unmarshal([]byte(frontmatterSplit), &parsedFrontmatter)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	markdown = strings.Join(strings.Split(filecontent, "---")[2:], "---")

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

	return parsedFrontmatter, parsedMarkdown.String(), true
}

func (p *Parser) DateParse(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		p.ErrorLogger.Fatal(err)
	}
	return parsedTime
}