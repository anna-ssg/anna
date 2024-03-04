package ssg

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"sort"

	//"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type LayoutConfig struct {
	Navbar  []string `yaml:"navbar"`
	BaseURL string   `yaml:"baseURL"`
}

type Frontmatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
	Draft bool   `yaml:"draft"`
    Type  string `yaml:"type"`
    Description string `yaml:"description"`
}

type Page struct {
	Filename    string
    Date        int64
	Frontmatter Frontmatter
	Body        template.HTML
	Layout      LayoutConfig
	Posts       []string
}

type Generator struct {
	ErrorLogger  *log.Logger
	mdFilesName  []string
	mdFilesPath  []string
	mdParsed     []Page
	LayoutConfig LayoutConfig
	MdPosts      []Page
	Draft        bool
}

func getFileNames(page []Page) []string {
    var filenames []string
    for _, p := range page {
        filenames = append(filenames, p.Filename)
    }
    return filenames
}

func (g *Generator) dateParse(date string) time.Time {
    parsedTime, err := time.Parse("2006-01-02", date)
    if err != nil {
        g.ErrorLogger.Fatal(err)
    }
    return parsedTime
}

// Write rendered HTML to disk
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

	g.parseConfig()
	g.mdPosts = []string{}
	g.readMdDir("content/")
	g.parseRobots()
	g.generateSitemap()
	g.copyStaticContent()
	templ := g.parseLayoutFiles()

	// Writing each parsed markdown file as a separate HTML file
	for i, page := range g.mdParsed {

		// Adding the names of all the files in posts/ dir to the page data
		g.mdParsed[i].Posts = getFileNames(g.MdPosts)
		page.Posts = getFileNames(g.MdPosts)

		filename, _ := strings.CutPrefix(g.mdFilesPath[i], "content/")

		// Creating subdirectories if the filepath contains '/'
		if strings.Contains(filename, "/") {
			// Extracting the directory path from the filepath
			dirPath, _ := strings.CutSuffix(filename, g.mdFilesName[i])
			dirPath = "rendered/" + dirPath

			err := os.MkdirAll(dirPath, 0750)
			if err != nil {
				g.ErrorLogger.Fatal(err)
			}
		}

		filename, _ = strings.CutSuffix(filename, ".md")
		filepath := "rendered/" + filename + ".html"
		var buffer bytes.Buffer

		// Storing the rendered HTML file to a buffer
		err = templ.ExecuteTemplate(&buffer, "page", page)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		// Flushing data from the buffer to the disk
		err := os.WriteFile(filepath, buffer.Bytes(), 0666)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}
	}

	var buffer bytes.Buffer
	// Rendering the 'posts.html' separately


    out := g.MdPosts

    sort.Slice(out, func(i, j int) bool {
        return out[i].Date > out[j].Date
    })


    type TemplateData struct {
        Generator       *Generator
        Frontmatter     Frontmatter
        Layout         LayoutConfig
    }
    data := TemplateData{
        Generator: g,
        Frontmatter: Frontmatter{Title: "Posts"},
        Layout: g.LayoutConfig,
    }


	err = templ.ExecuteTemplate(&buffer, "posts", data)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile("rendered/posts.html", buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) parseRobots() {
	tmpl, err := template.ParseFiles("layout/robots.txt")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, g.LayoutConfig)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	outputFile, err := os.Create("rendered/robots.txt")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) generateSitemap() {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// iterate over parsed markdown files
	for _, page := range g.mdParsed {
		url := g.LayoutConfig.BaseURL + page.Filename + ".html"
		buffer.WriteString(" <url>\n")
		buffer.WriteString("    <loc>" + url + "</loc>\n")
		buffer.WriteString("    <lastmod>" + page.Frontmatter.Date + "</lastmod>\n")
		buffer.WriteString(" </url>\n")
	}
	buffer.WriteString("</urlset>\n")
	outputFile, err := os.Create("rendered/sitemap.xml")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

func (g *Generator) parseMarkdownContent(filecontent string) (Frontmatter, string) {
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
			g.ErrorLogger.Fatal(err)
		}

		// we want to make sure that all filecontent is included and
		// not ignoring the horizontal markdown splitter "---"
		markdown = strings.Join(strings.Split(filecontent, "---")[2:], "")
	} else {
		markdown = filecontent
	}

	// Parsing markdown to HTML
	var parsedMarkdown bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &parsedMarkdown); err != nil {
		g.ErrorLogger.Fatal(err)
	}

	return parsedFrontmatter, parsedMarkdown.String()
}

// Copies the contents of the 'static/' directory to 'rendered/'
func (g *Generator) copyStaticContent() {
	g.copyDirectoryContents("static/", "rendered/static/")
}

// Parse 'config.yml' to configure the layout of the site
func (g *Generator) parseConfig() {
	configFile, err := os.ReadFile("layout/config.yml")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = yaml.Unmarshal(configFile, &g.LayoutConfig)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

// Parse all the ".html" layout files in the layout/ directory
func (g *Generator) parseLayoutFiles() *template.Template {
	// Parsing all files in the layout/ dir which match the "*.html" pattern
	templ, err := template.ParseGlob("layout/*.html")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Parsing all files in the partials/ dir which match the "*.html" pattern
	templ, err = templ.ParseGlob("layout/partials/*.html")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	return templ
}
