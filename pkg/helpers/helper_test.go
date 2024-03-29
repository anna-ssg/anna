package helpers_test

import (
	"io/fs"
	"log"
	"testing"
	"testing/fstest"
)

type Helper struct {
	ErrorLogger  *log.Logger
	SiteDataPath string
}

func TestCopyDirectoryContents(t *testing.T) {
	t.Run("Testinf copying contents from A to B", func(t *testing.T) {
		fsMock := fstest.MapFS{
			"index.md": {Data: []byte("# Index Page")},
			"about.md": {Data: []byte("# About Page")},
		}

		fs.WalkDir(fsMock, ".", func(path string, d fs.DirEntry, err error) error {
			return nil
		})

	})
}
