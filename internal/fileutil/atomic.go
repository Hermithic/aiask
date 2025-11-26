package fileutil

import (
	"os"
	"path/filepath"
)

// AtomicWriteFile writes data to a file atomically by writing to a temp file
// and then renaming it. This prevents corruption from concurrent writes or crashes.
func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error {
	// Create temp file in the same directory to ensure atomic rename works
	dir := filepath.Dir(filename)
	
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	
	// Create temp file
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	
	// Clean up temp file on any error
	defer func() {
		if tmpFile != nil {
			tmpFile.Close()
			os.Remove(tmpName)
		}
	}()
	
	// Set permissions
	if err := tmpFile.Chmod(perm); err != nil {
		return err
	}
	
	// Write data
	if _, err := tmpFile.Write(data); err != nil {
		return err
	}
	
	// Sync to disk
	if err := tmpFile.Sync(); err != nil {
		return err
	}
	
	// Close before rename
	if err := tmpFile.Close(); err != nil {
		return err
	}
	tmpFile = nil // Prevent deferred cleanup
	
	// Atomic rename
	if err := os.Rename(tmpName, filename); err != nil {
		os.Remove(tmpName) // Clean up on rename failure
		return err
	}
	
	return nil
}

