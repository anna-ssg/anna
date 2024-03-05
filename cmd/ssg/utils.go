package ssg

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

func (g *Generator) readMdDir(dirPath string) {
	// Listing all files in the dirPath directory
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Storing the markdown file names and paths
	for _, entry := range dirEntries {

		if entry.IsDir() {
			g.readMdDir(dirPath + entry.Name() + "/")
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := os.ReadFile(strings.Join([]string{dirPath, entry.Name()}, "/"))
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		frontmatter, body := g.parseMarkdownContent(string(content))

		if frontmatter.Draft && g.RenderDrafts {
			g.AddFileAndRender(dirPath, entry, frontmatter, body)
		} else if frontmatter.Draft && !g.RenderDrafts {
			continue
		}

		if !frontmatter.Draft {
			g.AddFileAndRender(dirPath, entry, frontmatter, body)
		}
	}
}

func (g *Generator) AddFileAndRender(dirPath string, entry fs.DirEntry, frontmatter Frontmatter, body string) {
	g.mdFilesName = append(g.mdFilesName, entry.Name())
	filepath := dirPath + entry.Name()
	g.mdFilesPath = append(g.mdFilesPath, filepath)

	var date int64
	if frontmatter.Date != "" {
		date = g.dateParse(frontmatter.Date).Unix()
	} else {
		date = 0
	}

	page := TemplateData{
		Date:        date,
		Filename:    strings.Split(entry.Name(), ".")[0],
		Frontmatter: frontmatter,
		Body:        template.HTML(body),
		Layout:      g.LayoutConfig,
	}

	key, _ := strings.CutPrefix(filepath, "content/")
	if frontmatter.Type == "post" {
		g.Posts = append(g.Posts, page)
	}

	g.Templates[template.URL(key)] = page
}

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
	for _, templateData := range g.Templates {
		url := g.LayoutConfig.BaseURL + "/" + templateData.Filename + ".html"
		buffer.WriteString(" <url>\n")
		buffer.WriteString("    <loc>" + url + "</loc>\n")
		buffer.WriteString("    <lastmod>" + templateData.Frontmatter.Date + "</lastmod>\n")
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

func (g *Generator) generateFeed() {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<feed xmlns=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("    <title>" + g.LayoutConfig.SiteTitle + "</title>\n")
	buffer.WriteString("    <link href=\"" + g.LayoutConfig.BaseURL + "/" + "\" rel=\"self\"/>\n")
	buffer.WriteString("    <updated>" + time.Now().Format(time.RFC3339) + "</updated>\n")

	// iterate over parsed markdown files that are non-draft posts
	for _, templateData := range g.Templates {
		if !templateData.Frontmatter.Draft {
			buffer.WriteString("    <entry>\n")
			buffer.WriteString("        <title>" + templateData.Frontmatter.Title + "</title>\n")
			buffer.WriteString("        <link href=\"" + g.LayoutConfig.BaseURL + "/posts/" + templateData.Filename + ".html\"/>\n")
			buffer.WriteString("        <id>" + g.LayoutConfig.BaseURL + "/" + templateData.Filename + ".html</id>\n")
			buffer.WriteString("        <updated>" + time.Unix(templateData.Date, 0).Format(time.RFC3339) + "</updated>\n")
			buffer.WriteString("        <content type=\"html\"><![CDATA[" + string(templateData.Body) + "]]></content>\n")
			buffer.WriteString("    </entry>\n")
		}
	}

	buffer.WriteString("</feed>\n")
	outputFile, err := os.Create("rendered/feed.atom")
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

func (g *Generator) dateParse(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	return parsedTime
}
