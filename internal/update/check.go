package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
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

// parseVersion splits a version string into its core version and pre-release parts
// e.g., "2.0.0-beta.1" -> core: "2.0.0", prerelease: "beta.1"
func parseVersion(version string) (core string, prerelease string) {
	parts := strings.SplitN(version, "-", 2)
	core = parts[0]
	if len(parts) > 1 {
		prerelease = parts[1]
	}
	return
}

// compareVersionCore compares two version core strings (e.g., "2.0.0" vs "2.1.0")
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareVersionCore(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	// Pad shorter version with zeros
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			fmt.Sscanf(aParts[i], "%d", &aNum)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bNum)
		}

		if aNum > bNum {
			return 1
		}
		if aNum < bNum {
			return -1
		}
	}
	return 0
}

// comparePrerelease compares two pre-release strings using semantic versioning rules
// Each dot-separated part is compared numerically if both are numbers, otherwise lexicographically
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func comparePrerelease(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		// If one has fewer parts, it comes first (per semver spec)
		if i >= len(aParts) {
			return -1
		}
		if i >= len(bParts) {
			return 1
		}

		aPart := aParts[i]
		bPart := bParts[i]

		// Try to parse both as numbers
		aNum, aErr := strconv.Atoi(aPart)
		bNum, bErr := strconv.Atoi(bPart)

		if aErr == nil && bErr == nil {
			// Both are numeric, compare as numbers
			if aNum > bNum {
				return 1
			}
			if aNum < bNum {
				return -1
			}
		} else if aErr == nil {
			// a is numeric, b is not - numeric comes first per semver
			return -1
		} else if bErr == nil {
			// b is numeric, a is not - numeric comes first per semver
			return 1
		} else {
			// Both are strings, compare lexicographically
			if aPart > bPart {
				return 1
			}
			if aPart < bPart {
				return -1
			}
		}
	}

	return 0
}

// isNewer compares version strings with support for pre-release versions
// Follows semantic versioning: release versions are newer than pre-release versions
// e.g., 2.0.0 > 2.0.0-beta.1, 2.0.0-beta.2 > 2.0.0-beta.1, 2.0.0-beta.10 > 2.0.0-beta.2
func isNewer(latest, current string) bool {
	if current == "dev" || current == "" {
		return false
	}

	latestCore, latestPrerelease := parseVersion(latest)
	currentCore, currentPrerelease := parseVersion(current)

	// Compare core versions first
	coreComparison := compareVersionCore(latestCore, currentCore)
	if coreComparison != 0 {
		return coreComparison > 0
	}

	// Core versions are equal, compare pre-release parts
	// A release version (no prerelease) is newer than a pre-release version
	if latestPrerelease == "" && currentPrerelease != "" {
		return true // latest is a release, current is a pre-release
	}
	if latestPrerelease != "" && currentPrerelease == "" {
		return false // latest is a pre-release, current is a release
	}
	if latestPrerelease == "" && currentPrerelease == "" {
		return false // both are releases, same version
	}

	// Both have pre-release parts, compare them properly
	return comparePrerelease(latestPrerelease, currentPrerelease) > 0
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

