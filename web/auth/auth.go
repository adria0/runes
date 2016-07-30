package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Auth contains the auth sessions
type Auth struct {
	sync.Mutex
	sessions map[string]int64
}

// New auth
func New() *Auth {
	return &Auth{
		sessions: make(map[string]int64),
	}
}

// IsAuthorized returns if the http context contains a valid cookie
func (a *Auth) IsAuthorized(c *gin.Context) bool {

	var sessionCookie *http.Cookie

	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == "token" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		return false
	}

	a.Lock()
	defer a.Unlock()

	now := time.Now().Unix()

	if expires, exists := a.sessions[sessionCookie.Value]; exists {
		if now > expires {
			delete(a.sessions, sessionCookie.Value)
			return false
		}
		return true
	}

	return false

}

// Authorize adds a new authorization cookie
func (a *Auth) Authorize(c *gin.Context) {

	token128 := fmt.Sprintf("%x%x%x%x",
		rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32())

	expires := time.Now().Unix() + 7*24*3600

	a.Lock()
	a.sessions[token128] = expires
	a.Unlock()

	cookie := http.Cookie{Name: "token", Value: token128}
	http.SetCookie(c.Writer, &cookie)
}
