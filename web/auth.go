//go:generate rice embed-go

package web

import (
	"log"
	"net/http"

	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/web/auth"
	"github.com/gin-gonic/gin"
)

var aa = auth.New()

func checkAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !aa.IsAuthorized(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
	}
}

func doGETLogin(c *gin.Context) {

	if server.Srv.Config.Auth.Type == config.AuthNone {
		aa.Authorize(c)
		c.Redirect(http.StatusSeeOther, "/w/default")
		return
	}

	if server.Srv.Config.Auth.Type == config.AuthGoogle {

		var err error

		c.HTML(http.StatusOK, "logingoauth2.tmpl", gin.H{
			"googleclientid": server.Srv.Config.Auth.GoogleClientID,
			"error":          err,
		})
		return

	}

	log.Fatalf("Server authentication type '%v' is not known.",
		server.Srv.Config.Auth.Type)

}

func doPOSTGoogleOauth2Login(c *gin.Context) {

	oauthtoken := c.DefaultPostForm("oauthtoken", "undefined")

	err := aa.AuthorizeGoogleOauth2(c, oauthtoken)
	if err != nil {
		c.HTML(http.StatusOK, "logingoauth2.tmpl", gin.H{
			"error":  err,
		})
	} else {
		c.Redirect(http.StatusSeeOther, "/w/default")
	}
}
