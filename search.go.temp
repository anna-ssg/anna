package main

import (
	"github.com/blevesearch/bleve"
)

type Content struct {
	Title   string
	Content string
}

func indexSiteContent(contents []Content) (bleve.Index, error) {
	// create a mapping
	mapping := bleve.NewIndexMapping()

	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		return nil, err
	}

	for _, content := range contents {
		err = index.Index(content.Title, content)
		if err != nil {
			return nil, err
		}
	}

	return index, nil
}

func searchContent(index bleve.Index, query string) (*bleve.SearchResult, error) {
	query = bleve.NewMatchQuery(query)
	search := bleve.NewSearchRequest(query)

	searchResults, err := index.Search(search)
	if err != nil {
		return nil, err
	}

	return searchResults, nil
}
