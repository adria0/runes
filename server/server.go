package server

import (
    "github.com/GeertJohan/go.rice"
    "github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/dict"
	"github.com/adriamb/gopad/store"
	"github.com/adriamb/gopad/web/render"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os/user"
	"strconv"
    "fmt"
)

type Server struct {
	config.Config
	Engine *gin.Engine
	Store  *store.Store
	Dict   *dict.Dict
}

var Srv *Server

var markdownRender = template.FuncMap{
	"markdown": func(s string) template.HTML {
		proc := string(render.Render(s, Srv.Dict))
		return template.HTML(proc)
	},
}

func templateReloader(c *gin.Context) {
	if tmpl, err := template.New("name").Funcs(markdownRender).ParseGlob("web/templates/*"); err == nil {
        fmt.Printf("%#v",tmpl)
        Srv.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
}

func NewServer(config config.Config) {

	store.InitCache()

    tbox, _ := rice.FindBox("../web/templates")
    if tbox == nil {

    }

	g := gin.New()
	g.Use(templateReloader, gin.Logger(), gin.Recovery())

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(usr.HomeDir + "/.gopad")

	server := Server{
		Engine: g,
		Config: config,
		Store:  store,
		Dict:   dict.New(store),
	}
/*
	if tmpl, err := template.New("name").Funcs(markdownRender).ParseGlob("web/templates/*"); err == nil {
		server.Engine.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
*/
	Srv = &server
}

func StartServer() {
	Srv.Engine.Run(":" + strconv.Itoa(Srv.Config.Port))
}
