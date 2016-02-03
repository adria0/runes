package model

type File struct {
    Id string
    Filename string
    Content []byte
}

type Entry struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
	Markdown  string  // the entry contents
}
