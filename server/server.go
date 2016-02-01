package server

import (
	"github.com/amassanet/gopad/store"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"html/template"
	"strconv"
)

type ServerConfiguration struct {
	Port   int
	Prefix string
}

type Server struct {
	Config ServerConfiguration
	Engine *gin.Engine
	Store  *store.Store
}

var Srv *Server

var funcName = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(blackfriday.MarkdownCommon([]byte(s)))
		return template.HTML(proc)
	},
}

func TemplateReloader(c *gin.Context) {
	if tmpl, err := template.New("name").Funcs(funcName).ParseGlob("web/templates/*"); err == nil {
		Srv.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
}

func NewServer(config ServerConfiguration) {

	g := gin.New()
	g.Use(TemplateReloader, gin.Logger(), gin.Recovery())

	server := Server{
		Engine: g,
		Config: config,
	}

	if tmpl, err := template.New("name").Funcs(funcName).ParseGlob("web/templates/*"); err == nil {
		server.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

	//server.Engine.LoadHTMLGlob("web/templates/*")
	server.Store = store.NewStore("/tmp/e")
	Srv = &server
}

func StartServer() {
    Srv.Engine.Run(":" + strconv.Itoa(Srv.Config.Port))
}
