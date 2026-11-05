package engine

import (
	"bytes"
	"cmp"
	"encoding/json/v2"
	"encoding/xml"
	"html/template"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anna-ssg/anna/pkg/parser"
)

type TagRootPageData struct {
	DeepDataMerge DeepDataMerge
	PageURL       template.URL
	PageData      parser.PageData
	TagNames      []string
}

type CollectionRootPageData struct {
	DeepDataMerge   DeepDataMerge
	PageURL         template.URL
	PageData        parser.PageData
	CollectionNames []string
}

func (e *Engine) RenderTags(fileOutPath string, templ *template.Template) {
	var tagsBuffer bytes.Buffer
	tagsBuffer.Grow(16 * 1024)

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

	tagRootPageData := parser.PageData{
		Frontmatter: parser.Frontmatter{Title: "Tags"},
	}

	tagPageData := TagRootPageData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       "tags.html",
		PageData:      tagRootPageData,
		TagNames:      tagNames,
	}

	// Rendering the page displaying all tags
	err := templ.ExecuteTemplate(&tagsBuffer, "all-tags", tagPageData)
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}

	// Flushing 'tags.html' to the disk
	err = os.WriteFile(fileOutPath+"rendered/tags.html", tagsBuffer.Bytes(), 0666)
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	e.DeepDataMerge.Tags = make(map[template.URL]parser.PageData)

	for tag := range e.DeepDataMerge.TagsMap {
		slices.SortFunc(e.DeepDataMerge.TagsMap[tag], func(a, b parser.PageData) int {
			return cmp.Compare(b.Date, a.Date)
		})
		tagString := string(tag)
		tagString, _ = strings.CutPrefix(tagString, "tags/")
		tagString, _ = strings.CutSuffix(tagString, ".html")

		e.DeepDataMerge.Tags[tag] = parser.PageData{
			Frontmatter: parser.Frontmatter{
				Title: tagString,
			},
		}
	}

	// Rendering the subpages with merged tagged posts
	for tag, taggedTemplates := range e.DeepDataMerge.TagsMap {
		wg.Add(1)
		go func(tag template.URL, taggedTemplates []parser.PageData) {
			defer wg.Done()

			e.RenderPage(fileOutPath, tag, templ, "tag-subpage")
		}(tag, taggedTemplates)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (e *Engine) RenderCollections(fileOutPath string, templ *template.Template) {
	var wg sync.WaitGroup

	e.DeepDataMerge.Collections = make(map[template.URL]parser.PageData)

	for collection := range e.DeepDataMerge.CollectionsMap {
		slices.SortFunc(e.DeepDataMerge.CollectionsMap[collection], func(a, b parser.PageData) int {
			return cmp.Compare(b.Date, a.Date)
		})

		if data, ok := e.DeepDataMerge.CollectionsHierarchyData[collection]; ok && data.IsIndexMdUserDefined {
			continue
		}

		collectionString := string(collection)
		collectionString, _ = strings.CutSuffix(collectionString, "index.html")
		collectionString, _ = strings.CutSuffix(collectionString, "/")
		_, collectionString, _ = strings.CutLast(collectionString, "/")

		e.DeepDataMerge.Collections[collection] = parser.PageData{
			Frontmatter: parser.Frontmatter{
				Title: collectionString,
				Date:  time.Now().Format(time.DateOnly),
			},
		}
	}

	for collection, collectionTemplates := range e.DeepDataMerge.CollectionsMap {
		if data, ok := e.DeepDataMerge.CollectionsHierarchyData[collection]; ok && data.IsIndexMdUserDefined {
			continue
		}

		wg.Add(1)
		go func(collection template.URL, collectionTemplates []parser.PageData) {
			defer wg.Done()

			layoutName := e.DeepDataMerge.CollectionsSubPageLayouts[collection]
			if layoutName == "" {
				layoutName = e.DeepDataMerge.Config.DefaultCollectionIndexLayout
			}

			e.RenderPage(fileOutPath, collection, templ, layoutName)
		}(collection, collectionTemplates)
	}

	wg.Wait()
}

func (e *Engine) GenerateJSONIndex(outFilePath string) {
	// This function creates an index of the site for search
	// It extracts data from the e.Pages slice
	// The index.json file is created during every VanillaRender()

	jsonFile, err := os.Create(outFilePath + "rendered/static/index.json")
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
	}()

	// Copying contents from e.Pages to new JsonMerged struct
	jsonIndexTemplate := make(map[template.URL]JSONIndexTemplate)
	for templateURL, PageData := range e.DeepDataMerge.Pages {
		jsonIndexTemplate[templateURL] = JSONIndexTemplate{
			CompleteURL: PageData.CompleteURL,
			Frontmatter: PageData.Frontmatter,
			Tags:        PageData.Frontmatter.Tags,
		}
	}

	e.DeepDataMerge.JSONIndex = jsonIndexTemplate

	// Marshal the contents of jsonMergedData
	jsonMergedMarshaledData, err := json.Marshal(jsonIndexTemplate)
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}

	_, err = jsonFile.Write(jsonMergedMarshaledData)
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
}

func (e *Engine) GenerateSitemap(outFilePath string) {
	var buffer bytes.Buffer
	buffer.Grow(16 * 1024)

	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// Sorting templates by key
	keys := make([]string, 0, len(e.DeepDataMerge.Pages))
	for k := range e.DeepDataMerge.Pages {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	tempTemplates := make(map[template.URL]parser.PageData)
	for _, templateURL := range keys {
		tempTemplates[template.URL(templateURL)] = e.DeepDataMerge.Pages[template.URL(templateURL)]
	}

	e.DeepDataMerge.Pages = tempTemplates

	// Iterate over parsed markdown files
	for _, PageData := range e.DeepDataMerge.Pages {
		url := e.DeepDataMerge.Config.BaseURL + "/" + string(PageData.CompleteURL)
		buffer.WriteString("\t<url>\n")

		buffer.WriteString("\t\t<loc>")
		buffer.WriteString(url)
		buffer.WriteString("</loc>\n")

		buffer.WriteString("\t\t<lastmod>")
		buffer.WriteString(PageData.Frontmatter.Date)
		buffer.WriteString("</lastmod>\n")

		buffer.WriteString("\t</url>\n")
	}
	buffer.WriteString("</urlset>\n")
	// e.SiteDataPath is the DirPath
	outputFile, err := os.Create(outFilePath)
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}

	defer func() {
		err = outputFile.Close()
		if err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
	}()

	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
}

func (e *Engine) GenerateFeed() {
	var buffer bytes.Buffer
	buffer.Grow(200 * 1024)

	buffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n")
	buffer.WriteString("<?xml-stylesheet href=\"/static/styles/feed.xsl\" type=\"text/xsl\"?>\n")
	buffer.WriteString("<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	buffer.WriteString("  <channel>\n")
	buffer.WriteString("   <title>")
	if err := xml.EscapeText(&buffer, []byte(e.DeepDataMerge.Config.SiteTitle)); err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	buffer.WriteString("</title>\n")

	buffer.WriteString("   <link>")
	buffer.WriteString(e.DeepDataMerge.Config.BaseURL)
	buffer.WriteString("/")
	buffer.WriteString("</link>\n")

	buffer.WriteString("   <description>Recent content on ")
	if err := xml.EscapeText(&buffer, []byte(e.DeepDataMerge.Config.SiteTitle)); err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	buffer.WriteString("</description>\n")
	buffer.WriteString("   <language>en-IN</language>\n")
	buffer.WriteString("   <webMaster>")
	if err := xml.EscapeText(&buffer, []byte(e.DeepDataMerge.Config.Author)); err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	buffer.WriteString("</webMaster>\n")
	buffer.WriteString("   <copyright>")
	if err := xml.EscapeText(&buffer, []byte(e.DeepDataMerge.Config.Copyright)); err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	buffer.WriteString("</copyright>\n")

	buffer.WriteString("   <lastBuildDate>")
	buffer.WriteString(time.Now().Format(time.RFC1123Z))
	buffer.WriteString("</lastBuildDate>\n")

	buffer.WriteString("   <atom:link href=\"")
	buffer.WriteString(e.DeepDataMerge.Config.BaseURL)
	buffer.WriteString("/feed.xml\" rel=\"self\" type=\"application/rss+xml\" />\n")

	var posts []parser.PageData
	for _, PageData := range e.DeepDataMerge.Pages {
		if !PageData.Frontmatter.Draft {
			posts = append(posts, PageData)
		}
	}

	// sort by publication date
	slices.SortFunc(posts, func(a, b parser.PageData) int {
		return cmp.Compare(b.Date, a.Date) // assuming Date is Unix timestamp
	})

	// Iterate over sorted posts
	for _, PageData := range posts {
		buffer.WriteString("    <item>\n")
		buffer.WriteString("      <title>")
		if err := xml.EscapeText(&buffer, []byte(PageData.Frontmatter.Title)); err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
		buffer.WriteString("</title>\n")

		buffer.WriteString("      <link>")
		buffer.WriteString(e.DeepDataMerge.Config.BaseURL)
		buffer.WriteString("/")
		buffer.WriteString(string(PageData.CompleteURL))
		buffer.WriteString("</link>\n")

		buffer.WriteString("      <pubDate>")
		buffer.WriteString(time.Unix(PageData.Date, 0).Format(time.RFC1123Z))
		buffer.WriteString("</pubDate>\n")

		buffer.WriteString("      <author>")
		if err := xml.EscapeText(&buffer, []byte(e.DeepDataMerge.Config.Author)); err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
		buffer.WriteString("</author>\n")

		buffer.WriteString("      <guid>")
		buffer.WriteString(e.DeepDataMerge.Config.BaseURL)
		buffer.WriteString("/")
		buffer.WriteString(string(PageData.CompleteURL))
		buffer.WriteString("</guid>\n")

		buffer.WriteString("      <description>")
		if err := xml.EscapeText(&buffer, []byte(PageData.Body)); err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
		buffer.WriteString("</description>\n")
		buffer.WriteString("    </item>\n")
	}

	buffer.WriteString("  </channel>\n")
	buffer.WriteString("</rss>\n")

	outputFile, err := os.Create(e.SiteDataPath + "rendered/feed.xml")
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			e.Logger.Fatal().Err(err).Send()
		}
	}()

	_, err = outputFile.Write(buffer.Bytes())
	if err != nil {
		e.Logger.Fatal().Err(err).Send()
	}
}
