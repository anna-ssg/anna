package anna

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ValidateHTMLContent() {
	root, err := filepath.Abs(SiteDataPath + "rendered")
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	fmt.Println("Walking directory at path:", root)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".html" {
			// Parse HTML file
			if err := parseHTMLFile(path); err != nil {
				fmt.Printf("Error parsing %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
}

func parseHTMLFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Load the HTML content into a GoQuery document
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}

	// Checking for semantic elements
	semanticElements := []string{"header", "nav", "article", "footer"}
	missingElements := make([]string, 0)

	for _, element := range semanticElements {
		if doc.Find(element).Length() == 0 {
			missingElements = append(missingElements, element)
		}
	}

	if len(missingElements) > 0 {
		fmt.Printf("File %s is missing the following semantic elements: %s\n", path, strings.Join(missingElements, ", "))
	} else {
		fmt.Printf("File %s has all the required semantic elements\n", path)
	}

	return nil
}
