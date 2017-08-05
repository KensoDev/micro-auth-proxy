package authproxy

import (
	"encoding/json"
	"fmt"
)

type Configuration struct {
	AuthenticationContextName string     `json:"authContext"`
	Upstreams                 []Upstream `json:"upstreams"`
	Users                     []User     `json:"users"`
}

type User struct {
	Username string `json:"username"`
	Restrict string `json:"restrict"`
}

type Upstream struct {
	Prefix   string `json:"prefix"`
	Location string `json:"location"`
	Type     string `json:"type"`
}

func (c *Configuration) GetAuthenticationContext() AuthenticationContext {
	if c.AuthenticationContextName == "github" {
		return NewGithubAuthContext(c)
	}

	return nil
}

func NewConfiguration(data []byte) (*Configuration, error) {
	config := &Configuration{}
	err := json.Unmarshal(data, config)

	if err != nil {
		return nil, fmt.Errorf("Problem with parsing the confi json file: %s", err.Error())
	}

	if len(config.Users) == 0 {
		return nil, fmt.Errorf("You have no users configured")

	}

	return config, nil
}

func (c *Configuration) GetRestrictionsForUsername(username string) string {
	for _, user := range c.Users {
		if user.Username == username {
			return user.Restrict
			break
		}
	}

	return "NotAllowed"
}

func (c *Configuration) ShouldRestrictUser(username string, method string) bool {
	allowedMethod := c.GetRestrictionsForUsername(username)

	// Allowed all methods
	if allowedMethod == "" {
		return true
	}

	return allowedMethod == method
}
