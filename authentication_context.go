package authproxy

import "net/http"

type AuthenticationContext interface {
	IsAccessTokenValidAndUserAuthorized(accessToken string) bool
	GetUserName(accessToken string) string
	GetHTTPEndpointPrefix() string
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
