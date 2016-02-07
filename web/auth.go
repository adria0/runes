package web

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/amassanet/gopad/server"
	"github.com/gin-gonic/gin"
    "math/rand"
    "time"
    "sync"
)

type userinfoStruct struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	Gender        string `json:"gender"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Link          string `json:"link"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

var (
    cookies = map[string]int64{}
    mutex = &sync.Mutex{}
)


func initAuthentication() {
    ticker := time.NewTicker(time.Hour)
    go func() {
        for _ = range ticker.C {
            now := time.Now().Unix()

            mutex.Lock()

            for k,expires := range cookies {
                if now > expires {
                    delete(cookies,k)
                    fmt.Println("Removed", k)
                }
            }
            mutex.Unlock()
        }
    }()
}

func getEmailFromOAuth2(accessToken string) (string, error) {
	client := &http.Client{}

	url := "https://www.googleapis.com/oauth2/v1/userinfo"

	res, err := client.Get(url + "?access_token=" + accessToken)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var userinfo userinfoStruct

	err = decoder.Decode(&userinfo)
	if err != nil {
		return "", err
	}

	return userinfo.Email, nil
}

func isAuthenticated(c *gin.Context) bool {

    validSession := func(token string) bool{
        now := time.Now().Unix()

        mutex.Lock()
        defer mutex.Unlock()

        if expires, exists := cookies[token] ; exists {
            if now > expires {
                delete(cookies,token)
                return false
            }
            return true
        }
        return false

    }

    valid := false

    for _, cookie := range c.Request.Cookies() {
        if cookie.Name == "token" {
            valid = validSession(cookie.Value)
            break
        }
    }

 	if !valid {
		c.Redirect(301, server.Srv.Config.Prefix+"/login")
	}

    return valid
}

func setAuthentication(c *gin.Context, oauthToken string) error{

    email, err := getEmailFromOAuth2(oauthToken)
    if err != nil {
        return err
    }
    if email != "adriamassanet@gmail.com" {
        return fmt.Errorf("Bad email %v",email)
    }

    token128 := fmt.Sprintf("%x%x%x%x",
        rand.Uint32(),rand.Uint32(),rand.Uint32(), rand.Uint32())
    expires := time.Now().Unix() + 7*24*3600

    mutex.Lock()
    cookies[token128]=expires
    mutex.Unlock()

    cookie := http.Cookie{Name: "token", Value: token128}
	http.SetCookie(c.Writer, &cookie)

    return nil
}


