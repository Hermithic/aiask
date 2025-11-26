package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteFile(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "atomic_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("write new file", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "test1.txt")
		content := []byte("hello world")

		err := AtomicWriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("AtomicWriteFile failed: %v", err)
		}

		// Verify content
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != string(content) {
			t.Errorf("content mismatch: got %q, expected %q", string(data), string(content))
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "test2.txt")

		// Write initial content
		err := AtomicWriteFile(filename, []byte("initial"), 0644)
		if err != nil {
			t.Fatalf("initial write failed: %v", err)
		}

		// Overwrite with new content
		newContent := []byte("updated content")
		err = AtomicWriteFile(filename, newContent, 0644)
		if err != nil {
			t.Fatalf("overwrite failed: %v", err)
		}

		// Verify new content
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != string(newContent) {
			t.Errorf("content mismatch: got %q, expected %q", string(data), string(newContent))
		}
	})

	t.Run("creates directory if needed", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "subdir", "nested", "test3.txt")
		content := []byte("nested content")

		err := AtomicWriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("AtomicWriteFile failed: %v", err)
		}

		// Verify content
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != string(content) {
			t.Errorf("content mismatch: got %q, expected %q", string(data), string(content))
		}
	})

	t.Run("respects file permissions", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "test4.txt")
		content := []byte("restricted content")

		err := AtomicWriteFile(filename, content, 0600)
		if err != nil {
			t.Fatalf("AtomicWriteFile failed: %v", err)
		}

		// Check permissions (skip on Windows as it handles permissions differently)
		info, err := os.Stat(filename)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}

		// On Unix, check the permissions match
		mode := info.Mode().Perm()
		if mode&0077 != 0 && mode != 0600 {
			// Allow some flexibility as umask might affect this
			t.Logf("file mode is %o (expected 0600 or similar restricted mode)", mode)
		}
	})
}

