package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitContext represents information about a git repository
type GitContext struct {
	IsRepo       bool
	Branch       string
	IsDirty      bool
	HasUntracked bool
	RemoteURL    string
}

// GetGitContext returns information about the current git repository
func GetGitContext() GitContext {
	ctx := GitContext{}

	// Check if we're in a git repository
	if !IsGitRepo() {
		return ctx
	}
	ctx.IsRepo = true

	// Get current branch
	ctx.Branch = GetGitBranch()

	// Check for uncommitted changes
	ctx.IsDirty = IsGitDirty()

	// Check for untracked files
	ctx.HasUntracked = HasUntrackedFiles()

	// Get remote URL
	ctx.RemoteURL = GetGitRemoteURL()

	return ctx
}

// IsGitRepo checks if the current directory is inside a git repository
func IsGitRepo() bool {
	// Check for .git directory in current or parent directories
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	for {
		gitDir := filepath.Join(cwd, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return true
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return false
}

// GetGitBranch returns the current git branch name
func GetGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// IsGitDirty checks if there are uncommitted changes
func IsGitDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	// Check if there are any changes (excluding untracked files for this check)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) >= 2 && line[0:2] != "??" && strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}

// HasUntrackedFiles checks if there are untracked files
func HasUntrackedFiles() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "??")
}

// GetGitRemoteURL returns the remote origin URL
func GetGitRemoteURL() string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetGitStatus returns a human-readable git status summary
func GetGitStatus() string {
	ctx := GetGitContext()
	if !ctx.IsRepo {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Git repository detected:\n")
	sb.WriteString("  Branch: " + ctx.Branch + "\n")

	if ctx.IsDirty {
		sb.WriteString("  Status: has uncommitted changes\n")
	} else {
		sb.WriteString("  Status: clean\n")
	}

	if ctx.HasUntracked {
		sb.WriteString("  Note: has untracked files\n")
	}

	return sb.String()
}

// IsGitRelatedPrompt checks if a prompt seems to be about git
func IsGitRelatedPrompt(prompt string) bool {
	gitKeywords := []string{
		"git", "commit", "push", "pull", "merge", "rebase", "branch",
		"checkout", "stash", "log", "diff", "status", "clone", "fetch",
		"remote", "tag", "reset", "revert", "cherry-pick", "squash",
	}

	promptLower := strings.ToLower(prompt)
	for _, keyword := range gitKeywords {
		if strings.Contains(promptLower, keyword) {
			return true
		}
	}
	return false
}

// GetRecentCommits returns the last N commit summaries
func GetRecentCommits(n int) string {
	if !IsGitRepo() {
		return ""
	}

	cmd := exec.Command("git", "log", "--oneline", "-n", string(rune('0'+n)))
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

