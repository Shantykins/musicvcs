# MusicVCS - Music Version Control System

MusicVCS is a specialized version control system designed for music production projects. It uses Git and Git LFS under the hood to manage audio files and track changes to music projects.

## Core Features

1. **Track Management**: Add and organize audio tracks by instrument type (drums, vocals, guitars, bass, keys, fx)
2. **Branch Mixes**: Create alternative mixes as branches to experiment with different arrangements
3. **Main Mix Promotion**: Promote a branch mix to become the main mix
4. **Mix Master Role**: Only designated mix masters can promote branch mixes to main
5. **Audio File Handling**: Uses Git LFS to efficiently manage large audio files
6. **Metadata Tracking**: Stores metadata for tracks and mixes

## How It Works

### Architecture

MusicVCS follows a layered architecture:

1. **CLI Layer** (`cmd/musicvcs/main.go`): Command-line interface for user interaction
2. **Repository Layer** (`internal/repo/repository.go`): Core functionality for managing the music project
3. **LFS Manager** (`internal/lfs/manager.go`): Handles Git LFS operations for large audio files
4. **Models** (`pkg/models/track.go`): Data structures for tracks and mixes

### Directory Structure

When initialized, MusicVCS creates the following structure:

```
/project-root
├── main-mix/            # The current main mix
├── branch-mix/          # Alternative mixes
│   └── [branch-name]/
│       ├── [audio-files]
│       └── branch-info.json
└── project-files/       # Original track files
    ├── drums/
    ├── vocals/
    ├── guitars/
    ├── bass/
    ├── keys/
    └── fx/
```

### Workflow

1. **Initialize**: Creates the directory structure and git repository
2. **Add Tracks**: Add audio files to the project, organized by instrument type
3. **Create Branch Mixes**: Make alternative arrangements/mixes as branches
4. **Add Tracks to Branches**: Assemble different tracks into branch mixes
5. **Promote to Main**: When a branch mix is ready, promote it to become the main mix

## How to Use MusicVCS

### Installation

```bash
git clone https://github.com/shantykins/musicvcs.git
cd musicvcs
go build -o musicvcs ./cmd/musicvcs
```

To make musicvcs available system-wide (like git), copy the executable to a directory in your PATH:

```bash
sudo cp musicvcs /usr/local/bin/
```

### Basic Commands

#### Initialize a Project

```bash
musicvcs -path=/path/to/project init
```

This creates the directory structure and initializes Git and Git LFS.

#### Add a Track

```bash
musicvcs -path=/path/to/project add-track -file=/path/to/audio.wav -type=drums
```

This copies the audio file to the project, categorizes it by type, and tracks it with Git LFS. 
Valid track types: drums, vocals, guitars, bass, keys, fx.

#### Create a Branch Mix

```bash
musicvcs -path=/path/to/project create-branch -name=alternative-mix -desc="Alternative arrangement with more bass"
```

Creates a new branch mix where you can assemble tracks in a different way.

#### Add a Track to a Branch Mix

```bash
musicvcs -path=/path/to/project add-to-branch -branch=alternative-mix -track=[track-id]
```

Adds the specified track to the branch mix. The track ID is output when you add a track.

#### List All Branch Mixes

```bash
musicvcs -path=/path/to/project list-branches
```

Shows all branch mixes with their descriptions and number of tracks.

#### Promote a Branch Mix to Main

```bash
musicvcs -path=/path/to/project promote -branch=alternative-mix
```

Makes the specified branch mix the new main mix. Only the designated mix master can do this.

### Advanced Features

- **Mix Master Control**: Set the mix master when starting the tool: `musicvcs -mix-master john.doe@example.com`
- **Multiple Branches**: Create as many branch mixes as needed to experiment with different arrangements
- **Git Integration**: All changes are tracked with Git, allowing for standard Git operations

## Sample Workflow: Vocal Recording Project

Here's a detailed workflow example showing how to use MusicVCS for a vocal recording project:

### 1. Setup Project and Add Tracks

```bash
# Initialize a new project
mkdir my_project
cd my_project
musicvcs -path=. init

# Add vocal tracks
musicvcs -path=. add-track -file=/path/to/Lead_Vocals.wav -type=vocals
# Note the returned track ID: b4f85868-2d3b-49cf-acc8-78aa0d497ecf

musicvcs -path=. add-track -file=/path/to/Harmony1.wav -type=vocals
# Note the returned track ID: 785f171e-f572-4675-92fa-cc28e9d4890c

musicvcs -path=. add-track -file=/path/to/Harmony2.wav -type=vocals
# Note the returned track ID: 52438fd6-0824-4351-8597-34f842496cf4

musicvcs -path=. add-track -file=/path/to/Harmony3.wav -type=vocals
# Note the returned track ID: c2e5b17e-718b-49a0-b917-8da0ed20a9cc
```

### 2. Create and Set Up Different Mix Branches

```bash
# Create a main mix branch
musicvcs -path=. create-branch -name=mix-master -desc="Master mix for song_title"

# Create a vocals-focused mix branch
musicvcs -path=. create-branch -name=vocals -desc="Vocals mix for song_title"

# Add tracks to the vocals branch
musicvcs -path=. add-to-branch -branch=vocals -track=b4f85868-2d3b-49cf-acc8-78aa0d497ecf
musicvcs -path=. add-to-branch -branch=vocals -track=785f171e-f572-4675-92fa-cc28e9d4890c
musicvcs -path=. add-to-branch -branch=vocals -track=52438fd6-0824-4351-8597-34f842496cf4
musicvcs -path=. add-to-branch -branch=vocals -track=c2e5b17e-718b-49a0-b917-8da0ed20a9cc

# Add the same tracks to the mix-master branch
musicvcs -path=. add-to-branch -branch=mix-master -track=b4f85868-2d3b-49cf-acc8-78aa0d497ecf
musicvcs -path=. add-to-branch -branch=mix-master -track=785f171e-f572-4675-92fa-cc28e9d4890c
musicvcs -path=. add-to-branch -branch=mix-master -track=52438fd6-0824-4351-8597-34f842496cf4
musicvcs -path=. add-to-branch -branch=mix-master -track=c2e5b17e-718b-49a0-b917-8da0ed20a9cc

# List all branches to verify
musicvcs -path=. list-branches
```

### 3. Make Different Versions in Different Branches

```bash
# Add processing notes to vocals branch
echo "Original vocal notes - modified with extra reverb and delay" > branch-mix/vocals/Lead_Vocals_notes.txt
git add branch-mix/vocals/Lead_Vocals_notes.txt
git commit -m "Add reverb and delay processing notes to vocals branch"

# Add different processing notes to mix-master branch
echo "Standard lead vocals with minimal processing" > branch-mix/mix-master/Lead_Vocals_notes.txt
git add branch-mix/mix-master/Lead_Vocals_notes.txt
git commit -m "Add standard processing notes to mix-master branch"
```

### 4. Promote Branch Mix to Main

```bash
# Promote the mix-master branch to main
musicvcs -path=. promote -branch=mix-master

# Verify main mix contents
cat main-mix/main-mix-info.json
```

### 5. Switch Between Versions and Branches

To switch to a specific version of a branch (going back in time):

```bash
# List commits to find the commit hash
git log --oneline -- branch-mix/vocals/

# Create a temporary branch at that point in time
git checkout -b temp-version 6f21783  # Replace with your specific commit hash

# Now you can view the files at that point in time
ls -la branch-mix/vocals/
```

To return to the latest version:

```bash
# Switch back to main branch with latest changes
git checkout main

# Verify the current state
ls -la branch-mix/vocals/
cat branch-mix/vocals/Lead_Vocals_notes.txt
```

### 6. Making Updates and Comparing Versions

After making changes to a branch, you can compare different versions:

```bash
# Make changes to the vocals branch
echo "Updated vocal processing: more reverb, less delay" > branch-mix/vocals/Lead_Vocals_notes.txt
git add branch-mix/vocals/Lead_Vocals_notes.txt
git commit -m "Update vocal processing settings"

# Compare different branches
diff branch-mix/vocals/Lead_Vocals_notes.txt branch-mix/mix-master/Lead_Vocals_notes.txt

# Compare with previous version of the same file
git diff HEAD~1 HEAD -- branch-mix/vocals/Lead_Vocals_notes.txt
```

This workflow demonstrates how MusicVCS gives you the power to manage different versions of your music project, try alternative processing approaches, and maintain a complete history of your work.

## Use Cases

MusicVCS is ideal for:
- Solo producers who want to track versions of their mixes
- Collaborative music production teams
- Studios managing multiple versions of projects

The system handles large audio files efficiently using Git LFS and provides a clear structure for organizing music production projects. 