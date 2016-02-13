package store

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

const (
	filesPath = "/files/"
)

var (
	errFileAlreadyExists = errors.New("File already exists")
)

// FileStore is the store for files
type FileStore struct {
	Config
}

// NewFileStore  inirializes a new store
func NewFileStore(config Config) *FileStore {
	if err := os.MkdirAll(config.path+filesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	return &FileStore{config}
}

// Write adds a new file
func (fs *FileStore) Write(file string, entryID string, reader io.Reader) (string, error) {

	filename := fmt.Sprintf("%v_%v", entryID, replaceFilenameChars(file))
	path := fs.path + filesPath + filename

	if _, err := os.Stat(path); err == nil {
		return "", errFileAlreadyExists
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return filename, nil
}

// Fullpath retrieves the full path for a file
func (fs *FileStore) Fullpath(filename string) string {
	return fs.path + filesPath + filename
}

// List  all files
func (fs *FileStore) List() ([]string, error) {

	fileInfos, err := ioutil.ReadDir(fs.path + filesPath)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortFileInfos(fileInfos))

	files := make([]string, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}
	return files, nil
}
