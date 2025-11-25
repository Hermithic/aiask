package context

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MaxDirEntries is the maximum number of directory entries to include
const MaxDirEntries = 50

// MaxFileNameLength is the maximum file name length to display
const MaxFileNameLength = 60

// DirEntry represents a directory entry
type DirEntry struct {
	Name  string
	IsDir bool
	Size  int64
}

// GetDirectoryContext returns a string describing the current directory contents
func GetDirectoryContext() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return ""
	}

	if len(entries) == 0 {
		return fmt.Sprintf("Current directory: %s (empty)", cwd)
	}

	var dirEntries []DirEntry
	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		size := int64(0)
		if err == nil {
			size = info.Size()
		}

		dirEntries = append(dirEntries, DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  size,
		})
	}

	// Sort directories first, then files
	sort.Slice(dirEntries, func(i, j int) bool {
		if dirEntries[i].IsDir != dirEntries[j].IsDir {
			return dirEntries[i].IsDir
		}
		return dirEntries[i].Name < dirEntries[j].Name
	})

	// Limit entries
	truncated := false
	if len(dirEntries) > MaxDirEntries {
		dirEntries = dirEntries[:MaxDirEntries]
		truncated = true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Current directory: %s\n", cwd))
	sb.WriteString("Contents:\n")

	for _, entry := range dirEntries {
		name := entry.Name
		if len(name) > MaxFileNameLength {
			name = name[:MaxFileNameLength-3] + "..."
		}

		if entry.IsDir {
			sb.WriteString(fmt.Sprintf("  [DIR] %s/\n", name))
		} else {
			sb.WriteString(fmt.Sprintf("  [FILE] %s (%s)\n", name, formatSize(entry.Size)))
		}
	}

	if truncated {
		sb.WriteString(fmt.Sprintf("  ... and more files (showing first %d)\n", MaxDirEntries))
	}

	return sb.String()
}

// GetCWD returns the current working directory
func GetCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

// GetCWDBaseName returns the base name of the current working directory
func GetCWDBaseName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Base(cwd)
}

// HasFiles checks if the current directory has files matching the pattern
func HasFiles(pattern string) bool {
	matches, err := filepath.Glob(pattern)
	return err == nil && len(matches) > 0
}

// CountFiles counts files matching the pattern in the current directory
func CountFiles(pattern string) int {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}
	return len(matches)
}

// formatSize formats a file size in human-readable format
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// IsFileRelatedPrompt checks if a prompt seems to be about files/directories
func IsFileRelatedPrompt(prompt string) bool {
	fileKeywords := []string{
		"file", "files", "directory", "directories", "folder", "folders",
		"find", "search", "list", "delete", "remove", "copy", "move",
		"rename", "compress", "zip", "unzip", "extract", "archive",
		"size", "large", "small", "old", "new", "recent", "modified",
		"create", "touch", "mkdir", "rmdir", "ls", "dir",
	}

	promptLower := strings.ToLower(prompt)
	for _, keyword := range fileKeywords {
		if strings.Contains(promptLower, keyword) {
			return true
		}
	}
	return false
}

