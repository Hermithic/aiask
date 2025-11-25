package safety

import (
	"regexp"
	"strings"
)

// DangerLevel represents the danger level of a command
type DangerLevel int

const (
	Safe DangerLevel = iota
	Caution
	Dangerous
	Critical
)

// DangerousPattern represents a pattern that indicates a dangerous command
type DangerousPattern struct {
	Pattern     *regexp.Regexp
	Description string
	Level       DangerLevel
}

// dangerousPatterns contains patterns for detecting dangerous commands
var dangerousPatterns = []DangerousPattern{
	// Critical - potentially catastrophic
	{regexp.MustCompile(`(?i)rm\s+(-[a-z]*f[a-z]*\s+)?(-[a-z]*r[a-z]*\s+)?(/|\*|~)`), "Recursive delete of root, all files, or home directory", Critical},
	{regexp.MustCompile(`(?i)rm\s+-[a-z]*rf[a-z]*\s+/`), "Recursive force delete from root", Critical},
	{regexp.MustCompile(`(?i)dd\s+.*of=/dev/(sd|hd|nvme)`), "Direct disk write operation", Critical},
	{regexp.MustCompile(`(?i)mkfs`), "Format filesystem", Critical},
	{regexp.MustCompile(`(?i):(){ :|:& };:`), "Fork bomb", Critical},
	{regexp.MustCompile(`(?i)>\s*/dev/sd`), "Overwrite disk device", Critical},
	{regexp.MustCompile(`(?i)chmod\s+(-[a-z]*R[a-z]*\s+)?777\s+/`), "Set world-writable permissions on root", Critical},

	// Dangerous - significant risk
	{regexp.MustCompile(`(?i)rm\s+-[a-z]*r`), "Recursive delete", Dangerous},
	{regexp.MustCompile(`(?i)rm\s+-[a-z]*f`), "Force delete without confirmation", Dangerous},
	{regexp.MustCompile(`(?i)del\s+/[sq]`), "Windows force/quiet delete", Dangerous},
	{regexp.MustCompile(`(?i)rmdir\s+/s`), "Windows recursive directory delete", Dangerous},
	{regexp.MustCompile(`(?i)drop\s+(table|database|schema)`), "SQL drop operation", Dangerous},
	{regexp.MustCompile(`(?i)truncate\s+table`), "SQL truncate operation", Dangerous},
	{regexp.MustCompile(`(?i)delete\s+from\s+\w+\s*($|;|where\s+1\s*=\s*1)`), "SQL delete without proper WHERE clause", Dangerous},
	{regexp.MustCompile(`(?i)>\s*/etc/`), "Overwrite system config file", Dangerous},
	{regexp.MustCompile(`(?i)curl.*\|\s*(ba)?sh`), "Piping remote content to shell", Dangerous},
	{regexp.MustCompile(`(?i)wget.*\|\s*(ba)?sh`), "Piping remote content to shell", Dangerous},
	{regexp.MustCompile(`(?i)git\s+(push|reset)\s+.*--force`), "Force git operation", Dangerous},
	{regexp.MustCompile(`(?i)git\s+clean\s+-[a-z]*f`), "Force git clean", Dangerous},

	// Caution - requires attention
	{regexp.MustCompile(`(?i)rm\s+`), "Delete operation", Caution},
	{regexp.MustCompile(`(?i)mv\s+.*\s+/dev/null`), "Move to /dev/null (delete)", Caution},
	{regexp.MustCompile(`(?i)chmod\s+`), "Permission change", Caution},
	{regexp.MustCompile(`(?i)chown\s+`), "Ownership change", Caution},
	{regexp.MustCompile(`(?i)kill\s+-9`), "Force kill process", Caution},
	{regexp.MustCompile(`(?i)pkill\s+`), "Kill processes by pattern", Caution},
	{regexp.MustCompile(`(?i)shutdown|reboot|halt|poweroff`), "System shutdown/reboot", Caution},
	{regexp.MustCompile(`(?i)systemctl\s+(stop|disable|mask)`), "Stop/disable system service", Caution},
	{regexp.MustCompile(`(?i)service\s+\w+\s+stop`), "Stop system service", Caution},
	{regexp.MustCompile(`(?i)iptables\s+-F`), "Flush firewall rules", Caution},
	{regexp.MustCompile(`(?i)git\s+reset\s+--hard`), "Hard git reset", Caution},
	{regexp.MustCompile(`(?i)git\s+checkout\s+--\s+\.`), "Discard all changes", Caution},
}

// AnalysisResult represents the result of analyzing a command for danger
type AnalysisResult struct {
	Level       DangerLevel
	Warnings    []string
	IsDangerous bool
}

// Analyze analyzes a command for potential dangers
func Analyze(command string) AnalysisResult {
	result := AnalysisResult{
		Level:       Safe,
		Warnings:    []string{},
		IsDangerous: false,
	}

	for _, pattern := range dangerousPatterns {
		if pattern.Pattern.MatchString(command) {
			result.Warnings = append(result.Warnings, pattern.Description)
			if pattern.Level > result.Level {
				result.Level = pattern.Level
			}
		}
	}

	result.IsDangerous = result.Level >= Dangerous

	return result
}

// GetLevelName returns a human-readable name for the danger level
func GetLevelName(level DangerLevel) string {
	switch level {
	case Safe:
		return "Safe"
	case Caution:
		return "Caution"
	case Dangerous:
		return "Dangerous"
	case Critical:
		return "CRITICAL"
	default:
		return "Unknown"
	}
}

// GetLevelColor returns an ANSI color code for the danger level
func GetLevelColor(level DangerLevel) string {
	switch level {
	case Safe:
		return "\033[32m" // Green
	case Caution:
		return "\033[33m" // Yellow
	case Dangerous:
		return "\033[31m" // Red
	case Critical:
		return "\033[1;31m" // Bold Red
	default:
		return "\033[0m"
	}
}

// RequiresConfirmation returns true if the command requires explicit confirmation
func RequiresConfirmation(command string) bool {
	result := Analyze(command)
	return result.Level >= Dangerous
}

// GetWarningMessage returns a formatted warning message for a dangerous command
func GetWarningMessage(command string) string {
	result := Analyze(command)
	if result.Level < Caution {
		return ""
	}

	var sb strings.Builder
	color := GetLevelColor(result.Level)
	reset := "\033[0m"

	sb.WriteString(color)
	sb.WriteString("⚠️  ")
	sb.WriteString(GetLevelName(result.Level))
	sb.WriteString(" Warning")
	sb.WriteString(reset)
	sb.WriteString("\n")

	for _, warning := range result.Warnings {
		sb.WriteString("   • ")
		sb.WriteString(warning)
		sb.WriteString("\n")
	}

	if result.Level >= Dangerous {
		sb.WriteString("\n")
		sb.WriteString(color)
		sb.WriteString("   Type 'yes' to confirm execution, or any other key to cancel.")
		sb.WriteString(reset)
	}

	return sb.String()
}

