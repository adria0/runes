package model

type Entry struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
	Markdown  string `json:"-"` // the entry contents
}
