// internal/lfs/manager.go
package lfs

import (
    "fmt"
    "os/exec"
    "path/filepath"
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
    // Install Git LFS in the repository
    cmd := exec.Command("git", "lfs", "install")
    cmd.Dir = m.RepoPath
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to initialize Git LFS: %w", err)
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
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("failed to track .%s files: %w", format, err)
        }
    }
    
    // Add .gitattributes to version control
    cmd = exec.Command("git", "add", ".gitattributes")
    cmd.Dir = m.RepoPath
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to add .gitattributes: %w", err)
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
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to track .%s files: %w", ext, err)
    }
    
    return nil
}

// GetObjectID retrieves the LFS object ID for a file
func (m *Manager) GetObjectID(filePath string) (string, error) {
    cmd := exec.Command("git", "lfs", "pointer", "--file", filePath)
    cmd.Dir = m.RepoPath
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get LFS pointer: %w", err)
    }
    
    // In real implementation, parse the output to extract the OID
    // For now, return simplified placeholder
    return string(output), nil
}

// Pull fetches LFS objects
func (m *Manager) Pull() error {
    cmd := exec.Command("git", "lfs", "pull")
    cmd.Dir = m.RepoPath
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to pull LFS objects: %w", err)
    }
    
    return nil
}

// Push uploads LFS objects
func (m *Manager) Push() error {
    cmd := exec.Command("git", "lfs", "push", "--all", "origin", "master")
    cmd.Dir = m.RepoPath
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to push LFS objects: %w", err)
    }
    
    return nil
}