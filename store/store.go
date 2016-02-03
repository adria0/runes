package store

import (
    "fmt"
	"encoding/json"
	"github.com/amassanet/gopad/model"
	"os"
    "io/ioutil"
	"time"
    "log"
    "io"
    "strings"
)

const (
    kOldEntriesPath = "/entries/old/"
    kEntriesPath = "/entries/"
    kFilesPath = "/files/"
    kJsonExt = ".json"
    kTxtExt = ".txt"
)

type StoreConfig struct {
	path string
}

type EntryStore struct {
    StoreConfig
}

type FileStore struct {
	StoreConfig
}

type Store struct {
	Entry *EntryStore
	File  *FileStore
}

var filenameRunes map[rune]rune

func NewStore(path string) *Store {
 	search := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.ÀÁÈÉÌÍÒÓÙÚàáèéìíòóùúÑñ")
	replace := []rune("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz0123456789-.aaeeiioouuaaeeiioouuNn")
	filenameRunes = make(map[rune]rune)

    for i := range search {
		filenameRunes[search[i]] = replace[i]
	}

    config := StoreConfig{path}

	store := Store{
        Entry: NewEntryStore(config),
		File:  NewFileStore(config),
	}

	return &store
}

// Structure
//   entryid is date    yyyymmddhhmmssmm
//   concepts are type  9999conceptname

func NewFileStore(config StoreConfig) *FileStore {
   if err := os.MkdirAll(config.path+kFilesPath,0744) ; err!= nil {
        log.Fatalf("Cannot create folder %v",err)
    }
    return &FileStore{config}
}

func NewEntryStore(config StoreConfig) *EntryStore {
    if err := os.MkdirAll(config.path+kEntriesPath ,0744) ; err!= nil {
        log.Fatalf("Cannot create folder %v",err)
    }
    if err := os.MkdirAll(config.path+kOldEntriesPath ,0744) ; err!= nil {
        log.Fatalf("Cannot create folder %v",err)
    }
    return &EntryStore{config}
}

func replaceFilenameChars(s string) string {
	r := []rune(s)
	for i := range r {
		if k, exist := filenameRunes[r[i]]; exist {
			r[i] = k

		} else {
			r[i] = '_'

		}
	}
	return string(r)
}

func (fs *FileStore) Write(file string, reader io.Reader ) (string, error) {

        filename := fmt.Sprintf("%v_%v",
            int32(time.Now().Unix()),
            replaceFilenameChars(file),
        )

        f, err := os.Create(fs.path+kFilesPath+filename)
 		if err != nil {
			return "",err
		}

        defer f.Close()

        if _, err := io.Copy(f,reader) ; err != nil {
            return "",err
        }

        return filename,nil

}

func (es *FileStore) Fullpath(filename string) string{

    return es.path+kFilesPath+filename
}

func (es *EntryStore) add(entry *model.Entry) error {

    filename := replaceFilenameChars(entry.Title) + "_" + entry.Id

    err := ioutil.WriteFile(es.path+kEntriesPath+filename+kTxtExt, []byte(entry.Markdown), 0644)

    if err != nil {
        return err
    }

    encoded, err := json.Marshal(entry)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(es.path+kEntriesPath+filename+kJsonExt, encoded, 0644)

    if err != nil {
        return err
    }

    return nil
}


func (es *EntryStore) Add(entry *model.Entry) error {

		entry.Id = time.Now().Format("20060102150405")

        err := es.add(entry)

        if err != nil {
			return err
		}

	return nil
}

func (es *EntryStore) Update(entry *model.Entry) error{

        filename,err := es.getFilenameForId (entry.Id)
		if err != nil {
			return err
		}

        now := time.Now().Format("20060102150405")
        os.Rename(
            es.path+kEntriesPath+filename+kTxtExt,
            es.path+kOldEntriesPath+filename+"_"+now+kTxtExt,
        )
        os.Rename(
            es.path+kEntriesPath+filename+kJsonExt,
            es.path+kOldEntriesPath+filename+"_"+now+kJsonExt,
        )


        err = es.add(entry)
		if err != nil {
			return err
		}

	return nil
}

func (es *EntryStore) get(filename string) (*model.Entry,error) {

    rawjson, err := ioutil.ReadFile(es.path + kEntriesPath + filename+kJsonExt)
    if err != nil {
        return nil,err
    }
    var entry model.Entry
    err = json.Unmarshal(rawjson, &entry)
    if err != nil {
        return nil,err
    }

    rawentry, err := ioutil.ReadFile(es.path + kEntriesPath  + filename+kTxtExt)
    if err != nil {
        return nil,err
    }
    entry.Markdown = string(rawentry)

    return &entry,nil
}

func (es *EntryStore) Get(Id string) (*model.Entry, error){


        filename,err := es.getFilenameForId(Id)
		if err != nil {
			return nil,err
		}

		entry, err := es.get(filename)
		if err != nil {
			return nil,err
		}

        return entry,nil
}

func (es *EntryStore) getFilenameForId(Id string) (string,error) {
    fileInfos, err := ioutil.ReadDir(es.path+kEntriesPath)

	if err != nil {
		return "", err
	}

    suffix := "_"+Id+kTxtExt
    for _, fileInfo := range fileInfos {
        if !fileInfo.IsDir() {
            name := fileInfo.Name()
            if strings.HasSuffix(name,suffix) {
                return name[:len(name)-len(kTxtExt)], nil
            }
        }
    }

    return "", fmt.Errorf("File not found for entry %v",Id)
}

func (es *EntryStore) List() ([]*model.Entry, error) {

    fileInfos, err := ioutil.ReadDir(es.path+kEntriesPath)

		if err != nil {
			return nil, err
		}

		entries := make([]*model.Entry, 0, len(fileInfos))

		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(),kTxtExt) {
                name := fileInfo.Name()
                entry, err := es.get(name[:len(name)-len(kTxtExt)])
				if err != nil {
                    return nil, err
				}
				entries = append(entries, entry)
			}
		}
        return entries[:], nil
}
