package authproxy

import (
	"fmt"
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
	for _, upstream := range config.Upstreams {

		uri, _ := url.Parse(upstream.Location)

		proxy := httputil.NewSingleHostReverseProxy(uri)

		listener := &Listener{
			Prefix:   upstream.Prefix,
			Location: upstream.Location,
			Proxy:    proxy,
			Hostname: uri.Hostname(),
		}

		http.Handle(listener.Prefix, listener)
	}
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("github_token")

	if err != nil {
		bytes, _ := publicAuthHtmlBytes()
		w.Write(bytes)

		return
	}

	token := cookie.Value

	if token == "" {
		bytes, _ := publicAuthHtmlBytes()
		w.Write(bytes)

		return
	}

	l.ServeAuthenticatedRequest(w, req)
}

func (l *Listener) ServeAuthenticatedRequest(w http.ResponseWriter, req *http.Request) {
	director := l.Proxy.Director

	l.Proxy.Director = func(req *http.Request) {
		director(req)
		fmt.Println("Hostname")
		fmt.Println(l.Hostname)
		req.Host = l.Hostname
	}

	l.Proxy.ServeHTTP(w, req)

	return
}
