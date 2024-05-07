package helpers_test

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/acmpesuecc/anna/v2/pkg/helpers"
)

var HelperTestDirPath = "../../test/helpers/"

func TestCopyDirectoryContents(t *testing.T) {
	t.Run("recursively copying directory contents", func(t *testing.T) {
		helper := helpers.Helper{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		helper.CopyDirectoryContents(HelperTestDirPath+"copy_dir/", HelperTestDirPath+"copy_dir/rendered/")

		baseDirFS := os.DirFS(HelperTestDirPath + "copy_dir/input_dir/")
		err := testfuncTraverseDirectory(baseDirFS, t)
		if err != nil {
			t.Error(err)
		}
	})
}

func testfuncTraverseDirectory(baseDirFS fs.FS, t *testing.T) error {
	err := fs.WalkDir(baseDirFS, ".", func(path string, dir fs.DirEntry, err error) error {
		if !dir.IsDir() {
			gotFile, err := os.ReadFile(HelperTestDirPath + "copy_dir/input_dir/" + path)
			if err != nil {
				t.Errorf("%v", err)
			}

			wantFile, err := os.ReadFile(HelperTestDirPath + "copy_dir/rendered/input_dir/" + path)
			if err != nil {
				t.Errorf("%v", err)
			}

			if !slices.Equal(gotFile, wantFile) {
				t.Errorf("The expected and generated files can be found in %s", HelperTestDirPath)
			}

		}
		return nil
	})
	return err
}

func TestCopyFiles(t *testing.T) {
	t.Run("copy file and create nested parent directories", func(t *testing.T) {
		helper := helpers.Helper{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		helper.CopyFiles(HelperTestDirPath+"copy_files/input.txt", HelperTestDirPath+"copy_files/rendered/output.txt")

		gotFile, err := os.ReadFile(HelperTestDirPath + "copy_files/input.txt")
		if err != nil {
			t.Errorf("%v", err)
		}

		wantFile, err := os.ReadFile(HelperTestDirPath + "copy_files/rendered/output.txt")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(gotFile, wantFile) {
			t.Errorf("The expected and generated files can be found in %s", HelperTestDirPath)
		}
	})
}
