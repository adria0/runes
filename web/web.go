package web

import (
	"github.com/amassanet/gopad/model"
	"github.com/amassanet/gopad/server"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"net/http"
	"os"
	"strings"
)

// InitWeb Initializes the web
func InitWeb() {

	server.Srv.Engine.StaticFS("/static", http.Dir("web/static"))

	server.Srv.Engine.GET("/", doGETRoot)
	server.Srv.Engine.GET("/entry/:id", doGETEntry)
	server.Srv.Engine.POST("/entry/:id", doPOSTEntry)
	server.Srv.Engine.POST("/markdown", doPOSTMarkdown)
    server.Srv.Engine.POST("/entry/:id/file", doPOSTUpload)
	server.Srv.Engine.GET("/file/:id", doGETFile)
}

type tmplEntry struct {
	*model.Entry
}

type dtoMarkdownRender struct {
	Markdown string `json:"markdown" binding:"required"`
}

func doPOSTUpload(c *gin.Context) {

    id := c.Param("id")

    file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
        c.JSON(http.StatusBadRequest , gin.H{"error":err.Error()})
		return
	}

	filename, err := server.Srv.Store.File.Write(fileHeader.Filename, id, file)
	if err != nil {
        c.JSON(http.StatusBadRequest , gin.H{"error":err.Error()})
		return
	}

	path := server.Srv.Config.Prefix + "/file/" + filename
	split := strings.Split(filename, ".")
	ext := split[len(split)-1]
	ico := ""

	if ext == "png" || ext == "jpg" || ext == "jpeg" || ext == "gif" {

	} else if _, err := os.Stat("web/static/ico/" + ext + ".png"); err == nil {
		ico = "/static/ico/" + ext + ".png"
	} else {
		ico = "/static/ico/default.png"
	}

	c.JSON(http.StatusOK, gin.H{"name": fileHeader.Filename, "path": path, "ico": ico})
}

func doGETFile(c *gin.Context) {
	file := server.Srv.Store.File.Fullpath(c.Param("id"))
	c.File(file)
}

func doPOSTMarkdown(c *gin.Context) {
	var json dtoMarkdownRender
	if c.BindJSON(&json) == nil {
		html := string(blackfriday.MarkdownCommon([]byte(json.Markdown)))
		c.JSON(http.StatusOK, gin.H{"html": html})
	}
}

func doGETRoot(c *gin.Context) {

	entries, err := server.Srv.Store.Entry.List()
	if err != nil {
		dumpError(c, err)
		return
	}

	htmlEntries := []tmplEntry{}
	for _, entry := range entries {
		htmlEntries = append(htmlEntries, tmplEntry{entry})
	}

	err = nil
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":   "Main website",
		"prefix":  server.Srv.Config.Prefix,
		"entries": htmlEntries,
		"error":   err,
	})
}

func buttonPressed(c *gin.Context, name string) bool {
	return c.DefaultPostForm(name, "undefined") != "undefined"
}

func dumpError(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "500.tmpl", gin.H{
		"prefix":  server.Srv.Config.Prefix,
		"message": err.Error(),
	})
}

func doGETEntry(c *gin.Context) {

	id := c.Param("id")

	var entry *model.Entry

	if id != "new" {
		var err error
		entry, err = server.Srv.Store.Entry.Get(id)
		if err != nil {
		    dumpError(c, err)
            return
		}
    } else {
        entry = &model.Entry{
            ID: server.Srv.Store.Entry.NewID(),
        }
	}

	c.HTML(http.StatusOK, "entry.tmpl", gin.H{
		"prefix": server.Srv.Config.Prefix,
		"entry":  entry,
	})

}

func doPOSTEntry(c *gin.Context) {

	if buttonPressed(c, "btnsave") {

		entry := model.Entry{
			ID:       c.Param("id"),
			Title:    c.DefaultPostForm("Title", "undefined"),
			Markdown: c.DefaultPostForm("Markdown", "undefined"),
		}

		err := server.Srv.Store.Entry.Store(&entry)
		if err != nil {
			dumpError(c, err)
			return
		}
		c.Redirect(301, server.Srv.Config.Prefix+"/")
		return
	}

	if buttonPressed(c, "btnback") {
		c.Redirect(301, server.Srv.Config.Prefix+"/")
		return
	}
}

