package authproxy

import (
	"io/ioutil"
	"testing"

	. "gopkg.in/check.v1"
)

func TestGithubAuthContext(t *testing.T) { TestingT(t) }

type GithubAuthContextSuite struct{}

var _ = Suite(&GithubAuthContextSuite{})

func (s *GithubAuthContextSuite) TestGithubUserParsing(c *C) {
	data, err := ioutil.ReadFile("fixtures/github-user.json")
	configuration := &Configuration{}
	authContext := NewGithubAuthContext(configuration)
	user, err := authContext.ParseUserResponse(data)
	c.Assert(err, IsNil)
	c.Assert(user.UserName, Equals, "KensoDev")
}
