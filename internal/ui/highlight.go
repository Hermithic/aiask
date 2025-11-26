package ui

import (
	"regexp"
	"strings"
)

// Highlighter provides syntax highlighting for shell commands
type Highlighter struct {
	patterns []highlightPattern
}

type highlightPattern struct {
	pattern *regexp.Regexp
	color   string
}

// NewHighlighter creates a new syntax highlighter
func NewHighlighter() *Highlighter {
	return &Highlighter{
		patterns: []highlightPattern{
			// Strings (double quotes)
			{regexp.MustCompile(`"[^"]*"`), "\033[33m"},       // Yellow
			// Strings (single quotes)
			{regexp.MustCompile(`'[^']*'`), "\033[33m"},       // Yellow
			// Comments
			{regexp.MustCompile(`#.*$`), "\033[90m"},          // Gray
			// Flags (short and long)
			{regexp.MustCompile(`\s(-{1,2}[a-zA-Z0-9_-]+)`), "\033[36m"}, // Cyan
			// Variables
			{regexp.MustCompile(`\$[a-zA-Z_][a-zA-Z0-9_]*`), "\033[35m"}, // Magenta
			{regexp.MustCompile(`\$\{[^}]+\}`), "\033[35m"},              // Magenta
			// Numbers
			{regexp.MustCompile(`\b\d+\b`), "\033[34m"},       // Blue
			// Paths
			{regexp.MustCompile(`(?:^|\s)((?:/|\.\.?/)[a-zA-Z0-9._/-]+)`), "\033[32m"}, // Green
			// Pipes and redirects
			{regexp.MustCompile(`[|><&]+`), "\033[31m"},       // Red
			// Glob patterns
			{regexp.MustCompile(`\*+|\?+`), "\033[33m"},       // Yellow
		},
	}
}

// Highlight applies syntax highlighting to a command
func (h *Highlighter) Highlight(command string) string {
	// First, identify all the regions that match patterns
	type match struct {
		start, end int
		color      string
	}

	var matches []match
	for _, p := range h.patterns {
		locs := p.pattern.FindAllStringIndex(command, -1)
		for _, loc := range locs {
			matches = append(matches, match{loc[0], loc[1], p.color})
		}
	}

	// Sort matches by start position (for proper ordering)
	// Simple bubble sort for small arrays
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].start < matches[i].start {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	// Remove overlapping matches (keep earlier ones)
	var filtered []match
	lastEnd := 0
	for _, m := range matches {
		if m.start >= lastEnd {
			filtered = append(filtered, m)
			lastEnd = m.end
		}
	}

	// Build the highlighted string
	if len(filtered) == 0 {
		return command
	}

	var result strings.Builder
	pos := 0
	for _, m := range filtered {
		// Add unhighlighted text before this match
		if m.start > pos {
			result.WriteString(command[pos:m.start])
		}
		// Add highlighted match
		result.WriteString(m.color)
		result.WriteString(command[m.start:m.end])
		result.WriteString(ColorReset)
		pos = m.end
	}
	// Add remaining text
	if pos < len(command) {
		result.WriteString(command[pos:])
	}

	return result.String()
}

// FormatKeyword highlights a keyword
func FormatKeyword(s string) string {
	return ColorBold + ColorBlue + s + ColorReset
}

// FormatPath highlights a file path
func FormatPath(s string) string {
	return ColorGreen + s + ColorReset
}

// FormatFlag highlights a command flag
func FormatFlag(s string) string {
	return ColorCyan + s + ColorReset
}

// FormatString highlights a string literal
func FormatString(s string) string {
	return ColorYellow + s + ColorReset
}

// FormatVariable highlights a variable
func FormatVariable(s string) string {
	return "\033[35m" + s + ColorReset // Magenta
}

// FormatNumber highlights a number
func FormatNumber(s string) string {
	return ColorBlue + s + ColorReset
}

