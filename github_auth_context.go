package authproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type GithubAuthContext struct {
	ClientID     string
	ClientSecret string
	Config       *Configuration
}

type GithubAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

func NewGithubAuthContext(config *Configuration) *GithubAuthContext {
	return &GithubAuthContext{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Config:       config,
	}
}

type GithubAuthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func (c *GithubAuthContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	client := &http.Client{}

	githubRequest := &GithubAuthRequest{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Code:         code,
	}

	jsonRequestBody, _ := json.Marshal(githubRequest)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonRequestBody))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	responseBody, _ := ioutil.ReadAll(resp.Body)

	githubResponse := &GithubAuthResponse{}

	err = json.Unmarshal(responseBody, githubResponse)

	if err != nil {
		log.Fatal(err)
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    "github_token",
		Value:   githubResponse.AccessToken,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, req, "/", 302)
}
