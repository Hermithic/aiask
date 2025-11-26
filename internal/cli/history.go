package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/history"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	historyLimit  int
	historySearch string
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View command history",
	Long: `View the history of generated commands.

Examples:
  aiask history              # Show recent history
  aiask history -n 5         # Show last 5 entries
  aiask history --search git # Search history for "git"
  aiask history clear        # Clear all history`,
	Run: runHistory,
}

var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear command history",
	Run:   runHistoryClear,
}

func init() {
	historyCmd.Flags().IntVarP(&historyLimit, "number", "n", 10, "Number of entries to show")
	historyCmd.Flags().StringVar(&historySearch, "search", "", "Search term to filter history")
	historyCmd.AddCommand(historyClearCmd)
}

func runHistory(cmd *cobra.Command, args []string) {
	h, err := history.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load history: %w", err))
		return
	}

	if len(h.Entries) == 0 {
		fmt.Println()
		fmt.Println(ui.InfoMessage("No command history found."))
		fmt.Printf("%sHistory is recorded when you generate commands with aiask.%s\n", ui.ColorDim, ui.ColorReset)
		return
	}

	var entries []history.Entry

	if historySearch != "" {
		entries = h.Search(historySearch)
		if len(entries) == 0 {
			fmt.Println()
			fmt.Printf("%sNo history entries matching '%s%s%s'%s\n", ui.ColorDim, ui.ColorCyan, historySearch, ui.ColorDim, ui.ColorReset)
			return
		}
	} else {
		entries = h.GetRecent(historyLimit)
	}

	fmt.Println()
	if historySearch != "" {
		fmt.Println(ui.Header(fmt.Sprintf("Search: %s", historySearch), 52))
	} else {
		fmt.Println(ui.Header("Command History", 52))
	}
	fmt.Println()

	for i, entry := range entries {
		printHistoryEntry(i+1, entry)
	}

	// Footer summary
	fmt.Println(ui.Divider(52))
	if historySearch != "" {
		fmt.Printf("%sFound %d matching entries%s\n", ui.ColorDim, len(entries), ui.ColorReset)
	} else if len(h.Entries) > historyLimit {
		fmt.Printf("%sShowing %d of %d entries. Use %s-n%s to see more.%s\n",
			ui.ColorDim, historyLimit, len(h.Entries), ui.ColorCyan, ui.ColorDim, ui.ColorReset)
	} else {
		fmt.Printf("%sTotal: %d entries%s\n", ui.ColorDim, len(entries), ui.ColorReset)
	}
	fmt.Println()
}

// printHistoryEntry formats and prints a single history entry
func printHistoryEntry(num int, entry history.Entry) {
	// Time formatting - show relative time for recent, absolute for older
	timeStr := formatRelativeTime(entry.Timestamp)

	// Status badge
	var statusBadge string
	if entry.Executed {
		statusBadge = fmt.Sprintf("%s%s executed%s", ui.ColorGreen, ui.IconCheck, ui.ColorReset)
	} else {
		statusBadge = fmt.Sprintf("%s%s copied%s", ui.ColorDim, ui.IconCopy, ui.ColorReset)
	}

	// Entry header with number and time
	fmt.Printf("%s%s%d%s  %s%s%s  %s\n",
		ui.ColorBold, ui.ColorCyan, num, ui.ColorReset,
		ui.ColorDim, timeStr, ui.ColorReset,
		statusBadge)

	// Prompt (truncated if too long)
	prompt := entry.Prompt
	if len(prompt) > 70 {
		prompt = prompt[:67] + "..."
	}
	fmt.Printf("   %s%s Prompt:%s  %s\n", ui.ColorDim, ui.IconArrow, ui.ColorReset, prompt)

	// Command with syntax highlighting
	highlighter := ui.NewHighlighter()
	command := entry.Command
	// Handle multi-line commands
	commandLines := strings.Split(command, "\n")
	if len(commandLines) == 1 {
		if len(command) > 70 {
			command = command[:67] + "..."
		}
		fmt.Printf("   %s%s Command:%s %s\n", ui.ColorDim, ui.IconTerminal, ui.ColorReset, highlighter.Highlight(command))
	} else {
		fmt.Printf("   %s%s Command:%s\n", ui.ColorDim, ui.IconTerminal, ui.ColorReset)
		for _, line := range commandLines {
			if strings.TrimSpace(line) != "" {
				if len(line) > 70 {
					line = line[:67] + "..."
				}
				fmt.Printf("      %s\n", highlighter.Highlight(line))
			}
		}
	}

	// Shell indicator
	fmt.Printf("   %s%s Shell:%s   %s\n", ui.ColorDim, ui.IconDot, ui.ColorReset, entry.Shell)

	fmt.Println()
}

// formatRelativeTime formats a timestamp as relative time for recent entries
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 02, 2006")
	}
}

func runHistoryClear(cmd *cobra.Command, args []string) {
	h, err := history.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load history: %w", err))
		return
	}

	count := len(h.Entries)
	if count == 0 {
		fmt.Println(ui.InfoMessage("No history to clear."))
		return
	}

	// Confirmation prompt
	fmt.Println()
	fmt.Printf("%s%s Warning:%s This will delete %d history entries.\n", ui.ColorYellow, ui.IconWarning, ui.ColorReset, count)

	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to clear all history",
		IsConfirm: true,
	}

	_, err = confirmPrompt.Run()
	if err != nil {
		fmt.Println("History not cleared.")
		return
	}

	h.Clear()

	if err := h.Save(); err != nil {
		ui.ShowError(fmt.Errorf("failed to save history: %w", err))
		return
	}

	fmt.Println()
	fmt.Println(ui.SuccessMessage(fmt.Sprintf("Cleared %d history entries", count)))
}
