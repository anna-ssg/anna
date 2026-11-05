package helpers_test

import (
	"io"
	"io/fs"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/anna-ssg/anna/pkg/helpers"

	"github.com/rs/zerolog"
)

var HelperTestDirPath = "../../test/helpers/"

func InitLogger(out io.Writer) *zerolog.Logger {
	console := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	}
	logger := zerolog.New(console).With().Timestamp().Caller().Stack().Logger()
	return &logger
}

func TestCopyDirectoryContents(t *testing.T) {
	t.Run("recursively copying directory contents", func(t *testing.T) {
		helper := helpers.Helper{
			Logger: InitLogger(os.Stdout),
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
			Logger: InitLogger(os.Stdout),
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
