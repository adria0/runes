package server

import (
	"github.com/amassanet/gopad/store"
	"github.com/amassanet/gopad/web/render"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
    "os/user"
    "log"
)

type Config struct {
	Port   int
	Prefix string
    Auth struct {
        GoogleClientID string
        AllowedEmails []string
    }
}

type Server struct {
	Config
	Engine *gin.Engine
	Store  *store.Store
}

var Srv *Server

var markdownRender = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(render.Render(s))
		return template.HTML(proc)
	},
}

func templateReloader(c *gin.Context) {
	if tmpl, err := template.New("name").Funcs(markdownRender).ParseGlob("web/templates/*"); err == nil {
		Srv.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
}

func NewServer(config Config) {

    store.InitCache()

	g := gin.New()
	g.Use(templateReloader, gin.Logger(), gin.Recovery())

	server := Server{
		Engine: g,
		Config: config,
	}

	if tmpl, err := template.New("name").Funcs(markdownRender ).ParseGlob("web/templates/*"); err == nil {
		server.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

    usr,err := user.Current()
    if err!= nil {
        log.Fatal(err)
    }

	server.Store = store.NewStore(usr.HomeDir+"/.gopad")
	Srv = &server
}

func StartServer() {
    Srv.Engine.Run(":" + strconv.Itoa(Srv.Config.Port))
}

