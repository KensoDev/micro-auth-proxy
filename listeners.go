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
	AuthContext *GithubAuthContext
}

func NewHttpListeners(config *Configuration) {
	authContext := NewGithubAuthContext(config)
	http.Handle("/callback", authContext)

	for _, upstream := range config.Upstreams {
		uri, _ := url.Parse(upstream.Location)

		proxy := httputil.NewSingleHostReverseProxy(uri)

		listener := &Listener{
			AuthContext: authContext,
			Prefix:      upstream.Prefix,
			Location:    upstream.Location,
			Proxy:       proxy,
			Hostname:    uri.Hostname(),
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

	l.ServeAuthenticatedRequest(w, req)
}

func (l *Listener) ServeAuthenticatedRequest(w http.ResponseWriter, req *http.Request) {
	director := l.Proxy.Director

	l.Proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = l.Hostname
	}

	l.Proxy.ServeHTTP(w, req)

	return
}
