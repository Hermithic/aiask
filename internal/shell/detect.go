package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ShellType represents the detected shell type
type ShellType string

const (
	ShellPowerShell ShellType = "powershell"
	ShellCmd        ShellType = "cmd"
	ShellBash       ShellType = "bash"
	ShellZsh        ShellType = "zsh"
	ShellFish       ShellType = "fish"
	ShellUnknown    ShellType = "unknown"
)

// ShellInfo contains information about the detected shell and OS
type ShellInfo struct {
	Shell ShellType
	OS    string
}

// Detect detects the current shell and operating system
func Detect() ShellInfo {
	info := ShellInfo{
		OS:    runtime.GOOS,
		Shell: ShellUnknown,
	}

	// Windows detection
	if runtime.GOOS == "windows" {
		info.Shell = detectWindowsShell()
		return info
	}

	// Unix-like systems (Linux, macOS)
	info.Shell = detectUnixShell()
	return info
}

// detectWindowsShell detects the shell on Windows
func detectWindowsShell() ShellType {
	// Check for PowerShell
	// PSModulePath is typically set in PowerShell environments
	if os.Getenv("PSModulePath") != "" {
		// Check if it's PowerShell Core (pwsh) or Windows PowerShell
		return ShellPowerShell
	}

	// Check for CMD
	// PROMPT is typically set in CMD
	if os.Getenv("PROMPT") != "" {
		return ShellCmd
	}

	// Check COMSPEC for CMD
	comspec := os.Getenv("COMSPEC")
	if strings.Contains(strings.ToLower(comspec), "cmd.exe") {
		return ShellCmd
	}

	// Check if running under WSL (Windows Subsystem for Linux)
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return detectUnixShell()
	}

	// Default to PowerShell on Windows
	return ShellPowerShell
}

// detectUnixShell detects the shell on Unix-like systems
func detectUnixShell() ShellType {
	// Check SHELL environment variable
	shell := os.Getenv("SHELL")
	if shell != "" {
		base := strings.ToLower(filepath.Base(shell))
		switch {
		case strings.Contains(base, "bash"):
			return ShellBash
		case strings.Contains(base, "zsh"):
			return ShellZsh
		case strings.Contains(base, "fish"):
			return ShellFish
		case strings.Contains(base, "pwsh"), strings.Contains(base, "powershell"):
			return ShellPowerShell
		}
	}

	// Check for specific environment variables
	if os.Getenv("BASH_VERSION") != "" {
		return ShellBash
	}
	if os.Getenv("ZSH_VERSION") != "" {
		return ShellZsh
	}
	if os.Getenv("FISH_VERSION") != "" {
		return ShellFish
	}

	// Default to bash on Unix
	return ShellBash
}

// GetOSName returns a human-readable OS name
func GetOSName() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}

// GetShellName returns a human-readable shell name
func GetShellName(shell ShellType) string {
	switch shell {
	case ShellPowerShell:
		return "PowerShell"
	case ShellCmd:
		return "Command Prompt (CMD)"
	case ShellBash:
		return "Bash"
	case ShellZsh:
		return "Zsh"
	case ShellFish:
		return "Fish"
	default:
		return "Unknown Shell"
	}
}

// GetShellDescription returns a description for the system prompt
func (s ShellInfo) GetDescription() string {
	shellName := GetShellName(s.Shell)
	osName := GetOSName()
	return shellName + " on " + osName
}

