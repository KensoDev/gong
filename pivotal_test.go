package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestPivotalClient(t *testing.T) { TestingT(t) }

type PivotalClientSuit struct{}

var _ = Suite(&PivotalClientSuit{})

func (p *PivotalClientSuit) TestBrowseWithCorrectBranchName(c * C) {
	pivotalClient := &PivotalClient{}

	url, err := pivotalClient.Browse("feature/124352-test-only")
	c.Assert(err, Equals, nil)
	c.Assert(url, Equals, "https://www.pivotaltracker.com/story/show/124352")
}

func (p *PivotalClientSuit) TestBrowseWithIncorrectBranchName(c * C) {
	pivotalClient := &PivotalClient{}

	url, err := pivotalClient.Browse("feature/test-only")
	c.Assert(err, Equals, nil)
	c.Assert(url, Equals, "https://www.pivotaltracker.com/story/show/")
}

func (p *PivotalClientSuit) TestGetPivotalIssueIDWithCorrectBranchName(c * C) {
	branchName := "feature/1234-test-only"
	issueId := GetPivotalIssueID(branchName)
	c.Assert(issueId, Equals, "1234")
}


func (p *PivotalClientSuit) TestGetPivotalIssueIDWithInCorrectBranchName(c * C) {
	branchName := "feature/test-only"
	issueId := GetPivotalIssueID(branchName)
	c.Assert(issueId, Equals, "")
}
