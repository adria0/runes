package store

import (
	"encoding/json"
	"github.com/amassanet/gopad/model"
	"io/ioutil"
	"time"
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

type Store struct {
	Entry *EntryStore
}

func NewStore(path string) *Store {
	config := StoreConfig{
		path: path,
	}
	store := Store{
		Entry: &EntryStore{config},
	}
	return &store
}

// Structure
//   entryid is date    yyyymmddhhmmssmm
//   concepts are type  9999conceptname

func (s *EntryStore) Add(entry *model.Entry) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		entry.Id = time.Now().Format("20060102150405")

		encoded, err := json.Marshal(entry)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		err = ioutil.WriteFile(s.path+"/"+entry.Id, encoded, 0644)

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		sendSuccessAndClose(ch, nil)
	}()

	return ch
}

func (s *EntryStore) Update(entry *model.Entry) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		encoded, err := json.Marshal(entry)
		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		err = ioutil.WriteFile(s.path+"/"+entry.Id, encoded, 0644)

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		sendSuccessAndClose(ch, nil)
	}()

	return ch
}

func (s *EntryStore) Get(Id string) StoreChannel {

	ch := make(StoreChannel)

	go func() {

		raw, err := ioutil.ReadFile(s.path + "/" + Id)
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

func (s *EntryStore) List() StoreChannel {

	ch := make(StoreChannel)

	go func() {

		fileInfos, err := ioutil.ReadDir(s.path)

		if err != nil {
			forwardErrorAndClose(ch, err)
			return
		}

		entries := make([]*model.Entry, 0, len(fileInfos))

		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				raw, err := ioutil.ReadFile(s.path + "/" + fileInfo.Name())
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
