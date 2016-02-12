package store

import (
    "os"
    "log"
    "io/ioutil"
)

var cacheDir = "/tmp/gopad/cache"
var tempDir = "/tmp/gopad/temp"

func InitCache() {
	if err := os.MkdirAll(cacheDir, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	if err := os.MkdirAll(tempDir, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
}

func WriteTemp(filename string, data []byte) (string,error) {
     path := GetTempPath(filename)
    err :=  ioutil.WriteFile(path,data,0644)
    return path, err
}

func GetTempPath(filename string) string {
    return tempDir + "/" + filename
}

func WriteCache(filename string, data []byte) (string, error) {
    path := GetCachePath(filename)
    err :=  ioutil.WriteFile(path,data,0644)
    return path, err
}

func ExistsCache(filename string) bool {
    if _, err := os.Stat(GetCachePath(filename)); err == nil {
        return true
    }
    return false
}

func GetCachePath(filename string) string {
    return cacheDir + "/" + filename
}



