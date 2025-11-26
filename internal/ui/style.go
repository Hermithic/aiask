package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Box drawing characters for consistent styling
const (
	BoxTopLeft     = "╭"
	BoxTopRight    = "╮"
	BoxBottomLeft  = "╰"
	BoxBottomRight = "╯"
	BoxHorizontal  = "─"
	BoxVertical    = "│"

	// Heavy box for headers
	HeavyTopLeft     = "┏"
	HeavyTopRight    = "┓"
	HeavyBottomLeft  = "┗"
	HeavyBottomRight = "┛"
	HeavyHorizontal  = "━"
	HeavyVertical    = "┃"

	// Icons
	IconCheck    = "✓"
	IconCross    = "✗"
	IconArrow    = "→"
	IconDot      = "●"
	IconCircle   = "○"
	IconWarning  = "⚠"
	IconInfo     = "ℹ"
	IconRocket   = "🚀"
	IconCopy     = "📋"
	IconEdit     = "✏"
	IconRefresh  = "↻"
	IconExit     = "⏻"
	IconTerminal = "❯"
	IconStar     = "★"
	IconTime     = "⏱"
)

// SelectTemplates returns styled templates for promptui Select
func SelectTemplates(labelText string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    fmt.Sprintf("%s{{ . }}%s", ColorBold, ColorReset),
		Active:   fmt.Sprintf("%s%s {{ .Label | cyan | bold }}%s  {{ .Desc | faint }}", ColorCyan, IconArrow, ColorReset),
		Inactive: "  {{ .Label }}  {{ .Desc | faint }}",
		Selected: fmt.Sprintf("%s%s {{ .Label | green }}%s", ColorGreen, IconCheck, ColorReset),
		Details: `
{{ "─────────────────────────────────────────" | faint }}
{{ .Desc | faint }}`,
	}
}

// SelectTemplatesSimple returns simpler templates without details section
func SelectTemplatesSimple() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    fmt.Sprintf("%s{{ . }}%s", ColorBold, ColorReset),
		Active:   fmt.Sprintf("%s%s {{ .Label | cyan | bold }}%s", ColorCyan, IconArrow, ColorReset),
		Inactive: "  {{ .Label }}",
		Selected: fmt.Sprintf("%s%s {{ .Label | green }}%s", ColorGreen, IconCheck, ColorReset),
	}
}

// PromptTemplates returns styled templates for promptui Prompt
func PromptTemplates() *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:  fmt.Sprintf("%s{{ . }}%s ", ColorBold, ColorReset),
		Valid:   fmt.Sprintf("%s{{ . }}%s ", ColorGreen, ColorReset),
		Invalid: fmt.Sprintf("%s{{ . }}%s ", ColorRed, ColorReset),
		Success: fmt.Sprintf("%s%s {{ . }}%s ", ColorGreen, IconCheck, ColorReset),
	}
}

// MenuItem represents a menu option
type MenuItem struct {
	Label string
	Desc  string
	Key   string // Keyboard shortcut
	Icon  string
}

// ActionMenuItem creates a menu item for actions
func ActionMenuItem(label, desc, key, icon string) MenuItem {
	return MenuItem{
		Label: fmt.Sprintf("%s %s", icon, label),
		Desc:  desc,
		Key:   key,
		Icon:  icon,
	}
}

// Box creates a styled box around text
func Box(title, content string, width int) string {
	var sb strings.Builder

	// Ensure minimum width
	if width < len(title)+4 {
		width = len(title) + 4
	}

	// Top border
	sb.WriteString(ColorCyan)
	sb.WriteString(BoxTopLeft)
	sb.WriteString(strings.Repeat(BoxHorizontal, width-2))
	sb.WriteString(BoxTopRight)
	sb.WriteString(ColorReset)
	sb.WriteString("\n")

	// Title line
	if title != "" {
		padding := width - 4 - len(title)
		leftPad := padding / 2
		rightPad := padding - leftPad
		sb.WriteString(ColorCyan)
		sb.WriteString(BoxVertical)
		sb.WriteString(ColorReset)
		sb.WriteString(" ")
		sb.WriteString(strings.Repeat(" ", leftPad))
		sb.WriteString(ColorBold)
		sb.WriteString(title)
		sb.WriteString(ColorReset)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(" ")
		sb.WriteString(ColorCyan)
		sb.WriteString(BoxVertical)
		sb.WriteString(ColorReset)
		sb.WriteString("\n")
	}

	// Content lines
	if content != "" {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			sb.WriteString(ColorCyan)
			sb.WriteString(BoxVertical)
			sb.WriteString(ColorReset)
			sb.WriteString(" ")

			// Truncate or pad line to fit
			if len(line) > width-4 {
				line = line[:width-7] + "..."
			}
			sb.WriteString(line)
			sb.WriteString(strings.Repeat(" ", width-4-len(line)))
			sb.WriteString(" ")
			sb.WriteString(ColorCyan)
			sb.WriteString(BoxVertical)
			sb.WriteString(ColorReset)
			sb.WriteString("\n")
		}
	}

	// Bottom border
	sb.WriteString(ColorCyan)
	sb.WriteString(BoxBottomLeft)
	sb.WriteString(strings.Repeat(BoxHorizontal, width-2))
	sb.WriteString(BoxBottomRight)
	sb.WriteString(ColorReset)

	return sb.String()
}

// Header creates a styled section header
func Header(title string, width int) string {
	if width < len(title)+4 {
		width = len(title) + 4
	}

	var sb strings.Builder

	// Top border
	sb.WriteString(ColorCyan)
	sb.WriteString(HeavyTopLeft)
	sb.WriteString(strings.Repeat(HeavyHorizontal, width-2))
	sb.WriteString(HeavyTopRight)
	sb.WriteString(ColorReset)
	sb.WriteString("\n")

	// Title line
	padding := width - 4 - len(title)
	leftPad := padding / 2
	rightPad := padding - leftPad
	sb.WriteString(ColorCyan)
	sb.WriteString(HeavyVertical)
	sb.WriteString(ColorReset)
	sb.WriteString(" ")
	sb.WriteString(strings.Repeat(" ", leftPad))
	sb.WriteString(ColorBold)
	sb.WriteString(ColorCyan)
	sb.WriteString(title)
	sb.WriteString(ColorReset)
	sb.WriteString(strings.Repeat(" ", rightPad))
	sb.WriteString(" ")
	sb.WriteString(ColorCyan)
	sb.WriteString(HeavyVertical)
	sb.WriteString(ColorReset)
	sb.WriteString("\n")

	// Bottom border
	sb.WriteString(ColorCyan)
	sb.WriteString(HeavyBottomLeft)
	sb.WriteString(strings.Repeat(HeavyHorizontal, width-2))
	sb.WriteString(HeavyBottomRight)
	sb.WriteString(ColorReset)

	return sb.String()
}

// Divider creates a horizontal divider line
func Divider(width int) string {
	return fmt.Sprintf("%s%s%s", ColorDim, strings.Repeat("─", width), ColorReset)
}

// SuccessMessage formats a success message
func SuccessMessage(message string) string {
	return fmt.Sprintf("%s%s %s%s", ColorGreen, IconCheck, message, ColorReset)
}

// ErrorMessage formats an error message with a box
func ErrorMessage(message string) string {
	return fmt.Sprintf("%s%s Error: %s%s", ColorRed, IconCross, message, ColorReset)
}

// WarningMessage formats a warning message
func WarningMessage(message string) string {
	return fmt.Sprintf("%s%s %s%s", ColorYellow, IconWarning, message, ColorReset)
}

// InfoMessage formats an info message
func InfoMessage(message string) string {
	return fmt.Sprintf("%s%s %s%s", ColorCyan, IconInfo, message, ColorReset)
}

// Badge creates a small status badge
func Badge(text string, color string) string {
	return fmt.Sprintf("%s[%s]%s", color, text, ColorReset)
}

// StatusBadge creates a status indicator
func StatusBadge(executed bool) string {
	if executed {
		return fmt.Sprintf("%s%s executed%s", ColorGreen, IconCheck, ColorReset)
	}
	return fmt.Sprintf("%s%s copied%s", ColorDim, IconCopy, ColorReset)
}

// KeyHint formats a keyboard shortcut hint
func KeyHint(key, action string) string {
	return fmt.Sprintf("%s[%s%s%s]%s %s", ColorDim, ColorYellow, key, ColorDim, ColorReset, action)
}

// FormatPrompt creates a styled input prompt
func FormatPrompt(prefix string) string {
	return fmt.Sprintf("%s%s%s%s ", ColorBold, ColorCyan, prefix, ColorReset)
}

// Faint makes text dimmed
func Faint(text string) string {
	return fmt.Sprintf("%s%s%s", ColorDim, text, ColorReset)
}

// Bold makes text bold
func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", ColorBold, text, ColorReset)
}

// Cyan colors text cyan
func Cyan(text string) string {
	return fmt.Sprintf("%s%s%s", ColorCyan, text, ColorReset)
}

// Green colors text green
func Green(text string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, text, ColorReset)
}

// Yellow colors text yellow
func Yellow(text string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, text, ColorReset)
}

// Red colors text red
func Red(text string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, text, ColorReset)
}

// ProgressDots returns an animated progress indicator character
func ProgressDots(frame int) string {
	dots := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return dots[frame%len(dots)]
}

// FormatDuration formats a duration for display
func FormatDuration(seconds float64) string {
	if seconds < 1 {
		return fmt.Sprintf("%.0fms", seconds*1000)
	}
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm %ds", minutes, secs)
}

// TruncateString truncates a string and adds ellipsis if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

