package server

import (
	"github.com/adriamb/gopad/dict"
	"github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/store"
	"github.com/gin-gonic/gin"
	"log"
	"os/user"
	"strconv"
	//"fmt"
)

type Server struct {
	config.Config
	Engine *gin.Engine
	Store  *store.Store
	Dict   *dict.Dict
}

var Srv *Server

func NewServer(config config.Config) {

	store.InitCache()

	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

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

	Srv = &server
}

func StartServer() {
	Srv.Engine.Run(":" + strconv.Itoa(Srv.Config.Port))
}


