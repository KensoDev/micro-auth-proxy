package authproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Listener struct {
	Prefix      string
	Location    string
	Proxy       *httputil.ReverseProxy
	Hostname    string
	AuthContext AuthenticationContext
	Config      *Configuration
}

func NewHttpListeners(config *Configuration) {
	authContext := config.GetAuthenticationContext()
	http.Handle(authContext.GetHTTPEndpointPrefix(), authContext)

	for _, upstream := range config.Upstreams {
		uri, _ := url.Parse(upstream.Location)

		proxy := httputil.NewSingleHostReverseProxy(uri)

		listener := &Listener{
			AuthContext: authContext,
			Prefix:      upstream.Prefix,
			Location:    upstream.Location,
			Proxy:       proxy,
			Hostname:    uri.Hostname(),
			Config:      config,
		}

		http.Handle(listener.Prefix, listener)
	}
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("github_token")

	if err != nil {
		auth, _ := publicAuthHtmlBytes()
		w.Write(auth)
		return
	}

	token := cookie.Value

	if !l.AuthContext.IsAccessTokenValidAndUserAuthorized(token) {
		denied, _ := publicDeniedHtmlBytes()
		w.Write(denied)
		return
	}

	l.ServeAuthenticatedRequest(w, req, token)
}

func (l *Listener) ServeAuthenticatedRequest(w http.ResponseWriter, req *http.Request, accessToken string) {
	username := l.AuthContext.GetUserName(accessToken)
	allowed := l.Config.ShouldRestrictUser(username, req.Method)

	if !allowed {
		http.Error(w, "Your user is not allowed to perform this action", http.StatusUnauthorized)
		return
	}

	director := l.Proxy.Director
	l.Proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = l.Hostname
	}

	l.Proxy.ServeHTTP(w, req)

	return
}
