package instance

import (
	"github.com/adriamb/runes/server/config"
	"github.com/adriamb/runes/store"
	"github.com/gin-gonic/gin"
)

// Server state definition
type Server struct {
	config.Config
	Engine *gin.Engine
	Store  *store.Store
}

// Srv is the global server state
var Srv *Server
