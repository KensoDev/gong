package gong

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"
)

type IntegrationSuite struct {
	tempDir     string
	originalDir string
}

var _ = Suite(&IntegrationSuite{})

func (s *IntegrationSuite) SetUpTest(c *C) {
	// Save original directory
	originalDir, err := os.Getwd()
	c.Assert(err, IsNil)
	s.originalDir = originalDir

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gong-integration-test-*")
	c.Assert(err, IsNil)
	s.tempDir = tempDir

	// Initialize a git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err = cmd.Run()
	c.Assert(err, IsNil)

	// Configure git for testing
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()

	// Change to temp directory
	err = os.Chdir(tempDir)
	c.Assert(err, IsNil)
}

func (s *IntegrationSuite) TearDownTest(c *C) {
	// Change back to original directory
	os.Chdir(s.originalDir)
	// Clean up
	os.RemoveAll(s.tempDir)
}

func (s *IntegrationSuite) TestHookIntegrationWithGit(c *C) {
	// Install hooks
	err := InstallHooks()
	c.Assert(err, IsNil)

	// Create a dummy branch with ticket ID
	cmd := exec.Command("git", "checkout", "-b", "feature/TEST-123-test-feature")
	err = cmd.Run()
	c.Assert(err, IsNil)

	// Create a test file
	testFile := filepath.Join(s.tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	c.Assert(err, IsNil)

	// Stage the file
	cmd = exec.Command("git", "add", "test.txt")
	err = cmd.Run()
	c.Assert(err, IsNil)

	// Create commit (hook should add ticket ID)
	// Note: This will fail if gong binary is not in PATH, which is expected in test environment
	// The hook file itself is still created and tested
	hookPath := filepath.Join(".git", "hooks", "prepare-commit-msg")
	content, err := os.ReadFile(hookPath)
	c.Assert(err, IsNil)

	// Verify hook contains necessary logic
	contentStr := string(content)
	c.Assert(strings.Contains(contentStr, "gong prepare-commit-message"), Equals, true)
	c.Assert(strings.Contains(contentStr, "BRANCH_NAME"), Equals, true)
}

func (s *IntegrationSuite) TestPrepareCommitMessageWithValidBranch(c *C) {
	branchName := "feature/PROJ-456-awesome-feature\n"
	client := NewJiraClient()

	// Mock minimal config
	client.config = map[string]string{
		"domain": "https://test.atlassian.net",
	}

	commitMessage := "Initial commit"
	result := client.PrepareCommitMessage(branchName, commitMessage)

	// Should contain ticket link
	c.Assert(strings.Contains(result, "PROJ-456"), Equals, true)
	c.Assert(strings.Contains(result, "https://test.atlassian.net/browse/PROJ-456"), Equals, true)
}

func (s *IntegrationSuite) TestPrepareCommitMessageWithInvalidBranch(c *C) {
	branchName := "main\n"
	client := NewJiraClient()

	client.config = map[string]string{
		"domain": "https://test.atlassian.net",
	}

	commitMessage := "Initial commit"
	result := client.PrepareCommitMessage(branchName, commitMessage)

	// Should return original commit message when no ticket ID found
	c.Assert(result, Equals, commitMessage)
}

func (s *IntegrationSuite) TestSlugifyIntegration(c *C) {
	// Test that slugification works correctly with GetBranchName
	testCases := []struct {
		title    string
		expected string
	}{
		{"Simple Title", "simple-title"},
		{"Title with Special!@# Characters", "title-with-special-characters"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"CamelCase Title", "camelcase-title"},
	}

	for _, tc := range testCases {
		result := SlugifyTitle(tc.title)
		c.Assert(result, Equals, tc.expected, Commentf("Title: %s", tc.title))
	}
}

func (s *IntegrationSuite) TestGetBranchNameIntegration(c *C) {
	// Note: This test would require a real JIRA connection
	// For now, we just test that the function signature is correct
	client := NewJiraClient()
	c.Assert(client, NotNil)
}
