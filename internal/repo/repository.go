// internal/repo/repository.go
package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shantykins/musicvcs/internal/lfs"
	"github.com/shantykins/musicvcs/pkg/models"
)

// Repository manages the version control operations
type Repository struct {
	RootPath         string
	MainMixPath      string
	BranchMixPath    string
	ProjectFilesPath string
	MixMaster        string // ID of the user who can update main mix
	LFS              *lfs.Manager
}

// New creates a new Repository instance
func New(rootPath string, mixMaster string) *Repository {
	repo := &Repository{
		RootPath:         rootPath,
		MainMixPath:      filepath.Join(rootPath, "main-mix"),
		BranchMixPath:    filepath.Join(rootPath, "branch-mix"),
		ProjectFilesPath: filepath.Join(rootPath, "project-files"),
		MixMaster:        mixMaster,
		LFS:              lfs.New(rootPath),
	}

	return repo
}

// Initialize sets up the repository structure
func (r *Repository) Initialize() error {
	// Create directory structure
	dirs := []string{
		r.MainMixPath,
		r.BranchMixPath,
		r.ProjectFilesPath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create instrument subdirectories
	instruments := []string{"drums", "vocals", "guitars", "bass", "keys", "fx"}
	for _, instrument := range instruments {
		instrumentDir := filepath.Join(r.ProjectFilesPath, instrument)
		if err := os.MkdirAll(instrumentDir, 0755); err != nil {
			return fmt.Errorf("failed to create instrument directory %s: %w", instrumentDir, err)
		}
	}

	// Initialize Git repository if not already initialized
	if _, err := os.Stat(filepath.Join(r.RootPath, ".git")); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = r.RootPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize Git repository: %w", err)
		}

		// Set Git configuration for commits
		if err := r.ensureGitConfig(); err != nil {
			return err
		}

		// Initialize Git LFS
		if err := r.LFS.Initialize(); err != nil {
			return err
		}

		// Initial commit
		cmd = exec.Command("git", "commit", "--allow-empty", "-m", "Initialize music version control")
		cmd.Dir = r.RootPath
		if err := cmd.Run(); err != nil {
			// Try to recover if the error is about user.name or user.email not being set
			if strings.Contains(err.Error(), "Please tell me who you are") {
				if err := r.ensureGitConfig(); err != nil {
					return err
				}
				// Try the commit again
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to make initial commit: %w", err)
				}
			} else {
				return fmt.Errorf("failed to make initial commit: %w", err)
			}
		}
	}

	return nil
}

// ensureGitConfig makes sure Git is configured with user.name and user.email
func (r *Repository) ensureGitConfig() error {
	// Check if user.name is set
	cmd := exec.Command("git", "config", "user.name")
	cmd.Dir = r.RootPath
	output, err := cmd.Output()

	if err != nil || len(output) == 0 {
		// Set user.name
		cmd = exec.Command("git", "config", "--local", "user.name", r.MixMaster)
		cmd.Dir = r.RootPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.name: %w", err)
		}
	}

	// Check if user.email is set
	cmd = exec.Command("git", "config", "user.email")
	cmd.Dir = r.RootPath
	output, err = cmd.Output()

	if err != nil || len(output) == 0 {
		// Set user.email to a placeholder if not set
		cmd = exec.Command("git", "config", "--local", "user.email", fmt.Sprintf("%s@musicvcs", r.MixMaster))
		cmd.Dir = r.RootPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.email: %w", err)
		}
	}

	return nil
}

// AddTrack adds a track to the project files
func (r *Repository) AddTrack(filePath string, trackType string, userId string) (string, error) {
	// Ensure Git is configured properly
	if err := r.ensureGitConfig(); err != nil {
		return "", err
	}

	// Validate track type
	validTypes := map[string]bool{
		"drums": true, "vocals": true, "guitars": true,
		"bass": true, "keys": true, "fx": true,
	}

	if !validTypes[trackType] {
		return "", fmt.Errorf("invalid track type: %s", trackType)
	}

	// Copy file to project directory
	fileName := filepath.Base(filePath)
	targetDir := filepath.Join(r.ProjectFilesPath, trackType)
	targetPath := filepath.Join(targetDir, fileName)

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %w", err)
	}

	// Copy file
	input, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(targetPath, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write target file: %w", err)
	}

	// Track file with Git LFS
	if err := r.LFS.TrackFile(targetPath); err != nil {
		return "", err
	}

	// Create track metadata
	trackID := uuid.New().String()
	metadata := models.TrackMetadata{
		ID:        trackID,
		FileName:  fileName,
		Type:      trackType,
		CreatedAt: time.Now(),
		CreatedBy: userId,
		Version:   1,
		// Other fields would be populated from actual audio analysis
		SampleRate:   44100, // Default values for example
		BitDepth:     16,
		Channels:     2,
		Tags:         []string{},
		Dependencies: []string{},
	}

	// Save metadata
	metadataPath := targetPath + ".json"
	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Add files to Git
	cmd := exec.Command("git", "add", targetPath, metadataPath)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to add files to Git: %w (output: %s)", err, string(out))
	}

	// Create commit
	commitMsg := fmt.Sprintf("Add %s track: %s", trackType, fileName)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		// Check if it's a "nothing to commit" error, which is not a real error
		if strings.Contains(string(out), "nothing to commit") {
			return trackID, nil
		}
		return "", fmt.Errorf("failed to commit files: %w (output: %s)", err, string(out))
	}

	return trackID, nil
}

// CreateBranchMix creates a new branch mix
func (r *Repository) CreateBranchMix(name string, description string, userId string) error {
	// Ensure Git is configured properly
	if err := r.ensureGitConfig(); err != nil {
		return err
	}

	branchDir := filepath.Join(r.BranchMixPath, name)

	// Check if branch already exists
	if _, err := os.Stat(branchDir); err == nil {
		return fmt.Errorf("branch mix already exists: %s", name)
	}

	// Create branch directory
	if err := os.MkdirAll(branchDir, 0755); err != nil {
		return fmt.Errorf("failed to create branch directory: %w", err)
	}

	// Create branch metadata
	metadata := models.BranchMetadata{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		CreatedBy:   userId,
		Tracks:      []string{},
	}

	// Save metadata
	metadataPath := filepath.Join(branchDir, "branch-info.json")
	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Add to Git
	cmd := exec.Command("git", "add", branchDir)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add branch to Git: %w (output: %s)", err, string(out))
	}

	// Create commit
	commitMsg := fmt.Sprintf("Create branch mix: %s", name)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		// Check if it's a "nothing to commit" error, which is not a real error
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("failed to commit branch: %w (output: %s)", err, string(out))
	}

	return nil
}

// AddToBranchMix adds a track to a branch mix
func (r *Repository) AddToBranchMix(branchName string, trackID string) error {
	// Ensure Git is configured properly
	if err := r.ensureGitConfig(); err != nil {
		return err
	}

	branchDir := filepath.Join(r.BranchMixPath, branchName)

	// Check if branch exists
	if _, err := os.Stat(branchDir); os.IsNotExist(err) {
		return fmt.Errorf("branch mix not found: %s", branchName)
	}

	// Find the track in project files
	var trackPath string
	var trackMetadata models.TrackMetadata

	// Search all instrument directories for the track
	instruments := []string{"drums", "vocals", "guitars", "bass", "keys", "fx"}
	for _, instrument := range instruments {
		instrumentDir := filepath.Join(r.ProjectFilesPath, instrument)

		err := filepath.Walk(instrumentDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Look for JSON metadata files
			if filepath.Ext(path) == ".json" {
				metadataBytes, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				var metadata models.TrackMetadata
				if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
					return nil
				}

				if metadata.ID == trackID {
					// Found the track - remove .json extension to get actual track path
					trackPath = path[:len(path)-5]
					trackMetadata = metadata
					return filepath.SkipDir // Stop searching
				}
			}

			return nil
		})

		if err != nil {
			continue
		}

		if trackPath != "" {
			break
		}
	}

	if trackPath == "" {
		return fmt.Errorf("track not found: %s", trackID)
	}

	// Copy track to branch directory
	targetPath := filepath.Join(branchDir, trackMetadata.FileName)

	input, err := os.ReadFile(trackPath)
	if err != nil {
		return fmt.Errorf("failed to read track file: %w", err)
	}

	if err := os.WriteFile(targetPath, input, 0644); err != nil {
		return fmt.Errorf("failed to write track to branch: %w", err)
	}

	// Update branch metadata
	metadataPath := filepath.Join(branchDir, "branch-info.json")
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read branch metadata: %w", err)
	}

	var branchMetadata models.BranchMetadata
	if err := json.Unmarshal(metadataBytes, &branchMetadata); err != nil {
		return fmt.Errorf("failed to unmarshal branch metadata: %w", err)
	}

	// Add track to branch metadata if not already present
	trackExists := false
	for _, id := range branchMetadata.Tracks {
		if id == trackID {
			trackExists = true
			break
		}
	}

	if !trackExists {
		branchMetadata.Tracks = append(branchMetadata.Tracks, trackID)
	}

	// Save updated metadata
	updatedMetadataJSON, err := json.MarshalIndent(branchMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, updatedMetadataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write updated metadata: %w", err)
	}

	// Add files to Git
	cmd := exec.Command("git", "add", targetPath, metadataPath)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add files to Git: %w (output: %s)", err, string(out))
	}

	// Create commit
	commitMsg := fmt.Sprintf("Add track %s to branch mix: %s", trackMetadata.FileName, branchName)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		// Check if it's a "nothing to commit" error, which is not a real error
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("failed to commit files: %w (output: %s)", err, string(out))
	}

	return nil
}

// PromoteToMain promotes a branch mix to the main mix
func (r *Repository) PromoteToMain(branchName string, userId string) error {
	// Ensure Git is configured properly
	if err := r.ensureGitConfig(); err != nil {
		return err
	}

	// Verify user is the mix master
	if userId != r.MixMaster {
		return fmt.Errorf("only the mix master can promote a branch mix to main")
	}

	branchDir := filepath.Join(r.BranchMixPath, branchName)

	// Check if branch exists
	if _, err := os.Stat(branchDir); os.IsNotExist(err) {
		return fmt.Errorf("branch mix not found: %s", branchName)
	}

	// Read branch metadata
	metadataPath := filepath.Join(branchDir, "branch-info.json")
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read branch metadata: %w", err)
	}

	var branchMetadata models.BranchMetadata
	if err := json.Unmarshal(metadataBytes, &branchMetadata); err != nil {
		return fmt.Errorf("failed to unmarshal branch metadata: %w", err)
	}

	// Create backup of current main mix if it exists
	if _, err := os.Stat(r.MainMixPath); err == nil {
		backupDir := fmt.Sprintf("%s-backup-%s", r.MainMixPath, time.Now().Format("20060102-150405"))
		if err := os.Rename(r.MainMixPath, backupDir); err != nil {
			return fmt.Errorf("failed to backup main mix: %w", err)
		}

		// Create new main mix directory
		if err := os.MkdirAll(r.MainMixPath, 0755); err != nil {
			return fmt.Errorf("failed to create main mix directory: %w", err)
		}
	} else {
		// Create main mix directory if it doesn't exist
		if err := os.MkdirAll(r.MainMixPath, 0755); err != nil {
			return fmt.Errorf("failed to create main mix directory: %w", err)
		}
	}

	// Copy all files from branch mix to main mix
	entries, err := os.ReadDir(branchDir)
	if err != nil {
		return fmt.Errorf("failed to read branch directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "branch-info.json" {
			continue
		}

		srcPath := filepath.Join(branchDir, entry.Name())
		dstPath := filepath.Join(r.MainMixPath, entry.Name())

		// Copy file
		input, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		if err := os.WriteFile(dstPath, input, 0644); err != nil {
			return fmt.Errorf("failed to write destination file: %w", err)
		}
	}

	// Create main mix metadata
	mainMixMetadata := models.MainMixMetadata{
		PromotedFrom: branchName,
		PromotedAt:   time.Now(),
		PromotedBy:   userId,
		Description:  branchMetadata.Description,
		Tracks:       branchMetadata.Tracks,
		Version:      1, // In future versions, we'd increment this
	}

	// Save main mix metadata
	mainMetadataPath := filepath.Join(r.MainMixPath, "main-mix-info.json")
	mainMetadataJSON, err := json.MarshalIndent(mainMixMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal main mix metadata: %w", err)
	}

	if err := os.WriteFile(mainMetadataPath, mainMetadataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write main mix metadata: %w", err)
	}

	// Add to Git
	cmd := exec.Command("git", "add", r.MainMixPath)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add main mix to Git: %w (output: %s)", err, string(out))
	}

	// Create commit
	commitMsg := fmt.Sprintf("Promote branch mix '%s' to main mix", branchName)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = r.RootPath
	if out, err := cmd.CombinedOutput(); err != nil {
		// Check if it's a "nothing to commit" error, which is not a real error
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("failed to commit main mix: %w (output: %s)", err, string(out))
	}

	return nil
}

// ListBranches lists all branch mixes
func (r *Repository) ListBranches() ([]models.BranchMetadata, error) {
	var branches []models.BranchMetadata

	entries, err := os.ReadDir(r.BranchMixPath)
	if err != nil {
		if os.IsNotExist(err) {
			return branches, nil
		}
		return nil, fmt.Errorf("failed to read branch directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadataPath := filepath.Join(r.BranchMixPath, entry.Name(), "branch-info.json")
		metadataBytes, err := os.ReadFile(metadataPath)
		if err != nil {
			continue
		}

		var metadata models.BranchMetadata
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			continue
		}

		branches = append(branches, metadata)
	}

	return branches, nil
}
