package store

import (
	"log"
    "os"
    "strings"
)

type FileInfos []os.FileInfo

func (fi FileInfos) Len() int {
    return len(fi)
}
func (fi FileInfos) Swap(i, j int) {
    fi[i], fi[j] = fi[j], fi[i]
}

func (fi FileInfos) Less(i, j int) bool {
    return strings.ToLower(fi[i].Name()) > strings.ToLower(fi[j].Name())
}

func assert(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}



