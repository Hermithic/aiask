package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/history"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/Hermithic/aiask/internal/update"
	"github.com/spf13/cobra"
)

var (
	// Version is set during build
	Version = "dev"

	// Command line flags
	verbose    bool
	jsonOutput bool
	useStdin   bool
	streaming  bool
)

// JSONOutput represents the JSON output format
type JSONOutput struct {
	Command  string `json:"command"`
	Shell    string `json:"shell"`
	OS       string `json:"os"`
	Prompt   string `json:"prompt"`
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
}

var rootCmd = &cobra.Command{
	Use:   "aiask [prompt]",
	Short: "AI-powered command line assistant",
	Long: `AIask is an AI-powered command line assistant that converts natural 
language into shell commands for PowerShell, CMD, Bash, and Zsh.

Example:
  aiask "list all files larger than 100MB"
  aiask "find and delete all .tmp files"
  aiask "compress the current directory into a zip file"

Environment Variables:
  AIASK_PROVIDER    - LLM provider (grok, openai, anthropic, gemini, ollama)
  AIASK_API_KEY     - API key for the provider
  AIASK_MODEL       - Model name to use
  AIASK_OLLAMA_URL  - Ollama server URL (default: http://localhost:11434)
  AIASK_TIMEOUT     - Request timeout in seconds (default: 60)`,
	Args: cobra.ArbitraryArgs,
	Run:  runMain,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output including debug information")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output result as JSON (non-interactive)")
	rootCmd.PersistentFlags().BoolVar(&useStdin, "stdin", false, "Read additional context from stdin")
	rootCmd.PersistentFlags().BoolVarP(&streaming, "stream", "s", false, "Stream the response as it generates")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(explainCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(completionCmd)
}

// Execute runs the root command
func Execute() error {
	// Check for updates in the background (non-blocking)
	checkForUpdates()
	return rootCmd.Execute()
}

// checkForUpdates checks for updates asynchronously and displays a message if available
func checkForUpdates() {
	// Load config to check if update checks are enabled
	cfg, err := config.Load()
	if err != nil || !cfg.CheckUpdates {
		return
	}

	// Check for updates in the background
	update.CheckForUpdatesAsync(Version, func(result *update.CheckResult) {
		if result.UpdateAvailable {
			// Only print if we're not in JSON mode
			if !jsonOutput {
				fmt.Println(update.FormatUpdateMessage(result))
			}
		}
	})
}

func runMain(cmd *cobra.Command, args []string) {
	// Check if prompt is provided
	if len(args) == 0 && !useStdin {
		fmt.Println("Usage: aiask \"your request here\"")
		fmt.Println()
		fmt.Println("Example: aiask \"list all files in the current directory\"")
		fmt.Println()
		fmt.Println("You can also pipe input:")
		fmt.Println("  cat error.log | aiask --stdin \"what's wrong here?\"")
		fmt.Println()
		fmt.Println("Run 'aiask config' to set up your LLM provider.")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		if !jsonOutput {
			fmt.Printf("Configuration error: %s\n", err)
			fmt.Println("Run 'aiask config' to set up your configuration.")
		} else {
			outputJSON(JSONOutput{}, err)
		}
		os.Exit(1)
	}

	// Detect shell
	shellInfo := shell.Detect()

	// Verbose output
	if verbose {
		printVerboseInfo(cfg, shellInfo)
	}

	// Create LLM provider
	provider, err := llm.NewProvider(cfg)
	if err != nil {
		if !jsonOutput {
			ui.ShowError(fmt.Errorf("failed to create LLM provider: %w", err))
		} else {
			outputJSON(JSONOutput{}, err)
		}
		os.Exit(1)
	}

	// Join args into a single prompt
	prompt := strings.Join(args, " ")

	// Read from stdin if requested or if stdin has data
	stdinContent := readStdin()
	if stdinContent != "" {
		if prompt != "" {
			prompt = prompt + "\n\nContext:\n" + stdinContent
		} else {
			prompt = "Analyze this output and suggest a solution:\n\n" + stdinContent
		}
		if verbose {
			fmt.Printf("%s[DEBUG] Read %d bytes from stdin%s\n", ui.ColorDim, len(stdinContent), ui.ColorReset)
		}
	}

	if verbose {
		fmt.Printf("%s[DEBUG] Prompt: %s%s\n", ui.ColorDim, truncateString(prompt, 200), ui.ColorReset)
	}

	// Run the main interaction loop
	runInteractionLoop(provider, prompt, shellInfo, cfg)
}

// readStdin reads from stdin if it's a pipe (not a terminal)
func readStdin() string {
	// Check if stdin has data (is a pipe, not a terminal)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}

	// If stdin is a character device (terminal), only read if --stdin flag is set
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		if !useStdin {
			return ""
		}
	}

	// Read stdin with a limit to avoid memory issues
	const maxStdinBytes = 50000 // ~50KB limit
	reader := bufio.NewReader(os.Stdin)
	var content strings.Builder
	bytesRead := 0

	for bytesRead < maxStdinBytes {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				content.WriteString(line)
			}
			break
		}
		content.WriteString(line)
		bytesRead += len(line)
	}

	result := strings.TrimSpace(content.String())

	// Truncate if too long
	if len(result) > maxStdinBytes {
		result = result[:maxStdinBytes] + "\n... (truncated)"
	}

	return result
}

// truncateString truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// printVerboseInfo prints debug information when verbose mode is enabled
func printVerboseInfo(cfg *config.Config, shellInfo shell.ShellInfo) {
	fmt.Printf("%s[DEBUG] Shell: %s%s\n", ui.ColorDim, shell.GetShellName(shellInfo.Shell), ui.ColorReset)
	fmt.Printf("%s[DEBUG] OS: %s%s\n", ui.ColorDim, shell.GetOSName(), ui.ColorReset)
	fmt.Printf("%s[DEBUG] Provider: %s%s\n", ui.ColorDim, cfg.Provider, ui.ColorReset)
	fmt.Printf("%s[DEBUG] Model: %s%s\n", ui.ColorDim, cfg.Model, ui.ColorReset)
	fmt.Printf("%s[DEBUG] Timeout: %v%s\n", ui.ColorDim, cfg.GetTimeout(), ui.ColorReset)
	if cfg.SystemPromptSuffix != "" {
		fmt.Printf("%s[DEBUG] System prompt suffix: %s%s\n", ui.ColorDim, cfg.SystemPromptSuffix, ui.ColorReset)
	}
}

// outputJSON outputs the result as JSON
func outputJSON(output JSONOutput, err error) {
	type jsonError struct {
		Error string `json:"error"`
	}

	if err != nil {
		data, _ := json.Marshal(jsonError{Error: err.Error()})
		fmt.Println(string(data))
		return
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

func runInteractionLoop(provider llm.Provider, prompt string, shellInfo shell.ShellInfo, cfg *config.Config) {
	for {
		// Generate command with configurable timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GetTimeout())

		startTime := time.Now()
		if !jsonOutput {
			fmt.Printf("\n%sGenerating command...%s\n", ui.ColorDim, ui.ColorReset)
		}

		command, err := provider.GenerateCommand(ctx, prompt, shellInfo)
		cancel()

		if verbose {
			fmt.Printf("%s[DEBUG] Response time: %v%s\n", ui.ColorDim, time.Since(startTime), ui.ColorReset)
		}

		if err != nil {
			if jsonOutput {
				outputJSON(JSONOutput{Prompt: prompt}, fmt.Errorf("failed to generate command: %w", err))
			} else {
				ui.ShowError(fmt.Errorf("failed to generate command: %w", err))
			}
			return
		}

		// Clean up the command (remove any markdown code blocks if present)
		command = cleanCommand(command)

		// JSON output mode - non-interactive
		if jsonOutput {
			outputJSON(JSONOutput{
				Command:  command,
				Shell:    string(shellInfo.Shell),
				OS:       shellInfo.OS,
				Prompt:   prompt,
				Provider: string(cfg.Provider),
				Model:    cfg.Model,
			}, nil)
			// Record in history (not executed)
			_ = history.AddEntry(prompt, command, string(shellInfo.Shell), false)
			return
		}

		// Display the command
		ui.DisplayCommand(command)

		// Get user action (with safety checks for dangerous commands)
		action := ui.PromptActionForCommand(command)

		switch action {
		case ui.ActionExecute:
			execErr, wantsRecovery := ui.ExecuteCommandWithErrorRecovery(command, shellInfo)
			// Record in history (executed)
			_ = history.AddEntry(prompt, command, string(shellInfo.Shell), execErr == nil)

			// If command failed and user wants help, generate a diagnostic prompt
			if wantsRecovery && execErr != nil {
				recoveryPrompt := fmt.Sprintf("The command '%s' failed with error: %s. How can I fix this?", command, execErr.Error())
				prompt = recoveryPrompt
				continue
			}
			return

		case ui.ActionCopy:
			err := ui.CopyToClipboard(command)
			// Record in history (not executed, just copied)
			_ = history.AddEntry(prompt, command, string(shellInfo.Shell), false)
			if err != nil {
				ui.ShowError(err)
			}
			return

		case ui.ActionEdit:
			editedCommand := ui.PromptEdit(command)
			ui.DisplayCommand(editedCommand)

			// Ask what to do with edited command (with safety checks)
			editAction := ui.PromptActionForCommand(editedCommand)
			switch editAction {
			case ui.ActionExecute:
				execErr, wantsRecovery := ui.ExecuteCommandWithErrorRecovery(editedCommand, shellInfo)
				// Record edited command in history
				_ = history.AddEntry(prompt, editedCommand, string(shellInfo.Shell), execErr == nil)

				if wantsRecovery && execErr != nil {
					recoveryPrompt := fmt.Sprintf("The command '%s' failed with error: %s. How can I fix this?", editedCommand, execErr.Error())
					prompt = recoveryPrompt
					continue
				}
			case ui.ActionCopy:
				err := ui.CopyToClipboard(editedCommand)
				// Record edited command in history (not executed)
				_ = history.AddEntry(prompt, editedCommand, string(shellInfo.Shell), false)
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

