package model

// Entry is a blog entry
type Entry struct {
	Workspace string
	ID        string
	Markdown  string
	Files     []string
}
