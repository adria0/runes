//go:generate rice embed-go

package web

import (
	"net/http"

	"github.com/adriamb/runes/server/instance"
	"github.com/gin-gonic/gin"
)

func doGETWorkspaces(c *gin.Context) {

	workspaces, err := instance.Srv.Store.Entry.ListWorkspaces()
	if err != nil {
		dumpError(c, err)
		return
	}

	c.HTML(http.StatusOK, "workspaces.tmpl", gin.H{
		"workspaces": workspaces,
	})
}

func doGETNewWorkspace(c *gin.Context) {

	ws := normalize(c.Param("ws"))

	err := instance.Srv.Store.Entry.CreateWorkspace(ws)
	if err != nil {
		dumpError(c, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/w/"+ws)
}

func doGETDeleteWorkspace(c *gin.Context) {

	ws := normalize(c.Param("ws"))

	err := instance.Srv.Store.Entry.DeleteWorkspace(ws)
	if err != nil {
		dumpError(c, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/w")
}
