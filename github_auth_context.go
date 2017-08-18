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
	ClientID          string
	ClientSecret      string
	Config            *Configuration
	CookirName        string
	ValidAccessTokens map[string]string
}

type GithubAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

func NewGithubAuthContext(config *Configuration) *GithubAuthContext {
	return &GithubAuthContext{
		ClientID:          os.Getenv("CLIENT_ID"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		Config:            config,
		ValidAccessTokens: map[string]string{},
	}
}

func (c *GithubAuthContext) GetLoginPage() ([]byte, error) {
	return publicAuthHtmlBytes()
}

func (c *GithubAuthContext) GetCookieName() string {
	return "github_token"
}

type mapf func(interface{}) string

func MapUserNames(in []User, f mapf) []string {
	newArray := []string{}

	for _, v := range in {
		newArray = append(newArray, f(v))
	}

	return newArray
}

type GithubUser struct {
	ID       int    `json:"id"`
	UserName string `json:"login"`
}

type GithubAuthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func (c *GithubAuthContext) IsAccessTokenValidAndUserAuthorized(accessToken string) bool {
	_, ok := c.ValidAccessTokens[accessToken]

	if ok {
		return true
	}

	responseBytes, err := c.GetUserDetailsFromGithub(accessToken)

	if err != nil {
		return false
	}

	githubUser, err := c.ParseUserResponse(responseBytes)

	if err != nil {
		return false
	}

	usernames := MapUserNames(c.Config.Users, func(user interface{}) string {
		return user.(User).Username
	})

	userExists := inArray(githubUser.UserName, usernames)

	if userExists {
		c.ValidAccessTokens[accessToken] = githubUser.UserName
	}

	return userExists
}

func (c *GithubAuthContext) GetUserName(accessToken string) string {
	username, _ := c.ValidAccessTokens[accessToken]
	return username
}

func (c *GithubAuthContext) ParseUserResponse(response []byte) (*GithubUser, error) {
	githubUser := &GithubUser{}
	err := json.Unmarshal(response, githubUser)

	return githubUser, err
}

func (c *GithubAuthContext) GetUserDetailsFromGithub(accessToken string) ([]byte, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("https://api.github.com/user?access_token=%s", accessToken)
	req, err := http.NewRequest("GET", uri, nil)
	resp, err := client.Do(req)

	responseBody, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	return responseBody, err
}

func (c *GithubAuthContext) GetHTTPEndpointPrefix() string {
	return "/callback"
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
		Name:    c.GetCookieName(),
		Value:   githubResponse.AccessToken,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, req, "/", 302)
}
