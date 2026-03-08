package gong

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

const prepareCommitMsgHook = `
# Gong: Auto-add JIRA ticket ID to commit messages
if command -v gong >/dev/null 2>&1; then
    BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
    if [ -n "$BRANCH_NAME" ]; then
        TICKET_LINK=$(echo "$BRANCH_NAME" | gong prepare-commit-message)
        if [ -n "$TICKET_LINK" ]; then
            # Append ticket link if not already present (use -F for literal string match)
            if ! grep -F -q "$TICKET_LINK" "$1"; then
                echo "" >> "$1"
                echo "$TICKET_LINK" >> "$1"
            fi
        fi
    fi
fi
`

// CheckAndInstallHooks checks if git hooks are installed and prompts to install if not
func CheckAndInstallHooks() error {
	// Check if we're in a git repo
	if !isGitRepo() {
		return nil // Silently skip if not in a git repo
	}

	hookPath := filepath.Join(".git", "hooks", "prepare-commit-msg")

	// Check if hook already has gong integration
	if hasGongHook(hookPath) {
		return nil // Already installed
	}

	// Prompt user
	prompt := promptui.Prompt{
		Label:     "Git hook not installed. Install prepare-commit-msg hook to auto-add ticket IDs to commits? (y/n)",
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil || strings.ToLower(result) != "y" {
		color.Yellow("Skipping hook installation. Run 'gong install-hooks' later to install.")
		return nil
	}

	return InstallHooks()
}

// InstallHooks installs or updates the git hooks
func InstallHooks() error {
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	hooksDir := filepath.Join(".git", "hooks")
	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	// Ensure hooks directory exists
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Check if hook already exists
	existingContent := ""
	hasExistingGong := false

	if _, err := os.Stat(hookPath); err == nil {
		// Hook exists, read it
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("failed to read existing hook: %w", err)
		}
		existingContent = string(content)
		hasExistingGong = hasGongHook(hookPath)
	}

	// Create or update hook
	var finalContent string

	if existingContent == "" {
		// New hook
		finalContent = "#!/bin/sh\n" + prepareCommitMsgHook
		color.Yellow("Installing git hook...")
	} else if hasExistingGong {
		// Update existing gong hook - remove old gong section and add new one
		color.Yellow("Updating existing gong hook...")
		// Remove old gong section (everything from "# Gong:" marker to end of file or next non-gong content)
		lines := strings.Split(existingContent, "\n")
		var filteredLines []string
		skipGong := false

		for _, line := range lines {
			if strings.Contains(line, "# Gong: Auto-add JIRA ticket ID") {
				skipGong = true
				continue
			}
			if skipGong {
				// Skip until we find a non-empty line that doesn't look like our hook
				if strings.TrimSpace(line) == "" ||
				   strings.Contains(line, "BRANCH_NAME") ||
				   strings.Contains(line, "TICKET_LINK") ||
				   strings.Contains(line, "gong prepare-commit-message") ||
				   strings.Contains(line, "command -v gong") {
					continue
				}
				skipGong = false
			}
			filteredLines = append(filteredLines, line)
		}

		finalContent = strings.Join(filteredLines, "\n") + "\n" + prepareCommitMsgHook
	} else {
		// Append to existing hook (no gong yet)
		color.Yellow("Found existing prepare-commit-msg hook. Appending gong integration...")
		finalContent = existingContent + "\n" + prepareCommitMsgHook
	}

	// Write hook
	if err := os.WriteFile(hookPath, []byte(finalContent), 0755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	color.Green("✓ Git hook installed successfully!")
	color.Green("  All commits will now automatically include the JIRA ticket ID")
	return nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// hasGongHook checks if the hook file already contains gong integration
func hasGongHook(hookPath string) bool {
	content, err := os.ReadFile(hookPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "# Gong: Auto-add JIRA ticket ID")
}
