package diskcache

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func getDiskcache() DiskCache {
	os.Setenv("HOME", "my/home/path")
	fs := afero.NewMemMapFs()
	return New(fs, "test-dir")
}

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
}

func TestNew(t *testing.T) {
	dc := getDiskcache()
	folderExists, err := afero.DirExists(dc.fs, "my/home/path/.diskcache-data/test-dir")
	is := is.New(t)
	is.NoErr(err)
	is.Equal(folderExists, true)
}

func TestSet(t *testing.T) {
	dc := getDiskcache()
	key := "abc"
	dc.Set(key, "my content")

	ok, err := afero.Exists(dc.fs, path.Join(dc.path, key))
	is := is.New(t)
	is.NoErr(err)
	is.True(ok) // file "abc" doesn't exists
}

func TestSetOverride(t *testing.T) {
	dc := getDiskcache()
	dataName := "mykey"
	dc.Set(dataName, []string{"my", "data"})
	dc.Set(dataName, "something else")

	var s string
	err := dc.Get(dataName, &s)
	is := is.New(t)
	is.NoErr(err)
	is.Equal(s, "something else")
}

func TestGetThanDontExists(t *testing.T) {
	dc := getDiskcache()
	dataName := "non_existant_key"
	err := dc.Get(dataName, nil)
	is := is.New(t)
	filePath := path.Join(dc.path, dataName)
	is.Equal(err, dataNotFoundError{dataName: "non_existant_key", filePath: filePath})
}

func TestGet(t *testing.T) {
	dc := getDiskcache()
	dataName := "mykey"
	in := []string{"my", "data"}
	dc.Set(dataName, in)

	out := []string{}
	err := dc.Get(dataName, &out)
	is := is.New(t)
	is.NoErr(err)
	is.Equal(out, in)
}

func TestGetIfMaxAgeNotExpired(t *testing.T) {
	/* This test is brittle, if it failes is because afero memory FS does not correctly set modtime to now for new files */
	dc := getDiskcache()

	dataName := "not_expired_yet"
	in := []string{"my", "data"}
	dc.Set(dataName, in)

	is := is.New(t)
	out := []string{}
	err := dc.GetIfMaxAge(dataName, &out, 100*time.Second)
	is.NoErr(err)
	is.Equal(out, in)
}

func TestGetIfMaxAgeExpired(t *testing.T) {
	dc := getDiskcache()
	dataName := "mykey"
	in := []string{"my", "data"}
	dc.Set(dataName, in)

	is := is.New(t)
	out := []string{}
	err := dc.GetIfMaxAge(dataName, &out, 1*time.Millisecond)
	is.Equal(err.Error(), "expired")
	is.Equal(out, []string{})
}
