package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/adriamb/gopad/model"
)

const (
	oldentriesPath = "/entries/old/"
	entriesPath    = "/entries/"
	jsonExt        = ".json"
	mdExt          = ".md"

	// DateTimeFormat is the  format used to generate version timestamps
	DateTimeFormat = "20060102150405"
)

var (
	errNotExists = errors.New("File does not exist")
)

// EntryStore is the store for entries
type EntryStore struct {
	Config
}

// NewEntryStore creates a new entry store
func NewEntryStore(config Config) *EntryStore {
	if err := os.MkdirAll(config.path+entriesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	if err := os.MkdirAll(config.path+oldentriesPath, 0744); err != nil {
		log.Fatalf("Cannot create folder %v", err)
	}
	return &EntryStore{config}
}

func (es *EntryStore) add(entry *model.Entry) error {

	filename := entry.ID + "_" + replaceFilenameChars(entry.Title)
	txtPath := es.path + entriesPath + filename + mdExt
	jsonPath := es.path + entriesPath + filename + jsonExt

	if _, err := os.Stat(txtPath); err == nil {
		panic("Oops, cannot override files!")
	}
	if _, err := os.Stat(jsonPath); err == nil {
		panic("Oops, cannot override files!")
	}

	err := ioutil.WriteFile(txtPath, []byte(entry.Markdown), 0644)

	if err != nil {
		return err
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jsonPath, encoded, 0644)

	if err != nil {
		return err
	}

	return nil
}

// NewID creates a new entry identifier
func (es *EntryStore) NewID() string {

	return time.Now().Format(DateTimeFormat)

}

// Store adds a new entry
func (es *EntryStore) Store(entry *model.Entry) error {

	filename, err := es.getFilenameForID(entry.ID)

	if err != nil && err != errNotExists {
		return err
	}

	if err == nil {
		now := time.Now().Format(DateTimeFormat)
		err := os.Rename(
			es.path+entriesPath+filename+mdExt,
			es.path+oldentriesPath+entry.ID+"_"+now+mdExt,
		)
		if err != nil {
			return err
		}
		err = os.Rename(
			es.path+entriesPath+filename+jsonExt,
			es.path+oldentriesPath+entry.ID+"_"+now+jsonExt,
		)
		if err != nil {
			return err
		}
	}

	err = es.add(entry)
	if err != nil {
		return err
	}

	return nil
}

func (es *EntryStore) get(filename string, old bool) (*model.Entry, error) {

	var path string
	if old {
		path = es.path + oldentriesPath
	} else {
		path = es.path + entriesPath
	}

	rawjson, err := ioutil.ReadFile(path + filename + jsonExt)
	if err != nil {
		return nil, err
	}
	var entry model.Entry
	err = json.Unmarshal(rawjson, &entry)
	if err != nil {
		return nil, err
	}

	rawentry, err := ioutil.ReadFile(path + filename + mdExt)
	if err != nil {
		return nil, err
	}
	entry.Markdown = string(rawentry)

	return &entry, nil
}

// Get retrieves the specified entry
func (es *EntryStore) Get(ID string) (*model.Entry, error) {

	filename, err := es.getFilenameForID(ID)
	if err != nil {
		return nil, err
	}

	entry, err := es.get(filename, false)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// GetVersions retrieves the versions of an entry
func (es *EntryStore) GetVersions(ID string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(es.path + oldentriesPath)

	if err != nil {
		return nil, err
	}

	versions := []string{}

	prefix := ID + "_"
	suffix := mdExt
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()
			if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
				versionLength := strings.Index(name[len(prefix):], ".")
				if versionLength == -1 {
					// Bad file?
					continue
				}
				version := name[len(prefix) : len(prefix)+versionLength]
				versions = append(versions, version)
			}
		}
	}

	return versions, nil
}

// GetVersion retrieves an specific version
func (es *EntryStore) GetVersion(ID string, version string) (*model.Entry, error) {

	entry, err := es.get(ID+"_"+version, true)
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

	prefix := ID + "_"
	suffix := mdExt
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()
			if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
				return name[:len(name)-len(suffix)], nil
			}
		}
	}

	return "", errNotExists
}

// List retruns a list of entries
func (es *EntryStore) List() ([]*model.Entry, error) {

	fileInfos, err := ioutil.ReadDir(es.path + entriesPath)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortFileInfos(fileInfos))

	entries := make([]*model.Entry, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), mdExt) {
			name := fileInfo.Name()
			pos := strings.Index(name, ".")
			if pos == -1 {
				return nil, fmt.Errorf("Bad filename %v", name)
			}
			entry, err := es.get(name[:pos], false)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}
	}
	return entries[:], nil
}

// SearchResult is the return type for searches
type SearchResult struct {
	ID      string
	Title   string
	Matches []string
}

// Search entries using a regular expression
func (es *EntryStore) Search(expr string) ([]SearchResult, error) {

	rg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	fileInfos, err := ioutil.ReadDir(es.path + entriesPath)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortFileInfos(fileInfos))

	results := []SearchResult{}

	for _, fileInfo := range fileInfos {

		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), mdExt) {

			name := fileInfo.Name()
			pos := strings.Index(name, ".")
			if pos == -1 {
				return nil, fmt.Errorf("Bad filename %v", name)
			}

			entry, err := es.get(name[:pos], false)
			if err != nil {
				return nil, err
			}

			matches := []string{}
			lines := strings.Split(entry.Markdown, "\n")
			for _, line := range lines {
				if rg.MatchString(line) {
					matches = append(matches, line)
				}
			}

			if len(matches) > 0 {
				results = append(results,
					SearchResult{entry.ID, entry.Title, matches})
			}
		}
	}
	return results, nil
}
