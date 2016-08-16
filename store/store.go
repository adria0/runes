package store

// Config uration for application store
type Config struct {
	path string
}

// Store s
type Store struct {
	Entry *EntryStore
}

// NewStore initializes a new store
func NewStore(path string) *Store {

	config := Config{path}

	store := Store{
		Entry: NewEntryStore(config),
	}

	return &store
}

func Normalize(filename string) string {

    // TODO

    return filename
}

