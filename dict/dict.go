package dict

import (
	"github.com/adriamb/gopad/store"
	"strings"
	"unicode"
)

type Dict struct {
	defs  map[string]string
	store *store.Store
}

func New(store *store.Store) *Dict {
	return &Dict{nil, store}
}

func (d *Dict) Rebuild() error {
	d.defs = nil
	return d.build()
}

func (d *Dict) Defs() (map[string]string, error) {

	if d.defs != nil {
		return d.defs, nil
	}

	err := d.build()
	if err != nil {
		return nil, err
	}

	return d.defs, nil

}

func (d *Dict) build() error {

	entries, err := d.store.Entry.List()
	if err != nil {
		return err
	}

	d.defs = make(map[string]string)

	for _, entry := range entries {

		lines := strings.Split(entry.Markdown, "\n")
		word := ""
		description := ""

		for i := 0; i < len(lines); i++ {

			line := strings.TrimSpace(lines[i])
			lineRunes := []rune(line)

			// pending definition
			if len(word) > 0 {

				if len(line) > 0 && !unicode.IsPunct(lineRunes[0]) {
					// more content
					description = description + line
					continue
				}

				d.defs[word] = description
				word = ""
				description = ""
			}

			// check for new definition
			if strings.HasPrefix(line, "-") && strings.Contains(line, "§:") {
				line := line[1:]
				split := strings.SplitN(line, "§:", 2)
				word = strings.TrimSpace(split[0])
				description = strings.TrimSpace(split[1])
			}
		}

		if len(word) > 0 {
			d.defs[word] = description
		}

	}

	return nil

}
