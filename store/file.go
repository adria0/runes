package store

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	filesPath      = "/files/"
)

type FileStore struct {
	Config
}

func NewFileStore(config Config) *FileStore {
	if err := os.MkdirAll(config.path+filesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	return &FileStore{config}
}

func (fs *FileStore) Write(file string, reader io.Reader) (string, error) {

	filename := fmt.Sprintf("%v_%v",
		int32(time.Now().Unix()),
		replaceFilenameChars(file),
	)

	f, err := os.Create(fs.path + filesPath + filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return filename, nil

}

func (fs *FileStore) Fullpath(filename string) string {
	return fs.path + filesPath + filename
}

