package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain <command>",
	Short: "Explain what a command does",
	Long: `Explain what a shell command does in plain English.

Examples:
  aiask explain "tar -xzvf archive.tar.gz"
  aiask explain "find . -name '*.log' -mtime +7 -delete"
  aiask explain "git rebase -i HEAD~3"`,
	Args: cobra.MinimumNArgs(1),
	Run:  runExplain,
}

func init() {
	// Already added in root.go init
}

func runExplain(cmd *cobra.Command, args []string) {
	command := strings.Join(args, " ")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("configuration error: %w", err))
		fmt.Println("Run 'aiask config' to set up your configuration.")
		os.Exit(1)
	}

	// Create LLM provider
	provider, err := llm.NewProvider(cfg)
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to create LLM provider: %w", err))
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("%s[DEBUG] Command to explain: %s%s\n", ui.ColorDim, command, ui.ColorReset)
	}

	fmt.Printf("\n%sAnalyzing command...%s\n\n", ui.ColorDim, ui.ColorReset)

	// Generate explanation
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetTimeout())
	defer cancel()

	explanation, err := provider.ExplainCommand(ctx, command)
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to explain command: %w", err))
		return
	}

	// Display the command
	fmt.Printf("%s%sCommand:%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset)
	fmt.Printf("  %s%s%s\n\n", ui.ColorGreen, command, ui.ColorReset)

	// Display the explanation
	fmt.Printf("%s%sExplanation:%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset)
	
	// Format the explanation with proper indentation
	lines := strings.Split(explanation, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("  %s\n", line)
		} else {
			fmt.Println()
		}
	}
	fmt.Println()
}

