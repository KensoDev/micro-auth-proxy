package authproxy

import (
	"log"
	"net/http"
	"os"
)

type AuthenticationContext interface {
	IsAccessTokenValidAndUserAuthorized(accessToken string) bool
	GetUserName(accessToken string) string
	GetHTTPEndpointPrefix() string
	GetCookieName() string
	GetLoginPage() ([]byte, error)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// GetenvOrDie is a safety wrapper around os.Getenv.
// This fatals if the key is unset or empty.
func GetenvOrDie(k string) string {
	o := os.Getenv(k)

	if o == "" {
		log.Fatalf("%s environment variable missing and is required.", k)
	}

	return o
}
