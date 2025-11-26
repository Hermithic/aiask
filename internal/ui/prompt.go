package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Hermithic/aiask/internal/safety"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/undo"
	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
)

// Action represents the user's chosen action
type Action int

const (
	ActionExecute Action = iota
	ActionCopy
	ActionEdit
	ActionReprompt
	ActionQuit
)

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorRed    = "\033[31m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// DisplayCommand shows the suggested command to the user with syntax highlighting
func DisplayCommand(command string) {
	fmt.Println()
	fmt.Printf("%s%s%s Suggested command:%s\n", ColorBold, ColorCyan, IconTerminal, ColorReset)
	fmt.Println(Divider(44))

	// Display the command with syntax highlighting
	highlighter := NewHighlighter()
	lines := strings.Split(command, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			highlighted := highlighter.Highlight(line)
			fmt.Printf("  %s\n", highlighted)
		}
	}
	fmt.Println()

	// Check for dangerous commands and display warning
	warning := safety.GetWarningMessage(command)
	if warning != "" {
		fmt.Println(warning)
		fmt.Println()
	}
}

// actionItem represents a selectable action in the menu
type actionItem struct {
	Label  string
	Action Action
	Icon   string
	Key    string
}

// PromptAction prompts the user for an action using an interactive menu
func PromptAction() Action {
	items := []actionItem{
		{Label: "Execute", Action: ActionExecute, Icon: IconRocket, Key: "e"},
		{Label: "Copy to clipboard", Action: ActionCopy, Icon: IconCopy, Key: "c"},
		{Label: "Edit command", Action: ActionEdit, Icon: IconEdit, Key: "d"},
		{Label: "New prompt", Action: ActionReprompt, Icon: IconRefresh, Key: "r"},
		{Label: "Quit", Action: ActionQuit, Icon: IconExit, Key: "q"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   fmt.Sprintf("%s%s {{ .Icon }} {{ .Label | cyan | bold }}%s %s({{ .Key }})%s", ColorCyan, IconArrow, ColorReset, ColorDim, ColorReset),
		Inactive: fmt.Sprintf("  {{ .Icon }} {{ .Label }} %s({{ .Key }})%s", ColorDim, ColorReset),
		Selected: fmt.Sprintf("%s%s {{ .Icon }} {{ .Label }}%s", ColorGreen, IconCheck, ColorReset),
	}

	// Custom searcher for keyboard shortcuts
	searcher := func(input string, index int) bool {
		item := items[index]
		input = strings.ToLower(input)
		// Match by key shortcut or label
		return strings.HasPrefix(strings.ToLower(item.Key), input) ||
			strings.HasPrefix(strings.ToLower(item.Label), input)
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("%sWhat would you like to do?%s", ColorBold, ColorReset),
		Items:     items,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
		// Start with cursor at position 0 so the user can press enter to execute immediately
		CursorPos:    0,
		HideHelp:     true,
		HideSelected: false,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		// If user presses Ctrl+C or there's an error, default to quit
		if err == promptui.ErrInterrupt {
			return ActionQuit
		}
		// Fallback to text-based prompt if interactive mode fails
		return promptActionFallback()
	}

	return items[idx].Action
}

// promptActionFallback provides text-based input as a fallback
func promptActionFallback() Action {
	fmt.Printf("%sWhat would you like to do?%s\n", ColorBold, ColorReset)
	fmt.Printf("  [%se%s]xecute  |  [%sc%s]opy  |  e[%sd%s]it  |  [%sr%s]e-prompt  |  [%sq%s]uit\n",
		ColorYellow, ColorReset,
		ColorYellow, ColorReset,
		ColorYellow, ColorReset,
		ColorYellow, ColorReset,
		ColorYellow, ColorReset)
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ActionQuit
	}

	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "e", "execute":
		return ActionExecute
	case "c", "copy":
		return ActionCopy
	case "d", "edit":
		return ActionEdit
	case "r", "re-prompt", "reprompt":
		return ActionReprompt
	case "q", "quit", "exit":
		return ActionQuit
	default:
		// Default to execute if user just presses enter
		if input == "" {
			return ActionExecute
		}
		fmt.Printf("%sInvalid option. Please try again.%s\n", ColorDim, ColorReset)
		return promptActionFallback()
	}
}

// PromptActionForCommand prompts the user for an action, with safety checks for the command
func PromptActionForCommand(command string) Action {
	action := PromptAction()

	// If executing a dangerous command, require explicit confirmation
	if action == ActionExecute && safety.RequiresConfirmation(command) {
		fmt.Printf("%s%sThis is a potentially dangerous command. Type 'yes' to confirm: %s", ColorBold, ColorRed, ColorReset)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return ActionQuit
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input != "yes" {
			fmt.Printf("%sExecution cancelled.%s\n", ColorYellow, ColorReset)
			return ActionQuit
		}
	}

	return action
}

// PromptEdit prompts the user to edit the command
func PromptEdit(currentCommand string) string {
	fmt.Printf("%s%s Edit command:%s\n", ColorBold, IconEdit, ColorReset)
	fmt.Printf("%sCurrent: %s%s\n", ColorDim, currentCommand, ColorReset)

	prompt := promptui.Prompt{
		Label:     "New command",
		Default:   currentCommand,
		AllowEdit: true,
		Templates: &promptui.PromptTemplates{
			Prompt:  fmt.Sprintf("%s{{ . }}:%s ", ColorCyan, ColorReset),
			Valid:   fmt.Sprintf("%s{{ . }}:%s ", ColorGreen, ColorReset),
			Invalid: fmt.Sprintf("%s{{ . }}:%s ", ColorRed, ColorReset),
			Success: fmt.Sprintf("%s%s {{ . }}:%s ", ColorGreen, IconCheck, ColorReset),
		},
	}

	result, err := prompt.Run()
	if err != nil {
		// On error or Ctrl+C, return the original command
		return currentCommand
	}

	if strings.TrimSpace(result) == "" {
		return currentCommand
	}

	return result
}

// PromptReprompt prompts the user for a new query
func PromptReprompt() string {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s%s New request%s", ColorBold, IconRefresh, ColorReset),
		Templates: &promptui.PromptTemplates{
			Prompt:  fmt.Sprintf("%s{{ . }}:%s ", ColorCyan, ColorReset),
			Valid:   fmt.Sprintf("%s{{ . }}:%s ", ColorGreen, ColorReset),
			Invalid: fmt.Sprintf("%s{{ . }}:%s ", ColorRed, ColorReset),
			Success: fmt.Sprintf("%s%s {{ . }}:%s ", ColorGreen, IconCheck, ColorReset),
		},
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("request cannot be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(result)
}

// CopyToClipboard copies the command to the clipboard
func CopyToClipboard(command string) error {
	err := clipboard.WriteAll(command)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	fmt.Println(SuccessMessage("Copied to clipboard!"))
	return nil
}

// ExecuteCommand executes the command in the current shell
func ExecuteCommand(command string, shellInfo shell.ShellInfo) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		if shellInfo.Shell == shell.ShellPowerShell {
			cmd = exec.Command("powershell", "-NoProfile", "-Command", command)
		} else {
			cmd = exec.Command("cmd", "/C", command)
		}
	default:
		// Unix-like systems - use dynamic shell path detection
		shellPath := getShellPath(shellInfo.Shell)
		cmd = exec.Command(shellPath, "-c", command)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%s%sExecuting...%s\n", ColorDim, ColorYellow, ColorReset)
	fmt.Println()

	err := cmd.Run()

	// Show undo suggestion after execution
	fmt.Println()
	undoSuggestion := undo.GetUndoSuggestion(command)
	if undoSuggestion.CanUndo {
		fmt.Println(undo.FormatUndoSuggestion(undoSuggestion))
	}

	return err
}

// ExecuteCommandWithErrorRecovery executes a command and offers help if it fails
func ExecuteCommandWithErrorRecovery(command string, shellInfo shell.ShellInfo) (error, bool) {
	err := ExecuteCommand(command, shellInfo)

	if err != nil {
		fmt.Println()
		fmt.Printf("%s%sCommand failed with error: %s%s\n", ColorRed, ColorBold, err.Error(), ColorReset)
		fmt.Printf("%sWould you like help diagnosing this error? [y/N]: %s", ColorYellow, ColorReset)

		reader := bufio.NewReader(os.Stdin)
		input, readErr := reader.ReadString('\n')
		if readErr != nil {
			return err, false
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			return err, true // Signal that error recovery is requested
		}
	}

	return err, false
}

// ShowError displays an error message
func ShowError(err error) {
	fmt.Println(ErrorMessage(err.Error()))
}

// ShowSpinner shows a simple loading indicator with elapsed time
func ShowSpinner(message string) func() {
	done := make(chan bool)
	startTime := time.Now()

	go func() {
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r" + strings.Repeat(" ", len(message)+20) + "\r")
				return
			default:
				elapsed := time.Since(startTime).Seconds()
				timeStr := FormatDuration(elapsed)
				fmt.Printf("\r%s%s %s %s[%s]%s", ColorCyan, ProgressDots(i), message, ColorDim, timeStr, ColorReset)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		done <- true
	}
}

// getShellPath returns the path to the shell executable using dynamic detection
func getShellPath(shellType shell.ShellType) string {
	// First, try to use the $SHELL environment variable if it matches the detected shell
	envShell := os.Getenv("SHELL")
	if envShell != "" {
		// Verify the env shell matches what we expect
		shellName := shellTypeToName(shellType)
		if strings.Contains(envShell, shellName) {
			return envShell
		}
	}

	// Try to find the shell in PATH using exec.LookPath
	shellName := shellTypeToName(shellType)
	if path, err := exec.LookPath(shellName); err == nil {
		return path
	}

	// Fallback to common paths
	fallbackPaths := map[shell.ShellType][]string{
		shell.ShellBash: {"/bin/bash", "/usr/bin/bash", "/usr/local/bin/bash"},
		shell.ShellZsh:  {"/bin/zsh", "/usr/bin/zsh", "/usr/local/bin/zsh"},
		shell.ShellFish: {"/usr/bin/fish", "/usr/local/bin/fish", "/bin/fish"},
	}

	if paths, ok := fallbackPaths[shellType]; ok {
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	// Ultimate fallback to /bin/sh
	return "/bin/sh"
}

// shellTypeToName converts a ShellType to the shell executable name
func shellTypeToName(shellType shell.ShellType) string {
	switch shellType {
	case shell.ShellBash:
		return "bash"
	case shell.ShellZsh:
		return "zsh"
	case shell.ShellFish:
		return "fish"
	default:
		return "sh"
	}
}

