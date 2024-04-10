package engine

import (
	"bytes"
	"cmp"
	"encoding/json"
	"html/template"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
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

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Rendering the subpages with merged tagged posts
	for tag, taggedTemplates := range e.TagsMap {
		wg.Add(1)
		go func(tag string, taggedTemplates []parser.TemplateData) {
			defer wg.Done()

			pagePath := "tags/" + tag
			templateData := parser.TemplateData{
				FilenameWithoutExtension: tag,
				Layout:                   e.LayoutConfig,
				Frontmatter: parser.Frontmatter{
					Title: tag,
				},
				SpecificTagTemplates: taggedTemplates,
			}

			e.RenderPage(fileOutPath, template.URL(pagePath), templateData, templ, "tag-subpage")
		}(tag, taggedTemplates)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (e *Engine) GenerateJSONIndex(outFilePath string) {
	// This function creates an index of the site for search
	// It extracts data from the e.Templates slice
	// The index.json file is created during every VanillaRender()

	jsonFile, err := os.Create(outFilePath + "/static/index.json")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer jsonFile.Close()

	// Copying contents from e.Templates to new JsonMerged struct
	jsonIndexTemplate := make(map[template.URL]JSONIndexTemplate)
	for templateURL, templateData := range e.Templates {
		jsonIndexTemplate[templateURL] = JSONIndexTemplate{
			CompleteURL:              templateData.CompleteURL,
			FilenameWithoutExtension: templateData.FilenameWithoutExtension,
			Frontmatter:              templateData.Frontmatter,
			Tags:                     templateData.Frontmatter.Tags,
		}
	}

	e.JSONIndex = jsonIndexTemplate

	// Marshal the contents of jsonMergedData
	jsonMergedMarshaledData, err := json.Marshal(jsonIndexTemplate)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	_, err = jsonFile.Write(jsonMergedMarshaledData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) GenerateSitemap(outFilePath string) {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// Sorting templates by key
	keys := make([]string, 0, len(e.Templates))
	for k := range e.Templates {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	tempTemplates := make(map[template.URL]parser.TemplateData)
	for _, templateURL := range keys {
		tempTemplates[template.URL(templateURL)] = e.Templates[template.URL(templateURL)]
	}

	e.Templates = tempTemplates

	// Iterate over parsed markdown files
	for _, templateData := range e.Templates {
		url := e.LayoutConfig.BaseURL + "/" + templateData.FilenameWithoutExtension + ".html"
		buffer.WriteString("\t<url>\n")
		buffer.WriteString("\t\t<loc>" + url + "</loc>\n")
		buffer.WriteString("\t\t<lastmod>" + templateData.Frontmatter.Date + "</lastmod>\n")
		buffer.WriteString("\t</url>\n")
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
	buffer.WriteString("<?xml-stylesheet href=\"/static/styles/feed.xsl\" type=\"text/xsl\"?>\n")
	buffer.WriteString("<feed xmlns=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("    <title>" + e.LayoutConfig.SiteTitle + "</title>\n")
	buffer.WriteString("    <link href=\"" + e.LayoutConfig.BaseURL + "/" + "\" rel=\"self\"/>\n")
	buffer.WriteString("    <updated>" + time.Now().Format(time.RFC3339) + "</updated>\n")

	// iterate over parsed markdown files that are non-draft posts
	for _, templateData := range e.Templates {
		if !templateData.Frontmatter.Draft {
			buffer.WriteString("<entry>\n")
			buffer.WriteString("        <title>" + templateData.Frontmatter.Title + "</title>\n")
			buffer.WriteString("        <link href=\"" + e.LayoutConfig.BaseURL + "/posts/" + templateData.FilenameWithoutExtension + ".html\"/>\n")
			buffer.WriteString("        <id>" + e.LayoutConfig.BaseURL + "/posts/" + templateData.FilenameWithoutExtension + ".html</id>\n")
			buffer.WriteString("        <updated>" + time.Unix(templateData.Date, 0).Format(time.RFC3339) + "</updated>\n")
			buffer.WriteString("        <content type=\"html\"><![CDATA[" + string(templateData.Body) + "]]></content>\n")
			buffer.WriteString("    </entry>\n")
		}
	}

	buffer.WriteString("</feed>\n")
	outputFile, err := os.Create(helpers.SiteDataPath + "rendered/feed.xml")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
