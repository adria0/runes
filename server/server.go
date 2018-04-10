package server

import (
	"log"

	"github.com/adriamb/runes/server/config"
	"github.com/adriamb/runes/server/instance"
	"github.com/adriamb/runes/store"
	"github.com/adriamb/runes/web"
	"github.com/gin-gonic/gin"
)

// Initialize the server
func startServer(config config.Config) {

	store.InitCache(config.CacheDir, config.TmpDir)

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

	store := store.NewStore(config.DataDir)
	err := store.Entry.Open()
	if err != nil {
		log.Fatal(err)
	}

	instance.Srv = &instance.Server{
		Engine: g,
		Config: config,
		Store:  store,
	}

	web.Initialize()
	if instance.Srv.Config.WebServer.CertFile != "" {
		err = instance.Srv.Engine.RunTLS(
			instance.Srv.Config.WebServer.Bind,
			instance.Srv.Config.WebServer.CertFile,
			instance.Srv.Config.WebServer.KeyFile,
		)
	} else {
		err = instance.Srv.Engine.Run(instance.Srv.Config.WebServer.Bind)
	}
	if err != nil {
		log.Fatal(err)
	}

}
