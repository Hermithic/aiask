package cli

import (
	"fmt"

	"github.com/Hermithic/aiask/internal/history"
	"github.com/Hermithic/aiask/internal/ui"
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
		fmt.Println("No command history found.")
		fmt.Println("History is recorded when you generate commands with aiask.")
		return
	}

	var entries []history.Entry

	if historySearch != "" {
		entries = h.Search(historySearch)
		if len(entries) == 0 {
			fmt.Printf("No history entries matching '%s'\n", historySearch)
			return
		}
	} else {
		entries = h.GetRecent(historyLimit)
	}

	fmt.Printf("%s╔══════════════════════════════════════════╗%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s║           Command History                ║%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════╝%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()

	for i, entry := range entries {
		// Time ago
		timeStr := entry.Timestamp.Format("2006-01-02 15:04")

		// Status indicator
		status := "📋"
		if entry.Executed {
			status = "✓"
		}

		fmt.Printf("%s%d. [%s] %s%s\n", ui.ColorBold, i+1, timeStr, status, ui.ColorReset)
		fmt.Printf("   %sPrompt:%s %s\n", ui.ColorDim, ui.ColorReset, entry.Prompt)
		fmt.Printf("   %sCommand:%s %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorGreen, entry.Command, ui.ColorReset)
		fmt.Printf("   %sShell:%s %s\n", ui.ColorDim, ui.ColorReset, entry.Shell)
		fmt.Println()
	}

	if historySearch == "" && len(h.Entries) > historyLimit {
		fmt.Printf("%sShowing %d of %d entries. Use -n to see more.%s\n",
			ui.ColorDim, historyLimit, len(h.Entries), ui.ColorReset)
	}
}

func runHistoryClear(cmd *cobra.Command, args []string) {
	h, err := history.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load history: %w", err))
		return
	}

	count := len(h.Entries)
	h.Clear()

	if err := h.Save(); err != nil {
		ui.ShowError(fmt.Errorf("failed to save history: %w", err))
		return
	}

	fmt.Printf("%s✓ Cleared %d history entries%s\n", ui.ColorGreen, count, ui.ColorReset)
}

