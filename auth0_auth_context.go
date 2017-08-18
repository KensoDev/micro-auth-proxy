package authproxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Auth0AuthContext struct {
	Config            *Configuration
	AuthDomain        string
	ClientID          string
	ClientSecret      string
	CallbackURL       string
	ValidAccessTokens map[string]string
	HTMLFile          []byte
}

func NewAuth0AuthContext(config *Configuration) *Auth0AuthContext {
	return &Auth0AuthContext{
		Config:            config,
		ClientID:          GetenvOrDie("AUTH0_CLIENT_ID"),
		ClientSecret:      GetenvOrDie("AUTH0_CLIENT_SECRET"),
		AuthDomain:        GetenvOrDie("AUTH0_DOMAIN"),
		CallbackURL:       GetenvOrDie("AUTH0_CALLBACK_URL"),
		ValidAccessTokens: map[string]string{},
	}
}

func (c *Auth0AuthContext) IsAccessTokenValidAndUserAuthorized(accessToken string) bool {
	_, ok := c.ValidAccessTokens[accessToken]

	if ok {
		return true
	}

	return false
}

func (c *Auth0AuthContext) GetUserName(accessToken string) string {
	username, _ := c.ValidAccessTokens[accessToken]
	return username
}

func (c *Auth0AuthContext) GetHTTPEndpointPrefix() string {
	return "/callback"
}

func (c *Auth0AuthContext) GetCookieName() string {
	return "auth0_token"
}

func (c *Auth0AuthContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conf := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.CallbackURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/authorize", c.AuthDomain),
			TokenURL: fmt.Sprintf("https://%s/oauth/token", c.AuthDomain),
		},
	}

	code := req.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken := token.AccessToken

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    c.GetCookieName(),
		Value:   accessToken,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get(fmt.Sprintf("https://%s/userinfo", c.AuthDomain))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var username string
	username = profile["email"].(string)

	usernames := MapUserNames(c.Config.Users, func(user interface{}) string {
		return user.(User).Username
	})

	fmt.Println(username)
	fmt.Println(usernames)

	if inArray(username, usernames) {
		c.ValidAccessTokens[accessToken] = username
	}

	http.Redirect(w, req, "/", 302)
}

func (c *Auth0AuthContext) GetLoginPage() ([]byte, error) {
	return c.HTMLFile, nil
}

func (c *Auth0AuthContext) RenderHTMLFile() error {
	tplBytes, err := publicAuth0HtmlTplBytes()
	if err != nil {
		return err
	}

	tpl := string(tplBytes)

	f, err := RenderTemplate(tpl, c)
	if err != nil {
		return err
	}

	c.HTMLFile = f
	return nil
}
