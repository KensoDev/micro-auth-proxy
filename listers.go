package authproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Listener struct {
	Prefix   string
	Location string
}

func NewHttpListeners(config *Configuration) {
	for _, upstream := range config.Upstreams {
		listener := &Listener{
			Prefix:   upstream.Prefix,
			Location: upstream.Location,
		}

		uri, _ := url.Parse(listener.Location)
		proxy := httputil.NewSingleHostReverseProxy(uri)

		http.Handle(listener.Prefix, proxy)
	}
}
