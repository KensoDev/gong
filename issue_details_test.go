package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestIssueDetails(t *testing.T) { TestingT(t) }

type IssueDetailsSuite struct{}

var _ = Suite(&IssueDetailsSuite{})

func (s *IssueDetailsSuite) TestGetIssueIDFromBranchName(c *C) {
	branchName := "feature/GLOB-1111-build-this-GID-1234hh"
	issueId, _ := GetIssueID(branchName)

	c.Assert(issueId, Equals, "GLOB-1111")
}
