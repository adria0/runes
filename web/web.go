//go:generate rice embed-go

package web

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/store"
	"github.com/adriamb/gopad/web/render"
	"github.com/gin-gonic/gin"
)

var markdownRender = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(render.Render(s))
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

func doGETFile(c *gin.Context) {
	file := server.Srv.Store.File.Fullpath(c.Param("id"))
	c.File(file)
}

func doGETCache(c *gin.Context) {
	file := store.GetCachePath(c.Param("id"))
	c.File(file)
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
