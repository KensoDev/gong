package gong

import (
	"os"
	"os/exec"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type HooksSuite struct {
	tempDir string
}

var _ = Suite(&HooksSuite{})

func (s *HooksSuite) SetUpTest(c *C) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gong-hooks-test-*")
	c.Assert(err, IsNil)
	s.tempDir = tempDir

	// Initialize a git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err = cmd.Run()
	c.Assert(err, IsNil)

	// Change to temp directory
	err = os.Chdir(tempDir)
	c.Assert(err, IsNil)
}

func (s *HooksSuite) TearDownTest(c *C) {
	// Clean up
	os.RemoveAll(s.tempDir)
}

func (s *HooksSuite) TestIsGitRepo(c *C) {
	// Should be true in our test git repo
	result := isGitRepo()
	c.Assert(result, Equals, true)

	// Create a non-git directory
	nonGitDir, err := os.MkdirTemp("", "non-git-*")
	c.Assert(err, IsNil)
	defer os.RemoveAll(nonGitDir)

	err = os.Chdir(nonGitDir)
	c.Assert(err, IsNil)

	result = isGitRepo()
	c.Assert(result, Equals, false)

	// Change back
	os.Chdir(s.tempDir)
}

func (s *HooksSuite) TestInstallHooksNewHook(c *C) {
	// Install hooks in fresh repo
	err := InstallHooks()
	c.Assert(err, IsNil)

	// Check that hook file was created
	hookPath := filepath.Join(".git", "hooks", "prepare-commit-msg")
	_, err = os.Stat(hookPath)
	c.Assert(err, IsNil)

	// Check that file is executable
	info, err := os.Stat(hookPath)
	c.Assert(err, IsNil)
	mode := info.Mode()
	c.Assert(mode&0111 != 0, Equals, true) // Has execute permission

	// Check content
	content, err := os.ReadFile(hookPath)
	c.Assert(err, IsNil)
	contentStr := string(content)

	c.Assert(contentStr, Matches, "(?s).*#!/bin/sh.*")
	c.Assert(contentStr, Matches, "(?s).*# Gong: Auto-add JIRA ticket ID.*")
	c.Assert(contentStr, Matches, "(?s).*gong prepare-commit-message.*")
}

func (s *HooksSuite) TestInstallHooksExistingHook(c *C) {
	// Create existing hook
	hooksDir := filepath.Join(".git", "hooks")
	os.MkdirAll(hooksDir, 0755)
	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	existingContent := `#!/bin/sh
# My existing hook
echo "Running my hook"
`
	err := os.WriteFile(hookPath, []byte(existingContent), 0755)
	c.Assert(err, IsNil)

	// Install gong hooks
	err = InstallHooks()
	c.Assert(err, IsNil)

	// Check that both old and new content exist
	content, err := os.ReadFile(hookPath)
	c.Assert(err, IsNil)
	contentStr := string(content)

	c.Assert(contentStr, Matches, "(?s).*# My existing hook.*")
	c.Assert(contentStr, Matches, "(?s).*echo \"Running my hook\".*")
	c.Assert(contentStr, Matches, "(?s).*# Gong: Auto-add JIRA ticket ID.*")
	c.Assert(contentStr, Matches, "(?s).*gong prepare-commit-message.*")
}

func (s *HooksSuite) TestInstallHooksIdempotent(c *C) {
	// Install hooks once
	err := InstallHooks()
	c.Assert(err, IsNil)

	hookPath := filepath.Join(".git", "hooks", "prepare-commit-msg")
	content1, err := os.ReadFile(hookPath)
	c.Assert(err, IsNil)

	// Install again
	err = InstallHooks()
	c.Assert(err, IsNil)

	// Content should be the same (not duplicated)
	content2, err := os.ReadFile(hookPath)
	c.Assert(err, IsNil)

	c.Assert(string(content1), Equals, string(content2))
}

func (s *HooksSuite) TestHasGongHook(c *C) {
	hookPath := filepath.Join(".git", "hooks", "prepare-commit-msg")

	// No hook file
	result := hasGongHook(hookPath)
	c.Assert(result, Equals, false)

	// Create hook without gong
	hooksDir := filepath.Join(".git", "hooks")
	os.MkdirAll(hooksDir, 0755)
	err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho test"), 0755)
	c.Assert(err, IsNil)

	result = hasGongHook(hookPath)
	c.Assert(result, Equals, false)

	// Create hook with gong
	err = os.WriteFile(hookPath, []byte("#!/bin/sh\n# Gong: Auto-add JIRA ticket ID\necho test"), 0755)
	c.Assert(err, IsNil)

	result = hasGongHook(hookPath)
	c.Assert(result, Equals, true)
}

func (s *HooksSuite) TestInstallHooksNotGitRepo(c *C) {
	// Create a non-git directory
	nonGitDir, err := os.MkdirTemp("", "non-git-*")
	c.Assert(err, IsNil)
	defer os.RemoveAll(nonGitDir)

	err = os.Chdir(nonGitDir)
	c.Assert(err, IsNil)

	// Should return error
	err = InstallHooks()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Matches, ".*not a git repository.*")

	// Change back
	os.Chdir(s.tempDir)
}
