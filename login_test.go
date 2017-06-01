package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestLogin(t *testing.T) { TestingT(t) }

type LoginSuite struct{}

var _ = Suite(&LoginSuite{})

func (s *LoginSuite) TestLoginDetailsLoad(c *C) {

	loginDetails := LoginDetails{
		Directory: "test/fixtures",
	}

	loginDetails, _ = Load(loginDetails)

	c.Assert(loginDetails.Username, Equals, "test-username")
	c.Assert(loginDetails.Password, Equals, "test-password")
	c.Assert(loginDetails.Domain, Equals, "https://test-domain")
}
