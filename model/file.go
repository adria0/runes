package model

// File is a file attached to a blog entry
type File struct {
	ID       string
	Filename string
	Content  []byte
}
