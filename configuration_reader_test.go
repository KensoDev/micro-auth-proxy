package authproxy

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestConfigurationReader(t *testing.T) { TestingT(t) }

type ConfigurationReaderSuite struct{}

var _ = Suite(&ConfigurationReaderSuite{})

func (s *ConfigurationReaderSuite) TestJsonReadFile(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, err := reader.ReadConfigurationFile()
	c.Assert(err, IsNil)
	c.Assert(len(data), Not(Equals), 0)
}
