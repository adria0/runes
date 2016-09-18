package store

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/adriamb/runes/store/model"
)

const (
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

func existsPath(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// NewEntryStore creates a new entry store
func NewEntryStore(config Config) *EntryStore {

	return &EntryStore{config}
}

// Create a new repository
func (es *EntryStore) Create() error {
	return os.MkdirAll(es.path, 0744)
}

// Open the repository
func (es *EntryStore) Open() error {
	if _, err := os.Stat(es.path); os.IsNotExist(err) {
		return fmt.Errorf("Repository %s does not exist ", es.path)
	}
	return nil
}

func (es *EntryStore) getMarkdownPath(workspace, ID, version string) string {
	if version == "" {
		return es.path + "/" + workspace + "/" + ID + ".md"
	}
	return es.path + "/" + workspace + "/" + ID + "/v/" + version + ".md"
}

// NewID creates a new entry identifier
func (es *EntryStore) NewID() string {

	return time.Now().Format(DateTimeFormat)

}

// StoreEntry adds a new entry
func (es *EntryStore) StoreEntry(entry *model.Entry) error {

	var err error

	mdPath := es.getMarkdownPath(entry.Workspace, entry.ID, "")

	if _, err := os.Stat(mdPath); err == nil {

		versionsPath := es.path + "/" + entry.Workspace + "/" + entry.ID + "/v"

		if !existsPath(versionsPath) {

			if err := os.MkdirAll(versionsPath, 0744); err != nil {

				log.Fatalf("Cannot create folder %v", err)

			}

		}

		now := time.Now().Format(DateTimeFormat)

		err := os.Rename(
			mdPath,
			es.getMarkdownPath(entry.Workspace, entry.ID, now),
		)

		if err != nil {
			return err
		}

	}

	err = ioutil.WriteFile(mdPath, []byte(entry.Markdown), 0644)

	if err != nil {
		return err
	}

	return nil
}

// GetEntry retrieves the specified entry
func (es *EntryStore) GetEntry(workspace, ID, version string) (*model.Entry, error) {

	mdPath := es.getMarkdownPath(workspace, ID, version)

	rawentry, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return nil, err
	}

	entry := &model.Entry{
		Markdown:  string(rawentry),
		Workspace: workspace,
		ID:        ID,
	}

	files, err := ioutil.ReadDir(es.path + "/" + workspace + "/" + ID + "/f")

	entry.Files = make([]string, 0, len(files))

	if err == nil {

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			entry.Files = append(entry.Files, file.Name())

		}
	}

	return entry, nil
}

// DeleteEntry removes the specified entry
func (es *EntryStore) DeleteEntry(workspace, ID string) error {

	mdPath := es.path + "/" + workspace + "/" + ID + ".md"
	folderPath := es.path + "/" + workspace + "/" + ID
	removedMdPath := es.path + "/" + workspace + "/_" + ID + ".md"
	removedFolderPath := es.path + "/" + workspace + "/_" + ID

	err := os.Rename(
		mdPath,
		removedMdPath,
	)

	if err != nil {
		return err
	}

	err = os.Rename(
		folderPath,
		removedFolderPath,
	)

	return err
}

// GetEntryVersions retrieves the versions of an entry
func (es *EntryStore) GetEntryVersions(workspace, ID string) ([]string, error) {

	versionsPath := es.path + "/" + workspace + "/" + ID + "/v"

	fileInfos, err := ioutil.ReadDir(versionsPath)

	versions := []string{}

	if err != nil {
		return versions, nil
	}

	for _, fileInfo := range fileInfos {

		if !fileInfo.IsDir() {

			name := fileInfo.Name()

			version := name[:len(name)-3] // strip .md
			versions = append(versions, version)

		}
	}

	return versions, nil
}

// ListEntries retruns a list of entries
func (es *EntryStore) ListEntries(workspace string) ([]*model.Entry, error) {

	fileInfos, err := ioutil.ReadDir(es.path + "/" + workspace)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortFileInfos(fileInfos))

	entries := make([]*model.Entry, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {

		name := fileInfo.Name()

		if !fileInfo.IsDir() && strings.HasSuffix(name, ".md") &&
			!strings.HasPrefix(name, "_") {

			entry, err := es.GetEntry(workspace, name[:len(name)-3], "")
			if err != nil {
				return nil, err
			}

            if strings.HasPrefix(entry.Markdown,"# ^") {
                entries = append([]*model.Entry { entry }, entries... )
            } else {
                entries = append(entries, entry)
            }
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

func extractTitleFromMarkdown(md string) string {
	endOfLineIndex := strings.Index(md, "\n")

	var title string

	if endOfLineIndex != -1 {

		title = md[:endOfLineIndex]

	} else {
		title = md

	}

	beginOfTextIndex := strings.Index(title, " ")

	if beginOfTextIndex != -1 {
		title = title[beginOfTextIndex+1:]

	}
	return title
}

// SearchEntries entries using a regular expression
func (es *EntryStore) SearchEntries(workspace, expr string) ([]SearchResult, error) {

	rg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	fileInfos, err := ioutil.ReadDir(es.path + "/" + workspace)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortFileInfos(fileInfos))

	results := []SearchResult{}

	for _, fileInfo := range fileInfos {

		name := fileInfo.Name()

		if !fileInfo.IsDir() && !strings.HasPrefix(name, "_") {

			entry, err := es.GetEntry(workspace, name[:len(name)-3], "")
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
				title := extractTitleFromMarkdown(entry.Markdown)
				results = append(results,
					SearchResult{entry.ID, title, matches})
			}
		}
	}
	return results, nil

}

// StoreFile adds a new file
func (es *EntryStore) StoreFile(workspace, ID string, filename string, reader io.Reader) error {

	filesPath := es.path + "/" + workspace + "/" + ID + "/f"

	if !existsPath(filesPath) {

		if err := os.MkdirAll(filesPath, 0744); err != nil {

			log.Fatalf("Cannot create folder %v", err)

		}

	}

	path := filesPath + "/" + filename

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("File already exists")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return err
	}

	return nil
}

// ListWorkspaces  returns the list of existing  workspaces
func (es *EntryStore) ListWorkspaces() ([]string, error) {

	entryDirs, err := ioutil.ReadDir(es.path + "/")

	if err != nil {
		return nil, err
	}

	list := []string{}

	for _, entryDir := range entryDirs {

		if !entryDir.IsDir() {
			continue
		}

		ws := entryDir.Name()
		if !strings.HasPrefix(ws, "_") && !strings.HasPrefix(ws, ".") {
			list = append(list, ws)
		}
	}

	return list, nil
}

// CreateWorkspace creates a new workspace
func (es *EntryStore) CreateWorkspace(ws string) error {

	path := es.path + "/" + ws
	return os.MkdirAll(path, 0744)

}

// DeleteWorkspace removes a workspace
func (es *EntryStore) DeleteWorkspace(ws string) error {

	path := es.path + "/" + ws
	removedPath := es.path + "/_" + ws

	return os.Rename(
		path,
		removedPath,
	)

}

// ListFiles  all files
func (es *EntryStore) ListFiles(workspace string) ([]FoundFile, error) {

	entryDirs, err := ioutil.ReadDir(es.path + "/" + workspace)

	if err != nil {
		return nil, err
	}

	list := []FoundFile{}

	for _, entryDir := range entryDirs {

		if !entryDir.IsDir() {
			continue
		}

		ID := entryDir.Name()

		files, err := ioutil.ReadDir(es.path + "/" + workspace + "/" + ID + "/f")

		if err != nil {
			continue
		}

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			list = append(list, FoundFile{
				Workspace: workspace,
				ID:        ID,
				Filename:  file.Name(),
			})

		}
	}

	return list, nil
}

// FilePath retrieves the full path for a file
func (es *EntryStore) FilePath(workspace, ID, filename string) string {

	path := es.path + "/" + workspace + "/" + ID + "/f/" + filename

	return path
}
