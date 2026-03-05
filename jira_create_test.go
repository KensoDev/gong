package gong

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

type JiraCreateSuite struct{}

var _ = Suite(&JiraCreateSuite{})

func TestJiraCreate(t *testing.T) { TestingT(t) }

func (s *JiraCreateSuite) TestDetermineBranchType(c *C) {
	client := NewJiraClient()

	testCases := []struct {
		issueType    string
		expectedType string
	}{
		{"Bug", "bugfix"},
		{"bug", "bugfix"},
		{"Defect", "bugfix"},
		{"Story", "feature"},
		{"Task", "feature"},
		{"Feature", "feature"},
		{"Epic", "epic"},
		{"Improvement", "enhancement"},
		{"Enhancement", "enhancement"},
		{"Sub-task", "task"},
		{"Subtask", "task"},
		{"Unknown Type", "feature"}, // default
	}

	for _, tc := range testCases {
		result := client.determineBranchType(tc.issueType)
		c.Assert(result, Equals, tc.expectedType, Commentf("Issue type: %s", tc.issueType))
	}
}

func (s *JiraCreateSuite) TestDetermineBranchTypeWithEnvVar(c *C) {
	client := NewJiraClient()

	// Set environment variable
	originalEnv := os.Getenv("GONG_DEFAULT_BRANCH_TYPE")
	defer func() {
		if originalEnv != "" {
			os.Setenv("GONG_DEFAULT_BRANCH_TYPE", originalEnv)
		} else {
			os.Unsetenv("GONG_DEFAULT_BRANCH_TYPE")
		}
	}()

	os.Setenv("GONG_DEFAULT_BRANCH_TYPE", "custom")

	// Should use env var regardless of issue type
	result := client.determineBranchType("Bug")
	c.Assert(result, Equals, "custom")

	result = client.determineBranchType("Story")
	c.Assert(result, Equals, "custom")
}

func (s *JiraCreateSuite) TestCreateGitBranchValidation(c *C) {
	client := NewJiraClient()

	// Test with empty branch name (should fail)
	err := client.createGitBranch("")
	c.Assert(err, NotNil)
}

func (s *JiraCreateSuite) TestGetIssueIDFromBranchName(c *C) {
	testCases := []struct {
		branchName string
		expectedID string
	}{
		{"feature/PROJ-123-some-feature", "PROJ-123"},
		{"bugfix/ABC-456-fix-bug", "ABC-456"},
		{"feature/XYZ-1-single-digit", "XYZ-1"},
		{"PROJ-789-no-prefix", "PROJ-789"},
		{"feature/no-ticket-here", ""},
		{"main", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		result := GetIssueID(tc.branchName)
		c.Assert(result, Equals, tc.expectedID, Commentf("Branch: %s", tc.branchName))
	}
}

func (s *JiraCreateSuite) TestCreateRequiresProjectKey(c *C) {
	client := NewJiraClient()

	// Should fail with empty project key
	_, err := client.Create("")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Matches, ".*project key is required.*")
}
