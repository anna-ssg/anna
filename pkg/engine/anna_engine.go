package engine

import (
	"bytes"
	"cmp"
	"html/template"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/acmpesuecc/anna/pkg/helpers"
	"github.com/acmpesuecc/anna/pkg/parser"
)

func (e *Engine) RenderTags(fileOutPath string, templ *template.Template) {
	var tagsBuffer bytes.Buffer

	// Extracting tag titles
	tags := make([]string, 0, len(e.TagsMap))
	for tag := range e.TagsMap {
		tags = append(tags, tag)
	}

	slices.SortFunc(tags, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	tagNames := parser.TemplateData{
		FilenameWithoutExtension: "Tags",
		Layout:                   e.LayoutConfig,
		Frontmatter:              parser.Frontmatter{Title: "Tags"},
		Tags:                     tags,
	}

	// Rendering the page displaying all tags
	err := templ.ExecuteTemplate(&tagsBuffer, "all-tags", tagNames)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing 'tags.html' to the disk
	err = os.WriteFile(fileOutPath+"rendered/tags.html", tagsBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Rendering the subpages with merged tagged posts
	for tag, taggedTemplates := range e.TagsMap {
		pagePath := "tags/" + tag
		templateData := parser.TemplateData{
			FilenameWithoutExtension: tag,
			Layout:                   e.LayoutConfig,
			Frontmatter: parser.Frontmatter{
				Title: tag,
			},
			SpecificTagTemplates: taggedTemplates,
		}

		e.ErrorLogger.Println(fileOutPath)
		e.ErrorLogger.Println(pagePath)
		e.RenderPage(fileOutPath, template.URL(pagePath), templateData, templ, "tag-subpage")
	}
}

func (e *Engine) GenerateSitemap(outFilePath string) {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// iterate over parsed markdown files
	for _, templateData := range e.Templates {
		url := e.LayoutConfig.BaseURL + "/" + templateData.FilenameWithoutExtension + ".html"
		buffer.WriteString(" <url>\n")
		buffer.WriteString("  <loc>" + url + "</loc>\n")
		buffer.WriteString("  <lastmod>" + templateData.Frontmatter.Date + "</lastmod>\n")
		buffer.WriteString(" </url>\n")
	}
	buffer.WriteString("</urlset>\n")
	// helpers.SiteDataPath is the DirPath
	outputFile, err := os.Create(outFilePath)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) GenerateFeed() {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<feed xmlns=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("    <title>" + e.LayoutConfig.SiteTitle + "</title>\n")
	buffer.WriteString("    <link href=\"" + e.LayoutConfig.BaseURL + "/" + "\" rel=\"self\"/>\n")
	buffer.WriteString("    <updated>" + time.Now().Format(time.RFC3339) + "</updated>\n")

	// iterate over parsed markdown files that are non-draft posts
	for _, templateData := range e.Templates {
		if !templateData.Frontmatter.Draft {
			buffer.WriteString("    <entry>\n")
			buffer.WriteString("        <title>" + templateData.Frontmatter.Title + "</title>\n")
			buffer.WriteString("        <link href=\"" + e.LayoutConfig.BaseURL + "/posts/" + templateData.FilenameWithoutExtension + ".html\"/>\n")
			buffer.WriteString("        <id>" + e.LayoutConfig.BaseURL + "/" + templateData.FilenameWithoutExtension + ".html</id>\n")
			buffer.WriteString("        <updated>" + time.Unix(templateData.Date, 0).Format(time.RFC3339) + "</updated>\n")
			buffer.WriteString("        <content type=\"html\"><![CDATA[" + string(templateData.Body) + "]]></content>\n")
			buffer.WriteString("    </entry>\n")
		}
	}

	buffer.WriteString("</feed>\n")
	outputFile, err := os.Create(helpers.SiteDataPath + "rendered/feed.atom")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
