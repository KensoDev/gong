package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestJiraClient(t *testing.T) { TestingT(t) }

type JiraClientSuite struct{}

var _ = Suite(&JiraClientSuite{})

func (s *JiraClientSuite) TestBrowseWithCorrectConfig(c *C) {
	config := map[string]string{
		"domain": "https://fake.atlassian.net",
	}

	jiraClient := &JiraClient{
		config: config,
	}

	url, err := jiraClient.Browse("feature/FAKE-1111-something-something")
	c.Assert(err, Equals, nil)
	c.Assert(url, Equals, "https://fake.atlassian.net/browse/FAKE-1111")
}

func (s *JiraClientSuite) TestBrowseWithIncorrectConfig(c *C) {
	config := map[string]string{}

	jiraClient := &JiraClient{
		config: config,
	}

	url, err := jiraClient.Browse("feature/FAKE-1111-something-something")
	c.Assert(err, NotNil)
	c.Assert(url, Equals, "")
}

func (s *JiraClientSuite) TestBrowseWithCorrectConfigButIncorrectBranchName(c *C) {
	config := map[string]string{
		"domain": "https://fake.atlassian.net",
	}

	jiraClient := &JiraClient{
		config: config,
	}

	url, err := jiraClient.Browse("feature/something-something")
	c.Assert(err, NotNil)
	c.Assert(url, Equals, "")
}
