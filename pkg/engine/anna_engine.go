package engine

import (
	"bytes"
	"cmp"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anna-ssg/anna/v3/pkg/parser"
)

type TagRootTemplateData struct {
	DeepDataMerge DeepDataMerge
	PageURL       template.URL
	TemplateData  parser.TemplateData
	TagNames      []string
}

type CollectionRootTemplateData struct {
	DeepDataMerge   DeepDataMerge
	PageURL         template.URL
	TemplateData    parser.TemplateData
	CollectionNames []string
}

func (e *Engine) RenderTags(fileOutPath string, templ *template.Template) {
	var tagsBuffer bytes.Buffer

	// Extracting tag titles
	tags := make([]template.URL, 0, len(e.DeepDataMerge.TagsMap))
	for tag := range e.DeepDataMerge.TagsMap {
		tags = append(tags, tag)
	}

	slices.SortFunc(tags, func(a, b template.URL) int {
		return cmp.Compare(strings.ToLower(string(a)), strings.ToLower(string(b)))
	})

	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagString := string(tag)
		tagString, _ = strings.CutPrefix(tagString, "tags/")
		tagString, _ = strings.CutSuffix(tagString, ".html")

		tagNames = append(tagNames, tagString)
	}

	tagRootTemplataData := parser.TemplateData{
		Frontmatter: parser.Frontmatter{Title: "Tags"},
	}

	tagTemplateData := TagRootTemplateData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       "tags.html",
		TemplateData:  tagRootTemplataData,
		TagNames:      tagNames,
	}

	// Rendering the page displaying all tags
	err := templ.ExecuteTemplate(&tagsBuffer, "all-tags", tagTemplateData)
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

	e.DeepDataMerge.Tags = make(map[template.URL]parser.TemplateData)

	for tag := range e.DeepDataMerge.TagsMap {
		slices.SortFunc(e.DeepDataMerge.TagsMap[tag], func(a, b parser.TemplateData) int {
			return cmp.Compare(b.Date, a.Date)
		})
		tagString := string(tag)
		tagString, _ = strings.CutPrefix(tagString, "tags/")
		tagString, _ = strings.CutSuffix(tagString, ".html")

		e.DeepDataMerge.Tags[tag] = parser.TemplateData{
			Frontmatter: parser.Frontmatter{
				Title: tagString,
			},
		}
	}

	// Rendering the subpages with merged tagged posts
	for tag, taggedTemplates := range e.DeepDataMerge.TagsMap {
		wg.Add(1)
		go func(tag template.URL, taggedTemplates []parser.TemplateData) {
			defer wg.Done()

			e.RenderPage(fileOutPath, tag, templ, "tag-subpage")
		}(tag, taggedTemplates)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (e *Engine) RenderCollections(fileOutPath string, templ *template.Template) {
	var collectionsBuffer bytes.Buffer

	// Extracting collection titles
	collections := make([]template.URL, 0, len(e.DeepDataMerge.CollectionsMap))
	for collection := range e.DeepDataMerge.CollectionsMap {
		collections = append(collections, collection)
	}

	slices.SortFunc(collections, func(a, b template.URL) int {
		return cmp.Compare(strings.ToLower(string(a)), strings.ToLower(string(b)))
	})

	collectionNames := make([]string, 0, len(collections))
	for _, collection := range collections {
		collectionString := string(collection)
		collectionString, _ = strings.CutPrefix(collectionString, "collections/")
		collectionString, _ = strings.CutSuffix(collectionString, ".html")

		collectionNames = append(collectionNames, string(collectionString))
	}

	collectionRootTemplataData := parser.TemplateData{
		Frontmatter: parser.Frontmatter{Title: "Collections"},
	}

	collectionTemplateData := CollectionRootTemplateData{
		DeepDataMerge:   e.DeepDataMerge,
		PageURL:         "collections.html",
		TemplateData:    collectionRootTemplataData,
		CollectionNames: collectionNames,
	}

	// Rendering the page displaying all collections
	err := templ.ExecuteTemplate(&collectionsBuffer, "all-collections", collectionTemplateData)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Flushing 'collections.html' to the disk
	err = os.WriteFile(fileOutPath+"rendered/collections.html", collectionsBuffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	e.DeepDataMerge.Collections = make(map[template.URL]parser.TemplateData)

	for collection := range e.DeepDataMerge.CollectionsMap {
		slices.SortFunc(e.DeepDataMerge.CollectionsMap[collection], func(a, b parser.TemplateData) int {
			return cmp.Compare(b.Date, a.Date)
		})

		collectionString := string(collection)
		collectionString, _ = strings.CutPrefix(collectionString, "collections/")
		collectionString, _ = strings.CutSuffix(collectionString, ".html")

		e.DeepDataMerge.Collections[collection] = parser.TemplateData{
			Frontmatter: parser.Frontmatter{
				Title: collectionString,
			},
		}
	}

	// Rendering the subpages with merged tagged posts
	for collection, collectionTemplates := range e.DeepDataMerge.CollectionsMap {
		wg.Add(1)
		go func(collection template.URL, collectionTemplates []parser.TemplateData) {
			defer wg.Done()

			layoutName := e.DeepDataMerge.CollectionsSubPageLayouts[collection]
			if layoutName == "" {
				layoutName = "collection-subpage"
			}

			e.RenderPage(fileOutPath, collection, templ, layoutName)
		}(collection, collectionTemplates)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (e *Engine) GenerateJSONIndex(outFilePath string) {
	// This function creates an index of the site for search
	// It extracts data from the e.Templates slice
	// The index.json file is created during every VanillaRender()

	jsonFile, err := os.Create(outFilePath + "rendered/static/index.json")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}()

	// Copying contents from e.Templates to new JsonMerged struct
	jsonIndexTemplate := make(map[template.URL]JSONIndexTemplate)
	for templateURL, templateData := range e.DeepDataMerge.Templates {
		jsonIndexTemplate[templateURL] = JSONIndexTemplate{
			CompleteURL: templateData.CompleteURL,
			Frontmatter: templateData.Frontmatter,
			Tags:        templateData.Frontmatter.Tags,
		}
	}

	e.DeepDataMerge.JSONIndex = jsonIndexTemplate

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
	keys := make([]string, 0, len(e.DeepDataMerge.Templates))
	for k := range e.DeepDataMerge.Templates {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	tempTemplates := make(map[template.URL]parser.TemplateData)
	for _, templateURL := range keys {
		tempTemplates[template.URL(templateURL)] = e.DeepDataMerge.Templates[template.URL(templateURL)]
	}

	e.DeepDataMerge.Templates = tempTemplates

	// Iterate over parsed markdown files
	for _, templateData := range e.DeepDataMerge.Templates {
		url := e.DeepDataMerge.LayoutConfig.BaseURL + "/" + string(templateData.CompleteURL)
		buffer.WriteString("\t<url>\n")
		buffer.WriteString("\t\t<loc>" + url + "</loc>\n")
		buffer.WriteString("\t\t<lastmod>" + templateData.Frontmatter.Date + "</lastmod>\n")
		buffer.WriteString("\t</url>\n")
	}
	buffer.WriteString("</urlset>\n")
	// e.SiteDataPath is the DirPath
	outputFile, err := os.Create(outFilePath)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}

	defer func() {
		err = outputFile.Close()
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}()

	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}

func (e *Engine) GenerateFeed() {
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n")
	buffer.WriteString("<?xml-stylesheet href=\"/static/styles/feed.xsl\" type=\"text/xsl\"?>\n")
	buffer.WriteString("<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("  <channel>\n")
	buffer.WriteString("   <title>")
	xml.EscapeText(&buffer, []byte(e.DeepDataMerge.LayoutConfig.SiteTitle))
	buffer.WriteString("</title>\n")
	buffer.WriteString("   <link>" + e.DeepDataMerge.LayoutConfig.BaseURL + "/" + "</link>\n")
	buffer.WriteString("   <description>Recent content on ")
	xml.EscapeText(&buffer, []byte(e.DeepDataMerge.LayoutConfig.SiteTitle))
	buffer.WriteString("</description>\n")
	buffer.WriteString("   <language>en-IN</language>\n")
	buffer.WriteString("   <webMaster>")
	xml.EscapeText(&buffer, []byte(e.DeepDataMerge.LayoutConfig.Author))
	buffer.WriteString("</webMaster>\n")
	buffer.WriteString("   <copyright>")
	xml.EscapeText(&buffer, []byte(e.DeepDataMerge.LayoutConfig.Copyright))
	buffer.WriteString("</copyright>\n")
	buffer.WriteString("   <lastBuildDate>" + time.Now().Format(time.RFC1123Z) + "</lastBuildDate>\n")
	buffer.WriteString("   <atom:link href=\"" + e.DeepDataMerge.LayoutConfig.BaseURL + "/feed.xml\" rel=\"self\" type=\"application/rss+xml\" />\n")

	var posts []parser.TemplateData
	for _, templateData := range e.DeepDataMerge.Templates {
		if !templateData.Frontmatter.Draft {
			posts = append(posts, templateData)
		}
	}

	// sort by publication date
	slices.SortFunc(posts, func(a, b parser.TemplateData) int {
		return cmp.Compare(b.Date, a.Date) // assuming Date is Unix timestamp
	})

	// Iterate over sorted posts
	for _, templateData := range posts {
		buffer.WriteString("    <item>\n")
		buffer.WriteString("      <title>")
		xml.EscapeText(&buffer, []byte(templateData.Frontmatter.Title))
		buffer.WriteString("</title>\n")
		buffer.WriteString("      <link>" + e.DeepDataMerge.LayoutConfig.BaseURL + "/" + string(templateData.CompleteURL) + "</link>\n")
		buffer.WriteString("      <pubDate>" + time.Unix(templateData.Date, 0).Format(time.RFC1123Z) + "</pubDate>\n")
		buffer.WriteString("      <author>")
		xml.EscapeText(&buffer, []byte(e.DeepDataMerge.LayoutConfig.Author))
		buffer.WriteString("</author>\n")
		buffer.WriteString("      <guid>" + e.DeepDataMerge.LayoutConfig.BaseURL + "/" + string(templateData.CompleteURL) + "</guid>\n")
		buffer.WriteString("      <description>")
		xml.EscapeText(&buffer, []byte(templateData.Body))
		buffer.WriteString("</description>\n")
		buffer.WriteString("    </item>\n")
	}

	buffer.WriteString("  </channel>\n")
	buffer.WriteString("</rss>\n")

	outputFile, err := os.Create(e.SiteDataPath + "rendered/feed.xml")
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}()

	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
