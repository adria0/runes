package web

import (
    "os"
    "io/ioutil"
	"github.com/adriamb/gopad/model"
)

func existsStaticMd(id string) bool {

	mdfile := "web/mdstatic/" + id + ".md"

    if _, err := os.Stat(mdfile); err == nil {
		return true
	}

	return false

}

func getStaticMdEntry(id string) (entry *model.Entry, err error) {

	var content []byte

	mdfile := "web/mdstatic/" + id + ".md"
    content, err = ioutil.ReadFile(mdfile)

	if err != nil {
		return nil, err
	}

	entry = &model.Entry{
		ID:        id,
		Title:     id,
		Timestamp: 0,
		Markdown:  string(content),
	}

	return entry, err

}
