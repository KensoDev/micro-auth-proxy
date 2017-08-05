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

func (s *ConfigurationSuite) TestJsonReadFileForUsers(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	c.Assert(len(config.Users), Equals, 2)
}

func (s *ConfigurationSuite) TestAuthenticationContext(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	c.Assert(config.AuthenticationContextName, Equals, "github")
}

func (s *ConfigurationSuite) TestJsonReadFileForUsersAndValidateError(c *C) {
	configLocation := "fixtures/testconfigwithnousers.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	_, err := NewConfiguration(data)
	c.Assert(err, Not(Equals), nil)
}

func (s *ConfigurationSuite) TestUserWithNoRestriction(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	c.Assert(len(config.Users), Equals, 2)
	c.Assert(config.Users[0].Username, Equals, "KensoDev")
	c.Assert(config.Users[0].Restrict, Equals, "")
	c.Assert(config.Users[1].Restrict, Equals, "GET")
}

func (s *ConfigurationSuite) GetUserRestrictedMethod(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	method := config.GetRestrictionsForUsername("KensoDev2")
	c.Assert(method, Equals, "Get")

	method = config.GetRestrictionsForUsername("KensoDev")
	c.Assert(method, Equals, "")
}

func (s *ConfigurationSuite) GetUserRestrictedMethodNotAllowed(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	method := config.GetRestrictionsForUsername("KensoDev23234234")
	c.Assert(method, Equals, "NotAllowed")
}

func (s *ConfigurationSuite) TestShouldAllowedMethodForUsername(c *C) {
	configLocation := "fixtures/testconfig.json"
	reader := NewConfigurationReader(configLocation)
	data, _ := reader.ReadConfigurationFile()

	config, _ := NewConfiguration(data)
	should := config.ShouldRestrictUser("KensoDev", "POST")
	c.Assert(should, Equals, true)

	should = config.ShouldRestrictUser("KensoDev2", "POST")
	c.Assert(should, Equals, false)
}
