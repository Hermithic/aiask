package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/history"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/ui"
)

// REPL represents an interactive REPL session
type REPL struct {
	cfg        *config.Config
	provider   llm.Provider
	shellInfo  shell.ShellInfo
	reader     *bufio.Reader
	history    []string
	historyIdx int
}

// New creates a new REPL instance
func New(cfg *config.Config, provider llm.Provider, shellInfo shell.ShellInfo) *REPL {
	return &REPL{
		cfg:        cfg,
		provider:   provider,
		shellInfo:  shellInfo,
		reader:     bufio.NewReader(os.Stdin),
		history:    []string{},
		historyIdx: -1,
	}
}

// Run starts the REPL loop
func (r *REPL) Run() {
	r.printWelcome()

	for {
		prompt := r.readPrompt()
		if prompt == "" {
			continue
		}

		// Handle special commands
		if strings.HasPrefix(prompt, "/") {
			if r.handleCommand(prompt) {
				continue
			}
			return // Exit on /exit or /quit
		}

		// Add to history
		r.history = append(r.history, prompt)
		r.historyIdx = len(r.history)

		// Generate and handle the command
		r.handlePrompt(prompt)
	}
}

// printWelcome prints the welcome message
func (r *REPL) printWelcome() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════╗%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s║       AIask Interactive Mode             ║%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════╝%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()
	fmt.Printf("%sCommands:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  /help     - Show available commands\n")
	fmt.Printf("  /history  - Show session history\n")
	fmt.Printf("  /clear    - Clear the screen\n")
	fmt.Printf("  /exit     - Exit interactive mode\n")
	fmt.Println()
	fmt.Printf("%sEnter your requests in natural language:%s\n", ui.ColorDim, ui.ColorReset)
	fmt.Println()
}

// readPrompt reads a prompt from the user
func (r *REPL) readPrompt() string {
	fmt.Printf("%s%s❯%s ", ui.ColorBold, ui.ColorCyan, ui.ColorReset)

	input, err := r.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n%sGoodbye!%s\n", ui.ColorDim, ui.ColorReset)
		os.Exit(0)
	}

	return strings.TrimSpace(input)
}

// handleCommand handles special REPL commands
func (r *REPL) handleCommand(cmd string) bool {
	cmd = strings.ToLower(strings.TrimPrefix(cmd, "/"))
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return true
	}

	switch parts[0] {
	case "help", "h", "?":
		r.showHelp()
		return true

	case "history", "hist":
		r.showHistory()
		return true

	case "clear", "cls":
		r.clearScreen()
		return true

	case "exit", "quit", "q":
		fmt.Printf("%sGoodbye!%s\n", ui.ColorDim, ui.ColorReset)
		return false

	case "config":
		r.showConfig()
		return true

	default:
		fmt.Printf("%sUnknown command: %s. Type /help for available commands.%s\n",
			ui.ColorYellow, parts[0], ui.ColorReset)
		return true
	}
}

// handlePrompt processes a user prompt
func (r *REPL) handlePrompt(prompt string) {
	// Generate command
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.GetTimeout())
	defer cancel()

	fmt.Printf("\n%sGenerating command...%s\n", ui.ColorDim, ui.ColorReset)

	command, err := r.provider.GenerateCommand(ctx, prompt, r.shellInfo)
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to generate command: %w", err))
		fmt.Println()
		return
	}

	// Clean up the command
	command = cleanCommand(command)

	// Display the command
	ui.DisplayCommand(command)

	// Get user action
	action := ui.PromptActionForCommand(command)

	switch action {
	case ui.ActionExecute:
		execErr, _ := ui.ExecuteCommandWithErrorRecovery(command, r.shellInfo)
		_ = history.AddEntry(prompt, command, string(r.shellInfo.Shell), execErr == nil)

	case ui.ActionCopy:
		err := ui.CopyToClipboard(command)
		_ = history.AddEntry(prompt, command, string(r.shellInfo.Shell), false)
		if err != nil {
			ui.ShowError(err)
		}

	case ui.ActionEdit:
		fmt.Printf("%sEdit the command:%s\n", ui.ColorBold, ui.ColorReset)
		fmt.Printf("%s%s%s\n", ui.ColorDim, command, ui.ColorReset)
		fmt.Print("> ")

		edited, _ := r.reader.ReadString('\n')
		edited = strings.TrimSpace(edited)
		if edited == "" {
			edited = command
		}

		ui.DisplayCommand(edited)
		editAction := ui.PromptActionForCommand(edited)
		if editAction == ui.ActionExecute {
			execErr, _ := ui.ExecuteCommandWithErrorRecovery(edited, r.shellInfo)
			_ = history.AddEntry(prompt, edited, string(r.shellInfo.Shell), execErr == nil)
		} else if editAction == ui.ActionCopy {
			err := ui.CopyToClipboard(edited)
			_ = history.AddEntry(prompt, edited, string(r.shellInfo.Shell), false)
			if err != nil {
				ui.ShowError(err)
			}
		}

	case ui.ActionQuit:
		// Just continue to next prompt in REPL mode
	}

	fmt.Println()
}

// showHelp displays help information
func (r *REPL) showHelp() {
	fmt.Println()
	fmt.Printf("%sAvailable Commands:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Println("  /help, /h, /?   - Show this help message")
	fmt.Println("  /history, /hist - Show session history")
	fmt.Println("  /clear, /cls    - Clear the screen")
	fmt.Println("  /config         - Show current configuration")
	fmt.Println("  /exit, /quit, /q - Exit interactive mode")
	fmt.Println()
	fmt.Printf("%sUsage:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Println("  Just type your request in natural language and press Enter.")
	fmt.Println("  Example: list all files larger than 100MB")
	fmt.Println()
}

// showHistory displays the session history
func (r *REPL) showHistory() {
	if len(r.history) == 0 {
		fmt.Printf("%sNo history yet.%s\n", ui.ColorDim, ui.ColorReset)
		return
	}

	fmt.Println()
	fmt.Printf("%sSession History:%s\n", ui.ColorBold, ui.ColorReset)
	for i, h := range r.history {
		fmt.Printf("  %d. %s\n", i+1, h)
	}
	fmt.Println()
}

// showConfig displays the current configuration
func (r *REPL) showConfig() {
	fmt.Println()
	fmt.Printf("%sCurrent Configuration:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  Provider: %s\n", r.cfg.Provider)
	fmt.Printf("  Model: %s\n", r.cfg.Model)
	fmt.Printf("  Shell: %s\n", shell.GetShellName(r.shellInfo.Shell))
	fmt.Printf("  OS: %s\n", shell.GetOSName())
	fmt.Printf("  Timeout: %v\n", r.cfg.GetTimeout())
	fmt.Println()
}

// clearScreen clears the terminal screen
func (r *REPL) clearScreen() {
	fmt.Print("\033[H\033[2J")
	r.printWelcome()
}

// cleanCommand removes markdown code blocks from the command
func cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash\n")
	command = strings.TrimPrefix(command, "```powershell\n")
	command = strings.TrimPrefix(command, "```cmd\n")
	command = strings.TrimPrefix(command, "```shell\n")
	command = strings.TrimPrefix(command, "```sh\n")
	command = strings.TrimPrefix(command, "```\n")
	command = strings.TrimSuffix(command, "\n```")
	command = strings.TrimSuffix(command, "```")
	return strings.TrimSpace(command)
}

