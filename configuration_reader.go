package authproxy

import (
	"fmt"
	"io/ioutil"
)

type ConfigurationReader struct {
	FileLocation string
}

func NewConfigurationReader(location string) *ConfigurationReader {
	return &ConfigurationReader{
		FileLocation: location,
	}
}

func (r *ConfigurationReader) ReadConfigurationFile() ([]byte, error) {
	data, err := ioutil.ReadFile(r.FileLocation)

	if err != nil {
		return nil, fmt.Errorf("Error reading the file: %s", err.Error())
	}

	return data, nil
}
