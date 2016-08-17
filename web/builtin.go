package web

import (
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/gopad/store"
)


func doGETBuiltin(c *gin.Context) {

	id := store.Normalize(c.Param("id"))

	var content string

	tbox, err := rice.FindBox("builtin")
	if err == nil {
		content = tbox.MustString(id + ".md")
	} else {
		bytes, err := ioutil.ReadFile("web/builtin/" + id + ".md")
		if err == nil {
		    content = string(bytes)
		}
	}

	if err != nil {
		dumpError(c, err)
		return
	}

	err = nil
	c.HTML(http.StatusOK, "builtin.tmpl", gin.H{
		"content": content,
		"error":   err,
	})
}

