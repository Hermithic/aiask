package history

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hermithic/aiask/internal/config"
	"gopkg.in/yaml.v3"
)

// MaxHistoryEntries is the maximum number of history entries to keep
const MaxHistoryEntries = 100

// Entry represents a single history entry
type Entry struct {
	Timestamp time.Time `yaml:"timestamp"`
	Prompt    string    `yaml:"prompt"`
	Command   string    `yaml:"command"`
	Shell     string    `yaml:"shell"`
	Executed  bool      `yaml:"executed"`
}

// History represents the command history
type History struct {
	Entries []Entry `yaml:"entries"`
}

// GetHistoryPath returns the path to the history file
func GetHistoryPath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "history.yaml"), nil
}

// Load loads the history from the history file
func Load() (*History, error) {
	historyPath, err := GetHistoryPath()
	if err != nil {
		return nil, err
	}

	// Return empty history if file doesn't exist
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return &History{Entries: []Entry{}}, nil
	}

	data, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	history := &History{}
	if err := yaml.Unmarshal(data, history); err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	return history, nil
}

// Save saves the history to the history file
func (h *History) Save() error {
	historyPath, err := GetHistoryPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(historyPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// Add adds a new entry to the history
func (h *History) Add(entry Entry) {
	// Set timestamp if not set
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Prepend to history (newest first)
	h.Entries = append([]Entry{entry}, h.Entries...)

	// Trim to max entries
	if len(h.Entries) > MaxHistoryEntries {
		h.Entries = h.Entries[:MaxHistoryEntries]
	}
}

// Clear clears all history entries
func (h *History) Clear() {
	h.Entries = []Entry{}
}

// GetRecent returns the N most recent entries
func (h *History) GetRecent(n int) []Entry {
	if n <= 0 || n > len(h.Entries) {
		return h.Entries
	}
	return h.Entries[:n]
}

// Search searches history entries by prompt substring
func (h *History) Search(query string) []Entry {
	var results []Entry
	for _, entry := range h.Entries {
		if containsIgnoreCase(entry.Prompt, query) || containsIgnoreCase(entry.Command, query) {
			results = append(results, entry)
		}
	}
	return results
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 len(substr) == 0 ||
		 findIgnoreCase(s, substr) >= 0)
}

func findIgnoreCase(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	
	// Simple case-insensitive search
	sLower := toLower(s)
	substrLower := toLower(substr)
	
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// AddEntry is a convenience function to add an entry and save
func AddEntry(prompt, command, shell string, executed bool) error {
	h, err := Load()
	if err != nil {
		return err
	}

	h.Add(Entry{
		Prompt:   prompt,
		Command:  command,
		Shell:    shell,
		Executed: executed,
	})

	return h.Save()
}

