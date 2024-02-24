package ssg

import (
	"io"
	"os"
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
