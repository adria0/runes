package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/runes/web/render"
)

func doGETBuiltin(c *gin.Context) {

	id := normalize(c.Param("id"))

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
		"htmlHeaders": render.HTMLHeaders(),
		"content":     content,
		"error":       err,
	})
}
