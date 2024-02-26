package ssg

import (
	"html/template"
	"io"
	"os"
	"strings"
)

// Copies the contents of the dirPath directory to outDirPath
func (g *Generator) copyDirectoryContents(dirPath string, outDirPath string) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	err = os.MkdirAll(outDirPath, 0750)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}

	// Copying the contents of the dirPath directory
	for _, entry := range dirEntries {
		if entry.IsDir() {
			g.copyDirectoryContents(dirPath+entry.Name()+"/", outDirPath+entry.Name()+"/")
		} else {
			g.copyFiles(dirPath+entry.Name(), outDirPath+entry.Name())
		}
	}
}

func (g *Generator) copyFiles(srcPath string, destPath string) {
	source, err := os.Open(srcPath)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		g.ErrorLogger.Fatal(err)
	}
}

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
			return
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		g.mdFilesName = append(g.mdFilesName, entry.Name())

		filepath := dirPath + entry.Name()
		g.mdFilesPath = append(g.mdFilesPath, filepath)

		// Parsing titles of md files in the posts folder
		if dirPath == "content/posts/" {
			g.mdPosts = append(g.mdPosts, (strings.Split(entry.Name(), ".")[0]))
		}
	}

	// Reading the markdown files into memory
	for i, filepath := range g.mdFilesPath {
		content, err := os.ReadFile(filepath)
		if err != nil {
			g.ErrorLogger.Fatal(err)
		}

		frontmatter, body := g.parseMarkdownContent(string(content))

		page := Page{
            Filename:    strings.Split(g.mdFilesName[i], ".")[0],
			Frontmatter: frontmatter,
			Body:        template.HTML(body),
			Layout:      g.LayoutConfig,
		}

        g.mdParsed = append(g.mdParsed, page) 
    }
}
