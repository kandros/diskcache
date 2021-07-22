package diskcache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/afero"
)

const dataFolterName = ".diskcache-data"

type DiskCache struct {
	name string
	path string
	fs   afero.Fs
}

func getDataFolderPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user HOME dir")
	}

	return homeDir + "/" + dataFolterName
}

func createFolderIfNotExists(fs afero.Fs, path string) {
	if _, err := fs.Stat(path); os.IsNotExist(err) {
		err := fs.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to crate folder: %s\n %v", path, err)
		}
	} else if err != nil {
		log.Fatalf("failed to crate folder: %s\n %v", path, err.Error())
	}
}

func New(fs afero.Fs, folderName string) DiskCache {
	p := path.Join(getDataFolderPath(), folderName)
	createFolderIfNotExists(fs, p)

	return DiskCache{name: folderName, path: p, fs: fs}
}

func (dc DiskCache) Set(dataName string, value interface{}) {
	filePath := path.Join(dc.path, dataName)
	file, err := dc.fs.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	if err != nil {
		log.Fatalf("Failed to open file %q, %v", dataName, err)
	}

	encoded := encode(value)
	_, err = file.Write(encoded)
	if err != nil {
		log.Fatalf("Failed to write to file %q, %v", dataName, err)
	}
}

func (dc DiskCache) Get(dataName string, value interface{}) error {
	filePath := path.Join(dc.path, dataName)
	file, err := dc.fs.OpenFile(filePath, os.O_RDONLY, os.ModePerm)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%q data does not exist at %q", dataName, filePath)
		}
		log.Fatalf("Failed to open file %q, %v", dataName, err)
	}
	defer file.Close()

	decode(file, value)

	return nil
}

var expiredErr = errors.New("expired")

func (dc DiskCache) GetIfMaxAge(dataName string, value interface{}, maxAge time.Duration) error {
	filePath := path.Join(dc.path, dataName)
	stat, err := dc.fs.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to read stat of file: %s", filePath)
	}

	expired := stat.ModTime().Second() > int(maxAge.Seconds())
	if expired {
		return expiredErr
	}

	file, err := dc.fs.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%q data does not exist at %q", dataName, filePath)
		}
		log.Fatalf("Failed to open file %q, %v", dataName, err)
	}

	decode(file, value)

	return nil
}

func encode(v interface{}) []byte {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(v)
	if err != nil {
		log.Fatalf("Failed to encode gob: %v", err)
	}
	return buf.Bytes()
}

func decode(r io.Reader, v interface{}) {
	err := gob.NewDecoder(r).Decode(v)
	if err != nil {
		log.Fatalf("Failed to decode gob: %v", err)
	}
}
