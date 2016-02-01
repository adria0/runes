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
	Markdown  string `json:"markdown"` // the entry contents

	Files []string `json:"files"` // array of shas256s
}
