package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/history"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/ui"
)

// REPL represents an interactive REPL session
type REPL struct {
	cfg          *config.Config
	provider     llm.Provider
	shellInfo    shell.ShellInfo
	reader       *bufio.Reader
	history      []string
	historyIdx   int
	commandCount int
	startTime    time.Time
}

// New creates a new REPL instance
func New(cfg *config.Config, provider llm.Provider, shellInfo shell.ShellInfo) *REPL {
	return &REPL{
		cfg:          cfg,
		provider:     provider,
		shellInfo:    shellInfo,
		reader:       bufio.NewReader(os.Stdin),
		history:      []string{},
		historyIdx:   -1,
		commandCount: 0,
		startTime:    time.Now(),
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
	fmt.Println(ui.Header("AIask Interactive Mode", 44))
	fmt.Println()

	// Status bar
	fmt.Printf("%s%s Provider:%s %s%s%s  ", ui.ColorDim, ui.IconDot, ui.ColorReset, ui.ColorCyan, r.cfg.Provider, ui.ColorReset)
	fmt.Printf("%s%s Model:%s %s%s%s  ", ui.ColorDim, ui.IconDot, ui.ColorReset, ui.ColorCyan, r.cfg.Model, ui.ColorReset)
	fmt.Printf("%s%s Shell:%s %s%s%s\n", ui.ColorDim, ui.IconDot, ui.ColorReset, ui.ColorCyan, shell.GetShellName(r.shellInfo.Shell), ui.ColorReset)
	fmt.Println()

	fmt.Printf("%sQuick Commands:%s\n", ui.ColorBold, ui.ColorReset)
	r.printCommand("/help", "Show all commands")
	r.printCommand("/history", "View session history")
	r.printCommand("/clear", "Clear screen")
	r.printCommand("/exit", "Exit REPL")
	fmt.Println()
	fmt.Println(ui.Divider(44))
	fmt.Printf("%sType your request in natural language:%s\n", ui.ColorDim, ui.ColorReset)
	fmt.Println()
}

// printCommand formats a command hint
func (r *REPL) printCommand(cmd, desc string) {
	fmt.Printf("  %s%-12s%s %s%s%s\n", ui.ColorCyan, cmd, ui.ColorReset, ui.ColorDim, desc, ui.ColorReset)
}

// readPrompt reads a prompt from the user with a styled prompt
func (r *REPL) readPrompt() string {
	// Show command count in prompt
	promptNum := len(r.history) + 1
	fmt.Printf("%s[%d]%s %s%s%s%s ", ui.ColorDim, promptNum, ui.ColorReset, ui.ColorBold, ui.ColorCyan, ui.IconTerminal, ui.ColorReset)

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
		r.showGoodbye()
		return false

	case "config":
		r.showConfig()
		return true

	case "stats":
		r.showStats()
		return true

	default:
		fmt.Println(ui.WarningMessage(fmt.Sprintf("Unknown command: %s. Type /help for available commands.", parts[0])))
		return true
	}
}

// handlePrompt processes a user prompt
func (r *REPL) handlePrompt(prompt string) {
	// Generate command with spinner
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.GetTimeout())
	defer cancel()

	stopSpinner := ui.ShowSpinner(fmt.Sprintf("Generating with %s", r.cfg.Model))

	startTime := time.Now()
	command, err := r.provider.GenerateCommand(ctx, prompt, r.shellInfo)
	elapsed := time.Since(startTime)

	stopSpinner()

	if err != nil {
		ui.ShowError(fmt.Errorf("failed to generate command: %w", err))
		fmt.Println()
		return
	}

	// Show generation time
	fmt.Printf("%sGenerated in %s%s\n", ui.ColorDim, ui.FormatDuration(elapsed.Seconds()), ui.ColorReset)

	// Clean up the command
	command = llm.CleanCommand(command)

	// Display the command
	ui.DisplayCommand(command)

	// Get user action
	action := ui.PromptActionForCommand(command)

	switch action {
	case ui.ActionExecute:
		r.commandCount++
		execErr, _ := ui.ExecuteCommandWithErrorRecovery(command, r.shellInfo)
		if histErr := history.AddEntry(prompt, command, string(r.shellInfo.Shell), execErr == nil); histErr != nil {
			fmt.Printf("%s[REPL] Failed to record history: %s%s\n", ui.ColorDim, histErr, ui.ColorReset)
		}

	case ui.ActionCopy:
		err := ui.CopyToClipboard(command)
		if histErr := history.AddEntry(prompt, command, string(r.shellInfo.Shell), false); histErr != nil {
			fmt.Printf("%s[REPL] Failed to record history: %s%s\n", ui.ColorDim, histErr, ui.ColorReset)
		}
		if err != nil {
			ui.ShowError(err)
		}

	case ui.ActionEdit:
		edited := ui.PromptEdit(command)

		ui.DisplayCommand(edited)
		editAction := ui.PromptActionForCommand(edited)
		if editAction == ui.ActionExecute {
			r.commandCount++
			execErr, _ := ui.ExecuteCommandWithErrorRecovery(edited, r.shellInfo)
			if histErr := history.AddEntry(prompt, edited, string(r.shellInfo.Shell), execErr == nil); histErr != nil {
				fmt.Printf("%s[REPL] Failed to record history: %s%s\n", ui.ColorDim, histErr, ui.ColorReset)
			}
		} else if editAction == ui.ActionCopy {
			err := ui.CopyToClipboard(edited)
			if histErr := history.AddEntry(prompt, edited, string(r.shellInfo.Shell), false); histErr != nil {
				fmt.Printf("%s[REPL] Failed to record history: %s%s\n", ui.ColorDim, histErr, ui.ColorReset)
			}
			if err != nil {
				ui.ShowError(err)
			}
		}

	case ui.ActionReprompt:
		// Prompt for a new request and process it
		newPrompt := ui.PromptReprompt()
		if newPrompt != "" {
			r.history = append(r.history, newPrompt)
			r.historyIdx = len(r.history)
			r.handlePrompt(newPrompt)
			return
		}

	case ui.ActionQuit:
		// Just continue to next prompt in REPL mode
	}

	fmt.Println()
}

// showHelp displays help information
func (r *REPL) showHelp() {
	fmt.Println()
	fmt.Println(ui.Header("Help", 44))
	fmt.Println()

	fmt.Printf("%sNavigation Commands:%s\n", ui.ColorBold, ui.ColorReset)
	r.printCommand("/help, /?", "Show this help message")
	r.printCommand("/history", "Show session history")
	r.printCommand("/clear", "Clear the screen")
	r.printCommand("/config", "Show current configuration")
	r.printCommand("/stats", "Show session statistics")
	r.printCommand("/exit, /q", "Exit interactive mode")
	fmt.Println()

	fmt.Printf("%sAction Keys (in menu):%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  %s↑/↓%s          Navigate menu options\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("  %sEnter%s        Select option\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("  %sType%s         Filter/search options\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()

	fmt.Printf("%sUsage:%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Printf("  Just type your request in natural language and press Enter.\n")
	fmt.Printf("  %sExample:%s list all files larger than 100MB\n", ui.ColorDim, ui.ColorReset)
	fmt.Println()
}

// showHistory displays the session history
func (r *REPL) showHistory() {
	if len(r.history) == 0 {
		fmt.Println(ui.InfoMessage("No history in this session yet."))
		return
	}

	fmt.Println()
	fmt.Println(ui.Header("Session History", 44))
	fmt.Println()

	for i, h := range r.history {
		fmt.Printf("  %s%d.%s %s\n", ui.ColorDim, i+1, ui.ColorReset, h)
	}
	fmt.Println()
}

// showConfig displays the current configuration
func (r *REPL) showConfig() {
	fmt.Println()
	fmt.Println(ui.Header("Configuration", 44))
	fmt.Println()

	fmt.Printf("  %sProvider:%s    %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, r.cfg.Provider, ui.ColorReset)
	fmt.Printf("  %sModel:%s       %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, r.cfg.Model, ui.ColorReset)
	fmt.Printf("  %sShell:%s       %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, shell.GetShellName(r.shellInfo.Shell), ui.ColorReset)
	fmt.Printf("  %sOS:%s          %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, shell.GetOSName(), ui.ColorReset)
	fmt.Printf("  %sTimeout:%s     %s%v%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, r.cfg.GetTimeout(), ui.ColorReset)
	fmt.Println()
}

// showStats displays session statistics
func (r *REPL) showStats() {
	sessionDuration := time.Since(r.startTime)

	fmt.Println()
	fmt.Println(ui.Header("Session Statistics", 44))
	fmt.Println()

	fmt.Printf("  %sSession Duration:%s  %s%s%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, ui.FormatDuration(sessionDuration.Seconds()), ui.ColorReset)
	fmt.Printf("  %sPrompts Entered:%s   %s%d%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, len(r.history), ui.ColorReset)
	fmt.Printf("  %sCommands Executed:%s %s%d%s\n", ui.ColorDim, ui.ColorReset, ui.ColorCyan, r.commandCount, ui.ColorReset)
	fmt.Println()
}

// showGoodbye displays a goodbye message with session summary
func (r *REPL) showGoodbye() {
	sessionDuration := time.Since(r.startTime)

	fmt.Println()
	fmt.Println(ui.Divider(44))
	fmt.Printf("%s%s Session Summary%s\n", ui.ColorBold, ui.IconStar, ui.ColorReset)
	fmt.Printf("  Duration: %s\n", ui.FormatDuration(sessionDuration.Seconds()))
	fmt.Printf("  Prompts: %d  |  Commands executed: %d\n", len(r.history), r.commandCount)
	fmt.Println(ui.Divider(44))
	fmt.Printf("%sGoodbye! %s%s\n", ui.ColorDim, ui.IconTerminal, ui.ColorReset)
}

// clearScreen clears the terminal screen
func (r *REPL) clearScreen() {
	fmt.Print("\033[H\033[2J")
	r.printWelcome()
}
