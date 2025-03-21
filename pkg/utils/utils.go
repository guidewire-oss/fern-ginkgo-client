package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindGitRoot looks for a .git folder in the startDir and works up the tree until it finds one.
func FindGitRoot(startDir string) (string, error) {
	dir, err := ToAbsolutePath(startDir)
	if err != nil {
		fmt.Printf("Unable to convert %s to an absolute path - %v\n", startDir, err)
		return "", err
	}
	for {
		if isValidGitRoot(dir) {
			return dir, nil // Found the valid Git root
		}

		// Move to the parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached filesystem root - the parent of root is root.
		}
		dir = parent
	}
	return "", fmt.Errorf("no valid .git directory found")
}

// isValidGitRoot checks if the directory is a real git repo and not a submodule reference
func isValidGitRoot(dir string) bool {
	// Ensure .git is a directory, not a file (submodules use a file reference)
	info, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ToAbsolutePath converts a given path to an absolute path, handling ".", relative paths, and "~".
func ToAbsolutePath(path string) (string, error) {
	// Expand "~" to home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, strings.TrimPrefix(path, "~"))
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
