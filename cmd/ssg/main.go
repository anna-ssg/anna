package ssg

import (
	"bytes"
    "fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

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
}

type Page struct {
    Filename    string
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
	mdPosts      []string
    mdNonDrafts  []string
    Draft        bool
}

func (g *Generator) draftChecker() {
    for _, parsedpage := range g.mdParsed {
        if !parsedpage.Frontmatter.Draft {
            if slices.Contains(g.mdPosts, parsedpage.Filename) {
                g.mdNonDrafts = append(g.mdNonDrafts, parsedpage.Filename)
            }
        }
    }
}

// Write rendered HTML to disk
func (g *Generator) RenderSite(addr string) {
	g.replaceBaseURL(addr)
	g.parseConfig()
	g.readMdDir("content/")
	g.copyStaticContent()
    if !g.Draft {
        g.draftChecker()
    } else {
        g.mdNonDrafts = g.mdPosts
    }

	// Creating the "rendered" directory if not present
	err := os.MkdirAll("rendered/", 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	templ := g.parseLayoutFiles()

	// Writing each parsed markdown file as a separate HTML file
	for i, page := range g.mdParsed {

		// Adding the names of all the files in posts/ dir to the page data
		g.mdParsed[i].Posts = g.mdNonDrafts
		page.Posts = g.mdNonDrafts

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
	err = templ.ExecuteTemplate(&buffer, "posts", g.mdParsed[0])
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Flushing 'posts.html' to the disk
	err = os.WriteFile("rendered/posts.html", buffer.Bytes(), 0666)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}
func (g *Generator) replaceBaseURL(addr string) {
	configFile, err := os.ReadFile("layout/config.yml")
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

    var config LayoutConfig
    if err := yaml.Unmarshal(configFile, &config); err != nil {
        g.ErrorLogger.Fatal(err)
    }

    config.BaseURL = "http://localhost:" + addr + "/"

    updatedConfig, err := yaml.Marshal(config)
    if err != nil {
        g.ErrorLogger.Fatal(err)
    }

    if err = os.WriteFile("layout/config.yml", updatedConfig, 0666); err != nil {
        g.ErrorLogger.Fatal(err)
    }
}

// Serves the rendered files over the address 'addr'
func (g *Generator) ServeSite(addr string) {
	fmt.Println("Serving content at", addr)
	err := http.ListenAndServe(":"+addr, http.FileServer(http.Dir("./rendered")))
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
	frontmatterSplit := strings.Split(filecontent, "---")[1]

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

	g.generateAbsoluteStaticLinks(&markdown)

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

// Make links to static assets load from root dir /
func (g *Generator) generateAbsoluteStaticLinks(mdBody *string) {
	re := regexp.MustCompile(`static\/`)
	absLink := "/" + "static/"
	*mdBody = re.ReplaceAllString(*mdBody, absLink)
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
