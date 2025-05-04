// internal/lfs/manager.go
package lfs

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager provides Git LFS functionality
type Manager struct {
	RepoPath string
}

// New creates a new LFS Manager
func New(repoPath string) *Manager {
	return &Manager{
		RepoPath: repoPath,
	}
}

// Initialize sets up Git LFS for the repository
func (m *Manager) Initialize() error {
	// Check if Git LFS is installed
	if err := m.checkLFSInstalled(); err != nil {
		return err
	}

	// Install Git LFS in the repository
	cmd := exec.Command("git", "lfs", "install")
	cmd.Dir = m.RepoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize Git LFS: %w (output: %s)", err, string(out))
	}

	// Configure LFS to track audio file formats
	formats := []string{
		// Audio formats
		"wav", "aif", "mp3", "ogg", "flac",
		// DAW project formats
		"als", "ptx", "sesx", "cpr", "rpp",
	}

	for _, format := range formats {
		cmd = exec.Command("git", "lfs", "track", fmt.Sprintf("*.%s", format))
		cmd.Dir = m.RepoPath
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to track .%s files: %w (output: %s)", format, err, string(out))
		}
	}

	// Add .gitattributes to version control
	cmd = exec.Command("git", "add", ".gitattributes")
	cmd.Dir = m.RepoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add .gitattributes: %w (output: %s)", err, string(out))
	}

	return nil
}

// checkLFSInstalled checks if Git LFS is installed on the system
func (m *Manager) checkLFSInstalled() error {
	cmd := exec.Command("git", "lfs", "version")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Git LFS is not installed: %w\nPlease install Git LFS: https://git-lfs.github.com", err)
	} else if !strings.Contains(string(out), "git-lfs") {
		return fmt.Errorf("Git LFS installation appears to be incomplete: %s", string(out))
	}
	return nil
}

// TrackFile ensures a specific file is tracked by Git LFS
func (m *Manager) TrackFile(filePath string) error {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file has no extension: %s", filePath)
	}

	// Remove leading dot from extension
	ext = ext[1:]

	cmd := exec.Command("git", "lfs", "track", fmt.Sprintf("*.%s", ext))
	cmd.Dir = m.RepoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to track .%s files: %w (output: %s)", ext, err, string(out))
	}

	return nil
}

// GetObjectID retrieves the LFS object ID for a file
func (m *Manager) GetObjectID(filePath string) (string, error) {
	cmd := exec.Command("git", "lfs", "pointer", "--file", filePath)
	cmd.Dir = m.RepoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get LFS pointer: %w (output: %s)", err, string(output))
	}

	// In real implementation, parse the output to extract the OID
	// For now, return simplified placeholder
	return string(output), nil
}

// Pull fetches LFS objects
func (m *Manager) Pull() error {
	cmd := exec.Command("git", "lfs", "pull")
	cmd.Dir = m.RepoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull LFS objects: %w (output: %s)", err, string(out))
	}

	return nil
}

// Push uploads LFS objects
func (m *Manager) Push() error {
	// First try to get the current branch
	branchCmd := exec.Command("git", "branch", "--show-current")
	branchCmd.Dir = m.RepoPath
	branchOutput, err := branchCmd.Output()

	branch := "main" // Default to main
	if err == nil && len(branchOutput) > 0 {
		branch = strings.TrimSpace(string(branchOutput))
	}

	cmd := exec.Command("git", "lfs", "push", "--all", "origin", branch)
	cmd.Dir = m.RepoPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push LFS objects: %w (output: %s)", err, string(out))
	}

	return nil
}
