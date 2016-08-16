//go:generate rice embed-go

package web

import (
	"net/http"

	"github.com/adriamb/gopad/model"
	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/store"
	"github.com/gin-gonic/gin"
)

type tmplEntry struct {
	*model.Entry
}

func doPOSTSearch(c *gin.Context) {

	ws := store.Normalize(c.Param("ws"))
	query := c.Request.FormValue("query")

	results, err := server.Srv.Store.Entry.SearchEntries(ws, query)
	if err != nil {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"ws":     ws,
			"prefix": server.Srv.Config.Prefix,
			"error":  err,
		})
		return
	}
	if len(results) == 0 {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"ws":     ws,
			"prefix": server.Srv.Config.Prefix,
			"info":   "No results",
		})
		return
	}
	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"ws":      ws,
		"prefix":  server.Srv.Config.Prefix,
		"results": results,
	})
}

func doGETEntries(c *gin.Context) {

	var entries []*model.Entry
	var err error

	ws := store.Normalize(c.Param("ws"))
	id := store.Normalize(c.Param("id"))

	if id != "" {

		var entry *model.Entry

		entry, err = server.Srv.Store.Entry.GetEntry(ws, id, "")

		if err == nil {
			entries = append(entries, entry)
		}

	} else {

		entries, err = server.Srv.Store.Entry.ListEntries(ws)

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
		"prefix":  server.Srv.Config.Prefix,
		"entries": htmlEntries,
		"ws":      ws,
		"error":   err,
	})
}

func doGETFiles(c *gin.Context) {

	ws := store.Normalize(c.Param("ws"))

	files, err := server.Srv.Store.Entry.ListFiles(ws)
	if err != nil {
		dumpError(c, err)
		return
	}

	c.HTML(http.StatusOK, "files.tmpl", gin.H{
		"ws":     ws,
		"prefix": server.Srv.Config.Prefix,
		"files":  files,
	})
}
