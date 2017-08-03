package authproxy

import (
	"encoding/json"
	"fmt"
)

type Configuration struct {
	Upstreams []Upstream `json:"upstreams"`
	Users     []string   `json:"users"`
}

type Upstream struct {
	Prefix   string `json:"prefix"`
	Location string `json:"location"`
	Type     string `json:"type"`
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
