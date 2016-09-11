//go:generate rice embed-go

package web

import (
	"net/http"

	"github.com/adriamb/runes/server/instance"
	"github.com/adriamb/runes/store/model"
	"github.com/gin-gonic/gin"

	"github.com/adriamb/runes/web/render"
)

type tmplEntry struct {
	*model.Entry
}

func doPOSTSearch(c *gin.Context) {

	ws := normalize(c.Param("ws"))
	query := c.Request.FormValue("query")

	results, err := instance.Srv.Store.Entry.SearchEntries(ws, query)
	if err != nil {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"ws":    ws,
			"error": err,
		})
		return
	}
	if len(results) == 0 {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"ws":   ws,
			"info": "No results",
		})
		return
	}
	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"ws":      ws,
		"results": results,
	})
}

func doGETEntries(c *gin.Context) {

	var entries []*model.Entry
	var err error

	ws := normalize(c.Param("ws"))
	id := normalize(c.Param("id"))

	if id != "" {

		var entry *model.Entry

		entry, err = instance.Srv.Store.Entry.GetEntry(ws, id, "")

		if err == nil {
			entries = append(entries, entry)
		}

	} else {

		entries, err = instance.Srv.Store.Entry.ListEntries(ws)

	}

	if err != nil {
		dumpError(c, err)
		return
	}

	htmlEntries := []tmplEntry{}
	for _, entry := range entries {
		htmlEntries = append(htmlEntries, tmplEntry{entry})
	}

	err = nil
	c.HTML(http.StatusOK, "entries.tmpl", gin.H{
		"htmlHeaders": render.HTMLHeaders(),
		"entries":     htmlEntries,
		"ws":          ws,
		"error":       err,
	})
}

func doGETFiles(c *gin.Context) {

	ws := normalize(c.Param("ws"))

	files, err := instance.Srv.Store.Entry.ListFiles(ws)
	if err != nil {
		dumpError(c, err)
		return
	}

	c.HTML(http.StatusOK, "files.tmpl", gin.H{
		"ws":    ws,
		"files": files,
	})
}
