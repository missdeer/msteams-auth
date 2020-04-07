package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	code string
)

type AccessTokenRequestBody struct {
	GrantType   string `json:"grant_type"`
	ClientID    string `json:"client_id"`
	Scope       string `json:"scope"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=46442420-1b26-4bd7-a997-183e1880bbd5&response_type=code&redirect_uri=https://msteam.ngrok.io/individual_user_consent/&response_mode=query&scope=https%3A%2F%2Fgraph.microsoft.com%2Fuser.read&state=12345")
	})

	r.GET("/individual_user_consent/", func(c *gin.Context) {
		responseCode := c.Query("code")
		log.Println("code", responseCode)
		if responseCode != "" {
			code = responseCode
			c.Redirect(http.StatusFound, "https://login.microsoftonline.com/afae2f63-1bcb-4d1f-b8c3-252a4cd3dd07/v2.0/adminconsent?client_id=46442420-1b26-4bd7-a997-183e1880bbd5&state=12345&redirect_uri=https://msteam.ngrok.io/individual_user_consent/&scope=https://graph.microsoft.com/user.read")
			return
		}
		admin_consent := c.Query("admin_consent")
		log.Println("admin_consent:", admin_consent)
		state := c.Query("state")
		log.Println("state:", state)
		scope := c.Query("scope")
		log.Println("scope:", scope)

		req, err := http.NewRequest("POST", "https://login.microsoftonline.com/afae2f63-1bcb-4d1f-b8c3-252a4cd3dd07/oauth2/v2.0/token",
			strings.NewReader(`client_secret=msHRpSOTQLP24lCk9afnSTejW%3DlV%3F8%3D%40&grant_type=authorization_code&client_id=46442420-1b26-4bd7-a997-183e1880bbd5&scope=https://graph.microsoft.com/user.read&redirect_uri=https://msteam.ngrok.io/individual_user_consent/&code=`+code))
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var response gin.H
		if err = json.Unmarshal(respBody, &response); err != nil {
			log.Fatal(err)
		}

		accessToken, ok := response["access_token"]
		if ok {
			log.Println("access_token", accessToken)
		}

		refreshToken, ok := response["refresh_token"]
		if ok {
			log.Println("refresh_token", refreshToken)
		}

		c.JSON(http.StatusOK, &response)
	})

	bind := os.Getenv("BINDADDR")
	if bind == "" {
		bind = "127.0.0.1:8765"
	}
	log.Fatal(r.Run(bind))
}
