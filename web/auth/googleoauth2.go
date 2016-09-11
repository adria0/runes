package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adriamb/runes/server/instance"
	"github.com/gin-gonic/gin"
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

func getEmailFromGoogle(accessToken string) (string, error) {
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

// AuthorizeGoogleOauth2 adds new autorization cookie from a oauth token
func (a *Auth) AuthorizeGoogleOauth2(c *gin.Context, oauthToken string) error {

	email, err := getEmailFromGoogle(oauthToken)
	if err != nil {
		return err
	}

	var found *string
	for _, allowed := range instance.Srv.Config.Auth.AllowedEmails {
		if email == allowed {
			found = &allowed
			break
		}
	}
	if found == nil {
		return fmt.Errorf("Bad email %v", email)
	}

	a.Authorize(c)

	return nil
}
