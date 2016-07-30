package server

import (
	"log"
	"strconv"

	"github.com/adriamb/gopad/dict"
	"github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/store"
	"github.com/gin-gonic/gin"
)

// Server state definition
type Server struct {
	config.Config
	Engine *gin.Engine
	Store  *store.Store
	Dict   *dict.Dict
}

// Srv is the global server state
var Srv *Server

// Initialize the server
func Initialize(config config.Config) {

	store.InitCache(config.CacheDir, config.TmpDir)

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

	store := store.NewStore(config.DataDir)

	server := Server{
		Engine: g,
		Config: config,
		Store:  store,
		Dict:   dict.New(store),
	}

	Srv = &server
}

// Start the server
func Start() {
	err := Srv.Engine.Run(":" + strconv.Itoa(Srv.Config.Port))
	if err != nil {
		log.Fatal(err)
	}
}
