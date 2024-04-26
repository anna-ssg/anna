package helpers_test

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/acmpesuecc/anna/pkg/helpers"
)

var HelperTestDirPath = "../../test/helpers/"

func TestCopyDirectoryContents(t *testing.T) {
	t.Run("recursively copying directory contents", func(t *testing.T) {
		helper := helpers.Helper{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		helper.CopyDirectoryContents(HelperTestDirPath+"copy_dir/", HelperTestDirPath+"copy_dir/rendered/")

		baseDirFS := os.DirFS(HelperTestDirPath + "copy_dir/input_dir/")
		TraverseDirectory(baseDirFS, t)
	})
}

func TraverseDirectory(baseDirFS fs.FS, t *testing.T) error {
	fs.WalkDir(baseDirFS, ".", func(path string, dir fs.DirEntry, err error) error {
		if !dir.IsDir() {
			got_file, err := os.ReadFile(HelperTestDirPath + "copy_dir/input_dir/" + path)
			if err != nil {
				t.Errorf("%v", err)
			}

			want_file, err := os.ReadFile(HelperTestDirPath + "copy_dir/rendered/input_dir/" + path)
			if err != nil {
				t.Errorf("%v", err)
			}

			if !slices.Equal(got_file, want_file) {
				t.Errorf("The expected and generated files can be found in %s", HelperTestDirPath)
			}

		}
		return nil
	})
	return nil
}

func TestCopyFiles(t *testing.T) {
	t.Run("copy file and create nested parent directories", func(t *testing.T) {
		helper := helpers.Helper{
			ErrorLogger: log.New(os.Stderr, "TEST ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		}
		helper.CopyFiles(HelperTestDirPath+"copy_files/input.txt", HelperTestDirPath+"copy_files/rendered/output.txt")

		got_file, err := os.ReadFile(HelperTestDirPath + "copy_files/input.txt")
		if err != nil {
			t.Errorf("%v", err)
		}

		want_file, err := os.ReadFile(HelperTestDirPath + "copy_files/rendered/output.txt")
		if err != nil {
			t.Errorf("%v", err)
		}

		if !slices.Equal(got_file, want_file) {
			t.Errorf("The expected and generated files can be found in %s", HelperTestDirPath)
		}
	})
}
