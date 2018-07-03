package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestBranch(t *testing.T) { TestingT(t) }

type BranchSuite struct{}

var _ = Suite(&BranchSuite{})

const issueType = "type"
const issueID = "id"
const issueTitle = "title somewhat long"

func (b *BranchSuite) TestDefaultPatternAndDefaultReplacementCharacterWhenConfigFieldsMissing(c *C) {
	config := map[string]string{}
	branch, err := NewBranch(config)

	branchName, err := branch.Name(issueType, issueID, issueTitle)

	c.Assert(err, Equals, nil)
	c.Assert(branchName, Equals, "type/id-title-somewhat-long")
}

func (b *BranchSuite) TestUseReplacementCharacterWhenConfigFieldPresent(c *C) {
	config := map[string]string {
		"branch_replacement_character": "_",
	}
	branch, err := NewBranch(config)

	branchName, err := branch.Name(issueType, issueID, issueTitle)

	c.Assert(err, Equals, nil)
	c.Assert(branchName, Equals, "type/id-title_somewhat_long")
}

func (b *BranchSuite) TestUsePatternWhenConfigFieldPresent(c *C) {
	config := map[string]string {
		"branch_pattern": "{{.IssueTitle}}+{{.IssueID}}-{{.IssueType}}",
	}
	branch, err := NewBranch(config)

	branchName, err := branch.Name(issueType, issueID, issueTitle)

	c.Assert(err, Equals, nil)
	c.Assert(branchName, Equals, "title-somewhat-long+id-type")
}

func (b *BranchSuite) TestFailWhenIssueTitleNotPresentInPattern(c *C) {
	config := map[string]string {
		"branch_pattern": "{{.IssueType}}/{{.IssueID}}",
	}

	_, err := NewBranch(config)

	c.Assert(err, NotNil)
}

func (b *BranchSuite) TestFailWhenIssueIDNotPresentInPattern(c *C) {
	config := map[string]string {
		"branch_pattern": "{{.IssueTitle}}+{{.IssueType}}",
	}

	_, err := NewBranch(config)

	c.Assert(err, 	NotNil)
}
