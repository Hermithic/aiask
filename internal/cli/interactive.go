package cli

import (
	"fmt"
	"os"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/repl"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i", "repl"},
	Short:   "Start interactive REPL mode",
	Long: `Start an interactive REPL (Read-Eval-Print Loop) mode for continuous
interaction with AIask without restarting.

Commands available in REPL:
  /help     - Show available commands
  /history  - Show session history
  /clear    - Clear the screen
  /config   - Show current configuration
  /exit     - Exit interactive mode`,
	Run: runInteractive,
}

func init() {
	// Command will be added in root.go
}

func runInteractive(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Configuration error: %s\n", err)
		fmt.Println("Run 'aiask config' to set up your configuration.")
		os.Exit(1)
	}

	// Detect shell
	shellInfo := shell.Detect()

	// Create LLM provider
	provider, err := llm.NewProvider(cfg)
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to create LLM provider: %w", err))
		os.Exit(1)
	}

	// Start REPL
	r := repl.New(cfg, provider, shellInfo)
	r.Run()
}

