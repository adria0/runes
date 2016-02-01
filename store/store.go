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
)

type StoreResult struct {
	Data interface{}
	Err  error
}

type StoreChannel chan StoreResult

type StoreConfig struct {
	path string
}

type EntryStore struct {
	StoreConfig
}

type FileStore struct {
	StoreConfig
	filenameRunes map[rune]rune
}

type Store struct {
	Entry *EntryStore
	File  *FileStore
}

func NewStore(path string) *Store {

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
 	search := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.ÀÁÈÉÌÍÒÓÙÚàáèéìíòóùúÑñ")
	replace := []rune("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz0123456789-.aaeeiioouuaaeeiioouuNn")
	filenameRunes := make(map[rune]rune)
	for i := range search {
		filenameRunes[search[i]] = replace[i]
	}
    if err := os.MkdirAll(config.path+"/files/",0744) ; err!= nil {
        log.Fatalf("Cannot create folder %v",err)
    }
    return &FileStore{config,filenameRunes}
}

func NewEntryStore(config StoreConfig) *EntryStore {
     if err := os.MkdirAll(config.path+"/entries/",0744) ; err!= nil {
        log.Fatalf("Cannot create folder %v",err)
    }
    return &EntryStore{config}
}

func (fs *FileStore) replaceFilenameChars(s string) string {
	r := []rune(s)
	for i := range r {
		if k, exist := fs.filenameRunes[r[i]]; exist {
			r[i] = k

		} else {
			r[i] = '_'

		}
	}
	return string(r)
}

func (fs *FileStore) Write(filename string, reader io.Reader ) StoreChannel {
	ch := make(StoreChannel)

	go func() {

        filename := fmt.Sprintf("%v_%v",
            int32(time.Now().Unix()),
            fs.replaceFilenameChars(filename),
        )

        f, err := os.Create(fs.path+"/files/"+filename)
 		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

        defer f.Close()

        if _, err := io.Copy(f,reader) ; err != nil {
            forwardErrorAndClose(ch, err)
            return
        }

		sendSuccessAndClose(ch, filename)
	}()

	return ch
}

func (es *FileStore) Fullpath(filename string) StoreChannel {

    ch := make(StoreChannel)

	go func() {

        fullpath := es.path+"/files/"+filename

		sendSuccessAndClose(ch, fullpath)
	}()

	return ch
}

func (es *EntryStore) Add(entry *model.Entry) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		entry.Id = time.Now().Format("20060102150405")

		encoded, err := json.Marshal(entry)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		err = ioutil.WriteFile(es.path+"/entries/"+entry.Id, encoded, 0644)

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		sendSuccessAndClose(ch, nil)
	}()

	return ch
}

func (es *EntryStore) Update(entry *model.Entry) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		encoded, err := json.Marshal(entry)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		err = ioutil.WriteFile(es.path+"/entries/"+entry.Id, encoded, 0644)

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		sendSuccessAndClose(ch, nil)
	}()

	return ch
}

func (es *EntryStore) Get(Id string) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		raw, err := ioutil.ReadFile(es.path + "/entries/" + Id)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}
		var entry model.Entry
		err = json.Unmarshal(raw, &entry)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		sendSuccessAndClose(ch, &entry)
	}()

	return ch
}

func (es *EntryStore) List() StoreChannel {

	ch := make(StoreChannel)

	go func() {

		fileInfos, err := ioutil.ReadDir(es.path+"/entries/")

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		entries := make([]*model.Entry, 0, len(fileInfos))

		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				raw, err := ioutil.ReadFile(es.path + "/entries/" + fileInfo.Name())
				if err != nil {
					forwardErrorAndClose(ch, err)
					return
				}
				var entry model.Entry
				err = json.Unmarshal(raw, &entry)
				if err != nil {
					forwardErrorAndClose(ch, err)
					return
				}
				entries = append(entries, &entry)
			}
		}

		sendSuccessAndClose(ch, entries[:])
	}()

	return ch
}
