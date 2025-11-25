package undo

import (
	"regexp"
	"strings"
)

// UndoSuggestion represents a suggestion for undoing a command
type UndoSuggestion struct {
	Original    string
	UndoCommand string
	Description string
	CanUndo     bool
}

// undoPatterns maps command patterns to undo suggestions
type undoPattern struct {
	Pattern     *regexp.Regexp
	UndoFunc    func(match []string) string
	Description string
}

var undoPatterns = []undoPattern{
	// Git operations
	{
		Pattern:     regexp.MustCompile(`^git\s+commit\s+`),
		UndoFunc:    func(m []string) string { return "git reset HEAD~1" },
		Description: "Undo the last commit (keeps changes staged)",
	},
	{
		Pattern:     regexp.MustCompile(`^git\s+add\s+(.+)$`),
		UndoFunc:    func(m []string) string { return "git reset " + m[1] },
		Description: "Unstage the added files",
	},
	{
		Pattern:     regexp.MustCompile(`^git\s+stash(\s+push)?$`),
		UndoFunc:    func(m []string) string { return "git stash pop" },
		Description: "Apply and remove the stash",
	},
	{
		Pattern:     regexp.MustCompile(`^git\s+checkout\s+-b\s+(\S+)`),
		UndoFunc:    func(m []string) string { return "git checkout - && git branch -d " + m[1] },
		Description: "Switch back and delete the new branch",
	},
	{
		Pattern:     regexp.MustCompile(`^git\s+merge\s+(\S+)`),
		UndoFunc:    func(m []string) string { return "git reset --hard HEAD~1" },
		Description: "Undo the merge (warning: discards changes)",
	},

	// File operations (when possible)
	{
		Pattern:     regexp.MustCompile(`^mv\s+(\S+)\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "mv " + m[2] + " " + m[1] },
		Description: "Move the file back",
	},
	{
		Pattern:     regexp.MustCompile(`^cp\s+(-[a-z]+\s+)?(\S+)\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "rm " + m[3] },
		Description: "Remove the copied file",
	},
	{
		Pattern:     regexp.MustCompile(`^mkdir\s+(-p\s+)?(\S+)$`),
		UndoFunc:    func(m []string) string { return "rmdir " + m[2] },
		Description: "Remove the created directory (if empty)",
	},
	{
		Pattern:     regexp.MustCompile(`^touch\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "rm " + m[1] },
		Description: "Remove the created file",
	},
	{
		Pattern:     regexp.MustCompile(`^ln\s+(-[a-z]+\s+)?(\S+)\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "rm " + m[3] },
		Description: "Remove the created link",
	},

	// Package manager operations
	{
		Pattern:     regexp.MustCompile(`^(apt|apt-get)\s+install\s+(.+)$`),
		UndoFunc:    func(m []string) string { return m[1] + " remove " + m[2] },
		Description: "Uninstall the package",
	},
	{
		Pattern:     regexp.MustCompile(`^brew\s+install\s+(.+)$`),
		UndoFunc:    func(m []string) string { return "brew uninstall " + m[1] },
		Description: "Uninstall the package",
	},
	{
		Pattern:     regexp.MustCompile(`^npm\s+install\s+(-[gG]\s+)?(.+)$`),
		UndoFunc:    func(m []string) string { return "npm uninstall " + m[1] + m[2] },
		Description: "Uninstall the package",
	},
	{
		Pattern:     regexp.MustCompile(`^pip\s+install\s+(.+)$`),
		UndoFunc:    func(m []string) string { return "pip uninstall " + m[1] },
		Description: "Uninstall the package",
	},

	// Service operations
	{
		Pattern:     regexp.MustCompile(`^systemctl\s+start\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "systemctl stop " + m[1] },
		Description: "Stop the service",
	},
	{
		Pattern:     regexp.MustCompile(`^systemctl\s+stop\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "systemctl start " + m[1] },
		Description: "Start the service",
	},
	{
		Pattern:     regexp.MustCompile(`^systemctl\s+enable\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "systemctl disable " + m[1] },
		Description: "Disable the service",
	},

	// Docker operations
	{
		Pattern:     regexp.MustCompile(`^docker\s+run\s+.*--name\s+(\S+)`),
		UndoFunc:    func(m []string) string { return "docker stop " + m[1] + " && docker rm " + m[1] },
		Description: "Stop and remove the container",
	},
	{
		Pattern:     regexp.MustCompile(`^docker\s+start\s+(\S+)$`),
		UndoFunc:    func(m []string) string { return "docker stop " + m[1] },
		Description: "Stop the container",
	},
}

// GetUndoSuggestion returns an undo suggestion for a command
func GetUndoSuggestion(command string) UndoSuggestion {
	command = strings.TrimSpace(command)

	for _, pattern := range undoPatterns {
		if matches := pattern.Pattern.FindStringSubmatch(command); matches != nil {
			return UndoSuggestion{
				Original:    command,
				UndoCommand: pattern.UndoFunc(matches),
				Description: pattern.Description,
				CanUndo:     true,
			}
		}
	}

	return UndoSuggestion{
		Original:    command,
		UndoCommand: "",
		Description: "No automatic undo available for this command",
		CanUndo:     false,
	}
}

// FormatUndoSuggestion returns a formatted string showing the undo suggestion
func FormatUndoSuggestion(suggestion UndoSuggestion) string {
	if !suggestion.CanUndo {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\033[2m") // Dim
	sb.WriteString("💡 To undo: ")
	sb.WriteString("\033[0m")
	sb.WriteString("\033[36m") // Cyan
	sb.WriteString(suggestion.UndoCommand)
	sb.WriteString("\033[0m")
	sb.WriteString("\n")
	sb.WriteString("\033[2m")
	sb.WriteString("   (")
	sb.WriteString(suggestion.Description)
	sb.WriteString(")")
	sb.WriteString("\033[0m")

	return sb.String()
}

