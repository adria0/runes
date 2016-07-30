//go:generate rice embed-go

package web

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/gopad/model"
	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/store"
	"github.com/adriamb/gopad/web/auth"
	"github.com/adriamb/gopad/web/render"
	"github.com/gin-gonic/gin"
)

var aa = auth.New()

var markdownRender = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(render.Render(s, server.Srv.Dict))
		return template.HTML(proc)
	},
}

func generateTemplate() *template.Template {
	templateList := []string{
		"500.tmpl", "entry.tmpl", "logingoauth2.tmpl",
		"search.tmpl", "entries.tmpl", "files.tmpl", "menu.tmpl",
	}

	tbox, tboxerr := rice.FindBox("templates")
	tmpl := template.New("name").Funcs(markdownRender)
	for _, name := range templateList {
		var content string
		if tboxerr == nil {
			content = tbox.MustString(name)
		} else {
			bytes, err := ioutil.ReadFile("web/templates/" + name)
			if err != nil {
				panic(err)
			}
			content = string(bytes)
		}
		_, err := tmpl.New(name).Parse(content)
		if err != nil {
			panic(err)
		}
	}
	return tmpl
}

// InitWeb Initializes the web
func InitWeb() {

	server.Srv.Engine.SetHTMLTemplate(generateTemplate())

	tbox, err := rice.FindBox("httpstatic")

	if err == nil {
		server.Srv.Engine.StaticFS("/static", tbox.HTTPBox())
	} else {
		server.Srv.Engine.StaticFS("/static", http.Dir("web/httpstatic"))
	}

	server.Srv.Engine.GET("/login", doGETLogin)

	authorized := server.Srv.Engine.Group("/")
	authorized.Use(checkAuthorization())

	authorized.GET("/", doGETEntries)
	authorized.POST("/logingoauth2", doPOSTGoogleOauth2Login)
	authorized.GET("/entries", doGETEntries)
	authorized.GET("/entries/:id", doGETEntries)
	authorized.GET("/entries/:id/edit", doGETEntry)
	authorized.POST("/entries/:id/edit", doPOSTEntry)
	authorized.POST("/markdown", doPOSTMarkdown)
	authorized.POST("/entries/:id/edit/files", doPOSTUpload)
	authorized.GET("/files/:id", doGETFile)
	authorized.GET("/files", doGETFiles)
	authorized.GET("/cache/:id", doGETCache)
	authorized.POST("/search", doPOSTSearch)
}

func checkAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !aa.IsAuthorized(c) {
			c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/login")
			return
		}
	}
}

type tmplEntry struct {
	*model.Entry
}

type dtoMarkdownRender struct {
	Markdown string `json:"markdown" binding:"required"`
}

func doPOSTSearch(c *gin.Context) {

	query := c.Request.FormValue("query")
	results, err := server.Srv.Store.Entry.Search(query)
	if err != nil {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"prefix": server.Srv.Config.Prefix,
			"error":  err,
		})
		return
	}
	if len(results) == 0 {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"prefix": server.Srv.Config.Prefix,
			"info":   "No results",
		})
		return
	}
	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"prefix":  server.Srv.Config.Prefix,
		"results": results,
	})
}

func doPOSTUpload(c *gin.Context) {

	id := c.Param("id")

	if existsStaticMd(id) {
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename, err := server.Srv.Store.File.Write(fileHeader.Filename, id, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := server.Srv.Config.Prefix + "/files/" + filename
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

func doGETCache(c *gin.Context) {
	file := store.GetCachePath(c.Param("id"))
	c.File(file)
}

func doPOSTMarkdown(c *gin.Context) {
	var json dtoMarkdownRender
	if c.BindJSON(&json) == nil {
		html := string(render.Render(json.Markdown, server.Srv.Dict))
		c.JSON(http.StatusOK, gin.H{"html": html})
	}
}

func doGETEntries(c *gin.Context) {

	var entries []*model.Entry
	var err error

	id := c.Param("id")

	if id != "" {

		var entry *model.Entry
		if existsStaticMd(id) {
			entry, err = getStaticMdEntry(id)
		} else {
			entry, err = server.Srv.Store.Entry.Get(id)
		}

		if err == nil {
			entries = append(entries, entry)
		}

	} else {

		entries, err = server.Srv.Store.Entry.List()

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
		"error":   err,
	})
}

func doGETFiles(c *gin.Context) {

	files, err := server.Srv.Store.File.List()
	if err != nil {
		dumpError(c, err)
		return
	}

	c.HTML(http.StatusOK, "files.tmpl", gin.H{
		"prefix": server.Srv.Config.Prefix,
		"files":  files,
	})
}

func doGETLogin(c *gin.Context) {

	if server.Srv.Config.Auth.Type == config.AuthNone {
		aa.Authorize(c)
		c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/entries")
		return
	}

	if server.Srv.Config.Auth.Type == config.AuthGoogle {

		var err error

		c.HTML(http.StatusOK, "logingoauth2.tmpl", gin.H{
			"prefix":         server.Srv.Config.Prefix,
			"googleclientid": server.Srv.Config.Auth.GoogleClientID,
			"error":          err,
		})
		return

	}

	log.Fatalf("Server authentication type '%v' is not known.",
		server.Srv.Config.Auth.Type)

}

func doPOSTGoogleOauth2Login(c *gin.Context) {

	oauthtoken := c.DefaultPostForm("oauthtoken", "undefined")

	err := aa.AuthorizeGoogleOauth2(c, oauthtoken)
	if err != nil {
		c.HTML(http.StatusOK, "logingoauth2.tmpl", gin.H{
			"prefix": server.Srv.Config.Prefix,
			"error":  err,
		})
	} else {
		c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/entries")
	}
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
	var err error
	var editable = true

	if id != "new" {

		if existsStaticMd(id) {
			entry, err = getStaticMdEntry(id)
			editable = false
		} else {
			entry, err = server.Srv.Store.Entry.Get(id)
		}

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
		"prefix":   server.Srv.Config.Prefix,
		"entry":    entry,
		"editable": editable,
	})

}

func doPOSTEntry(c *gin.Context) {

	id := c.Param("id")

	if existsStaticMd(id) {
		return
	}

	if buttonPressed(c, "btnsave") {

		entry := model.Entry{
			ID:       id,
			Title:    c.DefaultPostForm("Title", "undefined"),
			Markdown: c.DefaultPostForm("Markdown", "undefined"),
		}

		err := server.Srv.Store.Entry.Store(&entry)
		if err != nil {
			err = server.Srv.Dict.Rebuild()
		}
		if err != nil {
			dumpError(c, err)
			return
		}
		c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/entries")
		return
	}

	if buttonPressed(c, "btnback") {
		c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/entries")
		return
	}
}
