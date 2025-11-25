package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Version is set during build
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "aiask [prompt]",
	Short: "AI-powered command line assistant",
	Long: `AIask is an AI-powered command line assistant that converts natural 
language into shell commands for PowerShell, CMD, Bash, and Zsh.

Example:
  aiask "list all files larger than 100MB"
  aiask "find and delete all .tmp files"
  aiask "compress the current directory into a zip file"`,
	Args: cobra.ArbitraryArgs,
	Run:  runMain,
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func runMain(cmd *cobra.Command, args []string) {
	// Check if prompt is provided
	if len(args) == 0 {
		fmt.Println("Usage: aiask \"your request here\"")
		fmt.Println()
		fmt.Println("Example: aiask \"list all files in the current directory\"")
		fmt.Println()
		fmt.Println("Run 'aiask config' to set up your LLM provider.")
		os.Exit(1)
	}

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

	// Join args into a single prompt
	prompt := strings.Join(args, " ")

	// Run the main interaction loop
	runInteractionLoop(provider, prompt, shellInfo)
}

func runInteractionLoop(provider llm.Provider, prompt string, shellInfo shell.ShellInfo) {
	for {
		// Generate command
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		
		fmt.Printf("\n%sGenerating command...%s\n", ui.ColorDim, ui.ColorReset)
		
		command, err := provider.GenerateCommand(ctx, prompt, shellInfo)
		cancel()

		if err != nil {
			ui.ShowError(fmt.Errorf("failed to generate command: %w", err))
			return
		}

		// Clean up the command (remove any markdown code blocks if present)
		command = cleanCommand(command)

		// Display the command
		ui.DisplayCommand(command)

		// Get user action
		action := ui.PromptAction()

		switch action {
		case ui.ActionExecute:
			err := ui.ExecuteCommand(command, shellInfo)
			if err != nil {
				ui.ShowError(err)
			}
			return

		case ui.ActionCopy:
			err := ui.CopyToClipboard(command)
			if err != nil {
				ui.ShowError(err)
			}
			return

		case ui.ActionEdit:
			editedCommand := ui.PromptEdit(command)
			ui.DisplayCommand(editedCommand)
			
			// Ask what to do with edited command
			editAction := ui.PromptAction()
			switch editAction {
			case ui.ActionExecute:
				err := ui.ExecuteCommand(editedCommand, shellInfo)
				if err != nil {
					ui.ShowError(err)
				}
			case ui.ActionCopy:
				err := ui.CopyToClipboard(editedCommand)
				if err != nil {
					ui.ShowError(err)
				}
			}
			return

		case ui.ActionReprompt:
			newPrompt := ui.PromptReprompt()
			if newPrompt == "" {
				fmt.Println("No new prompt provided. Exiting.")
				return
			}
			prompt = newPrompt
			continue

		case ui.ActionQuit:
			fmt.Println("Goodbye!")
			return
		}
	}
}

// cleanCommand removes markdown code blocks and extra whitespace from the command
func cleanCommand(command string) string {
	// Remove markdown code blocks
	command = strings.TrimPrefix(command, "```bash\n")
	command = strings.TrimPrefix(command, "```powershell\n")
	command = strings.TrimPrefix(command, "```cmd\n")
	command = strings.TrimPrefix(command, "```shell\n")
	command = strings.TrimPrefix(command, "```sh\n")
	command = strings.TrimPrefix(command, "```\n")
	command = strings.TrimSuffix(command, "\n```")
	command = strings.TrimSuffix(command, "```")
	
	// Trim whitespace
	command = strings.TrimSpace(command)
	
	return command
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aiask version %s\n", Version)
	},
}

