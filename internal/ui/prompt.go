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

// DisplayCommand shows the suggested command to the user
func DisplayCommand(command string) {
	fmt.Println()
	fmt.Printf("%s%sSuggested command:%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Println()

	// Display the command with highlighting
	lines := strings.Split(command, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("  %s%s%s\n", ColorGreen, line, ColorReset)
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

// PromptAction prompts the user for an action
func PromptAction() Action {
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
		return PromptAction()
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
	fmt.Printf("%sEdit the command (current command shown below):%s\n", ColorBold, ColorReset)
	fmt.Printf("%s%s%s\n", ColorDim, currentCommand, ColorReset)
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return currentCommand
	}

	edited := strings.TrimSpace(input)
	if edited == "" {
		return currentCommand
	}
	return edited
}

// PromptReprompt prompts the user for a new query
func PromptReprompt() string {
	fmt.Printf("%sEnter your new request:%s\n", ColorBold, ColorReset)
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(input)
}

// CopyToClipboard copies the command to the clipboard
func CopyToClipboard(command string) error {
	err := clipboard.WriteAll(command)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	fmt.Printf("%s✓ Copied to clipboard!%s\n", ColorGreen, ColorReset)
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
		// Unix-like systems
		shellPath := "/bin/sh"
		if shellInfo.Shell == shell.ShellZsh {
			shellPath = "/bin/zsh"
		} else if shellInfo.Shell == shell.ShellBash {
			shellPath = "/bin/bash"
		} else if shellInfo.Shell == shell.ShellFish {
			shellPath = "/usr/bin/fish"
		}
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
	fmt.Printf("%sError: %s%s\n", ColorYellow, err.Error(), ColorReset)
}

// ShowSpinner shows a simple loading indicator
func ShowSpinner(message string) func() {
	done := make(chan bool)
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r" + strings.Repeat(" ", len(message)+5) + "\r")
				return
			default:
				fmt.Printf("\r%s%s %s%s", ColorCyan, chars[i%len(chars)], message, ColorReset)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		done <- true
	}
}

