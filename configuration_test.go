package authproxy

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestConfiguration(t *testing.T) { TestingT(t) }

type ConfigurationSuite struct{}

var _ = Suite(&ConfigurationSuite{})

func (s *ConfigurationSuite) TestJsonReadFile(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	c.Assert(len(config.Upstreams), Equals, 2)
}
