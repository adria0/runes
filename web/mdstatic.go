package web

import (
	"io/ioutil"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/gopad/model"
)

func existsStaticMd(id string) bool {

	tbox, err := rice.FindBox("mdstatic")
	if err == nil {
		_, err = tbox.String(id + ".md")
		return err == nil
	}

	if _, err := os.Stat("web/mdstatic/" + id + ".md"); err == nil {
		return true
	}
	return false

}

func getStaticMdEntry(id string) (entry *model.Entry, err error) {

	var content string

	tbox, err := rice.FindBox("mdstatic")
	if err == nil {
		content = tbox.MustString(id + ".md")
	} else {
		bytes, err := ioutil.ReadFile("web/mdstatic/" + id + ".md")
		if err != nil {
			return nil, err
		}
		content = string(bytes)
	}

	entry = &model.Entry{
		ID:        id,
		Title:     id,
		Timestamp: 0,
		Markdown:  content,
	}

	return entry, err

}
