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

### Basic Commands

#### Initialize a Project

```bash
./musicvcs init -path /path/to/project
```

This creates the directory structure and initializes Git and Git LFS.

#### Add a Track

```bash
./musicvcs add-track -file /path/to/audio.wav -type drums
```

This copies the audio file to the project, categorizes it by type, and tracks it with Git LFS. 
Valid track types: drums, vocals, guitars, bass, keys, fx.

#### Create a Branch Mix

```bash
./musicvcs create-branch -name alternative-mix -desc "Alternative arrangement with more bass"
```

Creates a new branch mix where you can assemble tracks in a different way.

#### Add a Track to a Branch Mix

```bash
./musicvcs add-to-branch -branch alternative-mix -track [track-id]
```

Adds the specified track to the branch mix. The track ID is output when you add a track.

#### List All Branch Mixes

```bash
./musicvcs list-branches
```

Shows all branch mixes with their descriptions and number of tracks.

#### Promote a Branch Mix to Main

```bash
./musicvcs promote -branch alternative-mix
```

Makes the specified branch mix the new main mix. Only the designated mix master can do this.

### Advanced Features

- **Mix Master Control**: Set the mix master when starting the tool: `./musicvcs -mix-master john.doe@example.com`
- **Multiple Branches**: Create as many branch mixes as needed to experiment with different arrangements
- **Git Integration**: All changes are tracked with Git, allowing for standard Git operations

## Use Cases

MusicVCS is ideal for:
- Solo producers who want to track versions of their mixes
- Collaborative music production teams
- Studios managing multiple versions of projects

The system handles large audio files efficiently using Git LFS and provides a clear structure for organizing music production projects. 