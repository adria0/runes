package model

// Entry is a blog entry
type Entry struct {
    Workspace string `json:"workspace"`
    ID        string `json:"id"`
    Markdown  string `json:"-"` // the entry contents
}
