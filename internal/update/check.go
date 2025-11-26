package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHubReleaseURL is the URL to check for releases
	GitHubReleaseURL = "https://api.github.com/repos/Hermithic/aiask/releases/latest"
	// CheckTimeout is the timeout for update checks
	CheckTimeout = 5 * time.Second
)

// ReleaseInfo contains information about a GitHub release
type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

// CheckResult contains the result of an update check
type CheckResult struct {
	UpdateAvailable bool
	CurrentVersion  string
	LatestVersion   string
	ReleaseURL      string
	ReleaseNotes    string
}

// CheckForUpdates checks if a newer version is available
func CheckForUpdates(currentVersion string) (*CheckResult, error) {
	client := &http.Client{
		Timeout: CheckTimeout,
	}

	req, err := http.NewRequest("GET", GitHubReleaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "aiask-update-checker")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentClean := strings.TrimPrefix(currentVersion, "v")

	return &CheckResult{
		UpdateAvailable: isNewer(latestVersion, currentClean),
		CurrentVersion:  currentVersion,
		LatestVersion:   release.TagName,
		ReleaseURL:      release.HTMLURL,
		ReleaseNotes:    truncateReleaseNotes(release.Body, 200),
	}, nil
}

// CheckForUpdatesAsync checks for updates in the background
func CheckForUpdatesAsync(currentVersion string, callback func(*CheckResult)) {
	go func() {
		result, err := CheckForUpdates(currentVersion)
		if err == nil && result != nil {
			callback(result)
		}
	}()
}

// isNewer compares version strings (simple comparison)
func isNewer(latest, current string) bool {
	if current == "dev" || current == "" {
		return false
	}

	// Simple version comparison (works for semver-like versions)
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		var latestNum, currentNum int
		fmt.Sscanf(latestParts[i], "%d", &latestNum)
		fmt.Sscanf(currentParts[i], "%d", &currentNum)

		if latestNum > currentNum {
			return true
		}
		if latestNum < currentNum {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}

// truncateReleaseNotes truncates release notes to a maximum length
func truncateReleaseNotes(notes string, maxLen int) string {
	if len(notes) <= maxLen {
		return notes
	}
	return notes[:maxLen] + "..."
}

// getUpdateCommand returns the platform-specific update command
func getUpdateCommand() string {
	switch runtime.GOOS {
	case "windows":
		return "winget upgrade Hermithic.aiask"
	case "darwin":
		return "brew upgrade aiask"
	case "linux":
		return "sudo apt update && sudo apt upgrade aiask"
	default:
		return "download from GitHub releases"
	}
}

// FormatUpdateMessage formats a message about available update
func FormatUpdateMessage(result *CheckResult) string {
	if !result.UpdateAvailable {
		return ""
	}

	return fmt.Sprintf(
		"\033[33m╔════════════════════════════════════════════╗\n"+
			"║  Update available: %s → %s\n"+
			"║  Run: \033[36m%s\033[33m\n"+
			"║  Or visit: %s\n"+
			"╚════════════════════════════════════════════╝\033[0m\n",
		result.CurrentVersion, result.LatestVersion,
		getUpdateCommand(),
		result.ReleaseURL,
	)
}

