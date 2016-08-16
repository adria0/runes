package web

import (
	"io/ioutil"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/gopad/model"
)

func isBuiltinWorkspace(ws string) bool {

    return ws == "builtin"

}

func getBuiltinMdEntry(id string) (entry *model.Entry, err error) {

	var content string

	tbox, err := rice.FindBox("builtin")
	if err == nil {
		content = tbox.MustString(id + ".md")
	} else {
		bytes, err := ioutil.ReadFile("web/builtin/" + id + ".md")
		if err != nil {
			return nil, err
		}
		content = string(bytes)
	}

	entry = &model.Entry{
		ID:        id,
		Markdown:  content,
		Workspace: "mdstatic",
	}

	return entry, err

}
