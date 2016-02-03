package store

type Config struct {
	path string
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

