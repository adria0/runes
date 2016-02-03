package model

type File struct {
	ID       string
	Filename string
	Content  []byte
}

type Entry struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
	Markdown  string // the entry contents
}
