//go:generate rice embed-go

package web

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/adriamb/runes/server/instance"
	"github.com/adriamb/runes/store"
	"github.com/adriamb/runes/web/render"
	"github.com/gin-gonic/gin"
)

var markdownRender = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(render.Render(s))
		return template.HTML(proc)
	},
	"raw": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func generateTemplate() *template.Template {
	templateList := []string{
		"500.tmpl", "builtin.tmpl", "entry.tmpl", "logingoauth2.tmpl",
		"search.tmpl", "entries.tmpl", "files.tmpl",
		"headbegin.tmpl", "headend.tmpl", "tail.tmpl", "workspaces.tmpl",
		"workspacemenu.tmpl",
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

// Initialize Initializes the web
func Initialize() {

	instance.Srv.Engine.SetHTMLTemplate(generateTemplate())

	tbox, err := rice.FindBox("httpstatic")

	if err == nil {
		instance.Srv.Engine.StaticFS("/static", tbox.HTTPBox())
	} else {
		instance.Srv.Engine.StaticFS("/static", http.Dir("web/httpstatic"))
	}

	instance.Srv.Engine.GET("/", doGETLogin)
	instance.Srv.Engine.GET("/login", doGETLogin)

	authorized := instance.Srv.Engine.Group("/")
	authorized.Use(checkAuthorization())

	authorized.GET("/w", doGETWorkspaces)

	authorized.GET("/w/:ws", doGETEntries)
	authorized.GET("/w/:ws/f", doGETFiles)
	authorized.GET("/w/:ws/new", doGETNewWorkspace)
	authorized.GET("/w/:ws/delete", doGETDeleteWorkspace)
	authorized.POST("/w/:ws/search", doPOSTSearch)

	authorized.GET("/w/:ws/e/:id", doGETEntries)
	authorized.POST("/w/:ws/e/:id/quickadd", doPOSTQuickadd)
	authorized.GET("/w/:ws/e/:id/edit", doGETEntryEdit)
	authorized.POST("/w/:ws/e/:id/delete", doPOSTEntryDelete)
	authorized.POST("/w/:ws/e/:id", doPOSTEntry)
	authorized.POST("/w/:ws/e/:id/f", doPOSTUpload)
	authorized.GET("/w/:ws/e/:id/f/:name", doGETFile)

	authorized.GET("/builtin/:id", doGETBuiltin)
	authorized.POST("/logingoauth2", doPOSTGoogleOauth2Login)
	authorized.POST("/render", doPOSTRender)
	authorized.GET("/cache/:hash", doGETCache)
}

func normalize(s string) string {

	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		rv := int(rs[i])
		if (rv >= int('A') && rv <= int('Z')) ||
			(rv >= int('a') && rv <= int('z')) ||
			(rv >= int('0') && rv <= int('9')) ||
			(rv == int('.')) {

		} else {
			rs[i] = '_'

		}

	}
	return string(rs)

}

func doGETFile(c *gin.Context) {

	ID := normalize(c.Param("id"))
	ws := normalize(c.Param("ws"))
	name := normalize(c.Param("name"))

	file := instance.Srv.Store.Entry.FilePath(ws, ID, name)
	c.File(file)
}

func doGETCache(c *gin.Context) {

	hash := normalize(c.Param("hash"))

	file := store.GetCachePath(hash)
	c.File(file)
}

func buttonPressed(c *gin.Context, name string) bool {
	return c.DefaultPostForm(name, "undefined") != "undefined"
}

func dumpError(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "500.tmpl", gin.H{
		"message": err.Error(),
	})
}
