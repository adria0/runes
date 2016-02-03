package store

import (
	"encoding/json"
	"fmt"
	"github.com/amassanet/gopad/model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
    oldentriesPath = "/entries/old/"
	entriesPath    = "/entries/"
	filesPath      = "/files/"
	jsonExt        = ".json"
	txtExt         = ".txt"
)

type Config struct {
	path string
}

type EntryStore struct {
	Config
}

type FileStore struct {
	Config
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

	config := Config{path}

	store := Store{
		Entry: NewEntryStore(config),
		File:  NewFileStore(config),
	}

	return &store
}

// Structure
//   entryid is date    yyyymmddhhmmssmm
//   concepts are type  9999conceptname

func NewFileStore(config Config) *FileStore {
	if err := os.MkdirAll(config.path+filesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	return &FileStore{config}
}

func NewEntryStore(config Config) *EntryStore {
	if err := os.MkdirAll(config.path+entriesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	if err := os.MkdirAll(config.path+oldentriesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
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

func (es *EntryStore) add(entry *model.Entry) error {

	filename := replaceFilenameChars(entry.Title) + "_" + entry.ID

	err := ioutil.WriteFile(es.path+entriesPath+filename+txtExt, []byte(entry.Markdown), 0644)

	if err != nil {
		return err
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(es.path+entriesPath+filename+jsonExt, encoded, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (es *EntryStore) Add(entry *model.Entry) error {

	entry.ID = time.Now().Format("20060102150405")

	err := es.add(entry)

	if err != nil {
		return err
	}

	return nil
}

func (es *EntryStore) Update(entry *model.Entry) error {

	filename, err := es.getFilenameForID(entry.ID)
	if err != nil {
		return err
	}

	now := time.Now().Format("20060102150405")
	os.Rename(
		es.path+entriesPath+filename+txtExt,
		es.path+oldentriesPath+filename+"_"+now+txtExt,
	)
	os.Rename(
		es.path+entriesPath+filename+jsonExt,
		es.path+oldentriesPath+filename+"_"+now+jsonExt,
	)

	err = es.add(entry)
	if err != nil {
		return err
	}

	return nil
}

func (es *EntryStore) get(filename string) (*model.Entry, error) {

	rawjson, err := ioutil.ReadFile(es.path + entriesPath + filename + jsonExt)
	if err != nil {
		return nil, err
	}
	var entry model.Entry
	err = json.Unmarshal(rawjson, &entry)
	if err != nil {
		return nil, err
	}

	rawentry, err := ioutil.ReadFile(es.path + entriesPath + filename + txtExt)
	if err != nil {
		return nil, err
	}
	entry.Markdown = string(rawentry)

	return &entry, nil
}

func (es *EntryStore) Get(ID string) (*model.Entry, error) {

	filename, err := es.getFilenameForID(ID)
	if err != nil {
		return nil, err
	}

	entry, err := es.get(filename)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (es *EntryStore) getFilenameForID(ID string) (string, error) {
	fileInfos, err := ioutil.ReadDir(es.path + entriesPath)

	if err != nil {
		return "", err
	}

	suffix := "_" + ID + txtExt
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()
			if strings.HasSuffix(name, suffix) {
				return name[:len(name)-len(txtExt)], nil
			}
		}
	}

	return "", fmt.Errorf("File not found for entry %v", ID)
}

func (es *EntryStore) List() ([]*model.Entry, error) {

	fileInfos, err := ioutil.ReadDir(es.path + entriesPath)

	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), txtExt) {
			name := fileInfo.Name()
			entry, err := es.get(name[:len(name)-len(txtExt)])
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}
	}
	return entries[:], nil
}
