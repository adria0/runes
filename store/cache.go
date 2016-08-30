package store

import (
	"io/ioutil"
	"log"
    "time"
    "fmt"
	"os"
)

var cacheDir string
var tempDir string

// InitCache Initializes cache
func InitCache(aCacheDir, aTempDir string) {
	cacheDir = aCacheDir
	tempDir = aTempDir
	if err := os.MkdirAll(cacheDir, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	if err := os.MkdirAll(tempDir, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
}

// WriteTemp creates a temporally file
func WriteTemp(data []byte) (string, error) {

	path := GetTempPath(fmt.Sprintf("tmp-%v",time.Now().UnixNano()))
	err := ioutil.WriteFile(path, data, 0644)
	return path, err
}

// GetTempPath gets  the path of a temporally file
func GetTempPath(filename string) string {
	return tempDir + "/" + filename
}

// WriteCache writes a cache file
func WriteCache(filename string, data []byte) (string, error) {
	path := GetCachePath(filename)
	err := ioutil.WriteFile(path, data, 0644)
	return path, err
}

// ExistsCache returns true if the file aleady exists
func ExistsCache(filename string) bool {
	if _, err := os.Stat(GetCachePath(filename)); err == nil {
		return true
	}
	return false
}

// GetCachePath return the path of a file in the cache
func GetCachePath(filename string) string {
	return cacheDir + "/" + filename
}
