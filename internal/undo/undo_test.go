package undo

import (
	"strings"
	"testing"
)

func TestGetUndoSuggestion(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		canUndo     bool
		undoContains string
	}{
		// Git operations
		{"git commit", "git commit -m 'message'", true, "git reset HEAD~1"},
		{"git add", "git add file.txt", true, "git reset file.txt"},
		{"git add all", "git add .", true, "git reset ."},
		{"git stash", "git stash", true, "git stash pop"},
		{"git stash push", "git stash push", true, "git stash pop"},
		{"git checkout -b", "git checkout -b feature", true, "git branch -d feature"},
		{"git merge", "git merge main", true, "git reset --hard HEAD~1"},

		// File operations
		{"mv file", "mv old.txt new.txt", true, "mv new.txt old.txt"},
		{"cp file", "cp src.txt dest.txt", true, "rm dest.txt"},
		{"cp recursive", "cp -r src/ dest/", true, "rm -r dest/"},
		{"mkdir", "mkdir newdir", true, "rmdir newdir"},
		{"mkdir -p", "mkdir -p path/to/dir", true, "rmdir path/to/dir"},
		{"touch", "touch newfile.txt", true, "rm newfile.txt"},
		{"ln", "ln -s target link", true, "rm link"},

		// Package managers
		{"apt install", "apt install nginx", true, "apt remove nginx"},
		{"apt-get install", "apt-get install nginx", true, "apt-get remove nginx"},
		{"brew install", "brew install wget", true, "brew uninstall wget"},
		{"npm install", "npm install express", true, "npm uninstall express"},
		{"npm install -g", "npm install -g typescript", true, "npm uninstall -g typescript"},
		{"pip install", "pip install requests", true, "pip uninstall requests"},

		// Services
		{"systemctl start", "systemctl start nginx", true, "systemctl stop nginx"},
		{"systemctl stop", "systemctl stop nginx", true, "systemctl start nginx"},
		{"systemctl enable", "systemctl enable nginx", true, "systemctl disable nginx"},

		// Docker
		{"docker run", "docker run --name mycontainer nginx", true, "docker stop mycontainer"},
		{"docker start", "docker start mycontainer", true, "docker stop mycontainer"},

		// No undo available
		{"ls command", "ls -la", false, ""},
		{"cat command", "cat file.txt", false, ""},
		{"echo command", "echo hello", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUndoSuggestion(tt.command)

			if result.CanUndo != tt.canUndo {
				t.Errorf("GetUndoSuggestion(%q).CanUndo = %v, expected %v", tt.command, result.CanUndo, tt.canUndo)
			}

			if tt.canUndo && !strings.Contains(result.UndoCommand, tt.undoContains) {
				t.Errorf("GetUndoSuggestion(%q).UndoCommand = %q, expected to contain %q", tt.command, result.UndoCommand, tt.undoContains)
			}

			if result.Original != tt.command {
				t.Errorf("GetUndoSuggestion(%q).Original = %q, expected %q", tt.command, result.Original, tt.command)
			}
		})
	}
}

func TestGetUndoSuggestionNpmGlobalFlag(t *testing.T) {
	// Test that npm -g flag is handled correctly
	result := GetUndoSuggestion("npm install -g typescript")
	if !result.CanUndo {
		t.Error("npm install -g should be undoable")
	}
	if !strings.Contains(result.UndoCommand, "-g") {
		t.Errorf("npm uninstall should contain -g flag, got: %q", result.UndoCommand)
	}
	if !strings.Contains(result.UndoCommand, "typescript") {
		t.Errorf("npm uninstall should contain package name, got: %q", result.UndoCommand)
	}
	// Check no double spaces
	if strings.Contains(result.UndoCommand, "  ") {
		t.Errorf("npm uninstall should not have double spaces, got: %q", result.UndoCommand)
	}
}

func TestGetUndoSuggestionCpRecursive(t *testing.T) {
	// Test that cp -r generates rm -r
	result := GetUndoSuggestion("cp -r source/ destination/")
	if !result.CanUndo {
		t.Error("cp -r should be undoable")
	}
	if !strings.Contains(result.UndoCommand, "rm -r") {
		t.Errorf("cp -r undo should use rm -r, got: %q", result.UndoCommand)
	}

	// Test that cp without -r generates rm without -r
	result = GetUndoSuggestion("cp source.txt destination.txt")
	if !result.CanUndo {
		t.Error("cp should be undoable")
	}
	if strings.Contains(result.UndoCommand, "rm -r") {
		t.Errorf("cp (non-recursive) undo should not use rm -r, got: %q", result.UndoCommand)
	}
}

func TestFormatUndoSuggestion(t *testing.T) {
	// Test with undo available
	suggestion := UndoSuggestion{
		Original:    "git commit -m 'test'",
		UndoCommand: "git reset HEAD~1",
		Description: "Undo the last commit",
		CanUndo:     true,
	}
	result := FormatUndoSuggestion(suggestion)
	if result == "" {
		t.Error("FormatUndoSuggestion should return non-empty string when CanUndo is true")
	}
	if !strings.Contains(result, "git reset HEAD~1") {
		t.Error("FormatUndoSuggestion should contain the undo command")
	}

	// Test without undo available
	suggestion = UndoSuggestion{
		Original:    "ls -la",
		UndoCommand: "",
		Description: "No automatic undo available",
		CanUndo:     false,
	}
	result = FormatUndoSuggestion(suggestion)
	if result != "" {
		t.Errorf("FormatUndoSuggestion should return empty string when CanUndo is false, got: %q", result)
	}
}

