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

	server.Srv.Engine.GET("/", doGETLogin)
	server.Srv.Engine.GET("/login", doGETLogin)

	authorized := server.Srv.Engine.Group("/")
	authorized.Use(checkAuthorization())

	authorized.GET("/w/:ws", doGETEntries)
	authorized.GET("/w/:ws/f", doGETFiles)
	authorized.POST("/w/:ws/search", doPOSTSearch)

	authorized.GET("/w/:ws/e/:id", doGETEntries)
	authorized.GET("/w/:ws/e/:id/edit", doGETEntry)
	authorized.POST("/w/:ws/e/:id", doPOSTEntry)
	authorized.POST("/w/:ws/e/:id/f", doPOSTUpload)
	authorized.GET("/w/:ws/e/:id/f/:name", doGETFile)

	authorized.POST("/logingoauth2", doPOSTGoogleOauth2Login)
	authorized.POST("/render", doPOSTRender)
	authorized.GET("/cache/:hash", doGETCache)
}

func doGETFile(c *gin.Context) {

	ID := store.Normalize(c.Param("id"))
	ws := store.Normalize(c.Param("ws"))
	name := store.Normalize(c.Param("name"))

	file := server.Srv.Store.Entry.FilePath(ws, ID, name)
	c.File(file)
}

func doGETCache(c *gin.Context) {

	hash := store.Normalize(c.Param("hash"))

	file := store.GetCachePath(hash)
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
