package diskcache

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestGetDataFolderPath(t *testing.T) {
	is := is.New(t)
	os.Setenv("HOME", "my/home/path")
	s := getDataFolderPath()
	is.Equal(s, "my/home/path/.diskcache-data")
}

func TestCreateFolderIfNotExists(t *testing.T) {
	path := "my-folder/my-child-folder"
	fs := afero.NewMemMapFs()
	is := is.New(t)
	folderExists, err := afero.DirExists(fs, path)
	is.NoErr(err)
	is.Equal(folderExists, false)

	createFolderIfNotExists(fs, path)
	folderExists, err = afero.DirExists(fs, path)
	is.NoErr(err)
	is.True(folderExists)

	// _, err := appFS.Stat(name)

	// // create test files and directories
	// appFS.MkdirAll("src/a", 0755)
	// afero.WriteFile(appFS, "src/a/b", []byte("file b"), 0644)
	// afero.WriteFile(appFS, "src/c", []byte("file c"), 0644)
	// name := "src/c"
	// _, err := appFS.Stat(name)
	// if os.IsNotExist(err) {
	// 	t.Errorf("file \"%s\" does not exist.\n", name)
	// }
}

func TestNew(t *testing.T) {
	is := is.New(t)
	os.Setenv("HOME", "my/home/path")
	fs := afero.NewMemMapFs()
	New(fs, "gino")
	folderExists, err := afero.DirExists(fs, "my/home/path/.diskcache-data/gino")
	is.NoErr(err)
	is.Equal(folderExists, true)

	// createFolderIfNotExists(fs, path)
	// folderExists, err = afero.DirExists(fs, path)
	// is.NoErr(err)
	// is.True(folderExists)
}
