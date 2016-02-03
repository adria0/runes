package server

import (
	"github.com/amassanet/gopad/store"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"html/template"
	"strconv"
    "os/user"
    "log"
)

type Config struct {
	Port   int
	Prefix string
}

type Server struct {
	Config
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

func templateReloader(c *gin.Context) {
	if tmpl, err := template.New("name").Funcs(funcName).ParseGlob("web/templates/*"); err == nil {
		Srv.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
}

func NewServer(config Config) {

	g := gin.New()
	g.Use(templateReloader, gin.Logger(), gin.Recovery())

	server := Server{
		Engine: g,
		Config: config,
	}

	if tmpl, err := template.New("name").Funcs(funcName).ParseGlob("web/templates/*"); err == nil {
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
