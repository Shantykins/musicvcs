// cmd/musicvcs/main.go
package main

import (
    "flag"
    "fmt"
    "os"
    
    "github.com/shantykins/musicvcs/internal/repo"
)

func main() {
    var (
        rootPath  string
        mixMaster string
    )
    
    // Define flags
    flag.StringVar(&rootPath, "path", ".", "Path to the music project")
    flag.StringVar(&mixMaster, "mix-master", os.Getenv("USER"), "User ID of the mix master")
    
    // Define command flags
    initCmd := flag.NewFlagSet("init", flag.ExitOnError)
    
    addTrackCmd := flag.NewFlagSet("add-track", flag.ExitOnError)
    addTrackFile := addTrackCmd.String("file", "", "Path to the audio file")
    addTrackType := addTrackCmd.String("type", "", "Type of track (drums, vocals, guitars, bass, keys, fx)")
    
    createBranchCmd := flag.NewFlagSet("create-branch", flag.ExitOnError)
    createBranchName := createBranchCmd.String("name", "", "Name of the branch mix")
    createBranchDesc := createBranchCmd.String("desc", "", "Description of the branch mix")
    
    addToBranchCmd := flag.NewFlagSet("add-to-branch", flag.ExitOnError)
    addToBranchName := addToBranchCmd.String("branch", "", "Name of the branch mix")
    addToBranchTrack := addToBranchCmd.String("track", "", "ID of the track")
    
    promoteCmd := flag.NewFlagSet("promote", flag.ExitOnError)
    promoteBranch := promoteCmd.String("branch", "", "Name of the branch mix to promote")
    
    listBranchesCmd := flag.NewFlagSet("list-branches", flag.ExitOnError)
    
    // Parse top-level flags
    flag.Parse()
    
    // Check if a command was provided
    if flag.NArg() < 1 {
        fmt.Println("Command required")
        fmt.Println("Available commands: init, add-track, create-branch, add-to-branch, promote, list-branches")
        os.Exit(1)
    }
    
    // Get command
    command := flag.Arg(0)
    
    // Create repository
    repository := repo.New(rootPath, mixMaster)
    
    // Handle commands
    switch command {
    case "init":
        initCmd.Parse(flag.Args()[1:])
        if err := repository.Initialize(); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Music version control system initialized successfully!")
        
    case "add-track":
        addTrackCmd.Parse(flag.Args()[1:])
        if *addTrackFile == "" {
            fmt.Println("Error: File path required")
            os.Exit(1)
        }
        if *addTrackType == "" {
            fmt.Println("Error: Track type required")
            os.Exit(1)
        }
        
        trackID, err := repository.AddTrack(*addTrackFile, *addTrackType, mixMaster)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Track added successfully with ID: %s\n", trackID)
        
    case "create-branch":
        createBranchCmd.Parse(flag.Args()[1:])
        if *createBranchName == "" {
            fmt.Println("Error: Branch name required")
            os.Exit(1)
        }
        
        if err := repository.CreateBranchMix(*createBranchName, *createBranchDesc, mixMaster); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Branch mix '%s' created successfully\n", *createBranchName)
        
    case "add-to-branch":
        addToBranchCmd.Parse(flag.Args()[1:])
        if *addToBranchName == "" {
            fmt.Println("Error: Branch name required")
            os.Exit(1)
        }
        if *addToBranchTrack == "" {
            fmt.Println("Error: Track ID required")
            os.Exit(1)
        }
        
        if err := repository.AddToBranchMix(*addToBranchName, *addToBranchTrack); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Track added to branch mix '%s' successfully\n", *addToBranchName)
        
    case "promote":
        promoteCmd.Parse(flag.Args()[1:])
        if *promoteBranch == "" {
            fmt.Println("Error: Branch name required")
            os.Exit(1)
        }
        
        if err := repository.PromoteToMain(*promoteBranch, mixMaster); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Branch mix '%s' promoted to main mix successfully\n", *promoteBranch)
        
    case "list-branches":
        listBranchesCmd.Parse(flag.Args()[1:])
        
        branches, err := repository.ListBranches()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        
        if len(branches) == 0 {
            fmt.Println("No branch mixes found")
        } else {
            fmt.Println("Branch mixes:")
            for _, branch := range branches {
                fmt.Printf("- %s: %s (%d tracks)\n", 
                    branch.Name, branch.Description, len(branch.Tracks))
            }
        }
        
    default:
        fmt.Printf("Unknown command: %s\n", command)
        fmt.Println("Available commands: init, add-track, create-branch, add-to-branch, promote, list-branches")
        os.Exit(1)
    }
}