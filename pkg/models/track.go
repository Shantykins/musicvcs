// pkg/models/track.go
package models

import (
    "time"
)

// TrackMetadata stores information about an audio track
type TrackMetadata struct {
    ID          string    `json:"id"`
    FileName    string    `json:"fileName"`
    Type        string    `json:"type"` // vocals, drums, guitar, etc.
    CreatedAt   time.Time `json:"createdAt"`
    CreatedBy   string    `json:"createdBy"`
    Duration    float64   `json:"duration"` // in seconds
    SampleRate  int       `json:"sampleRate"`
    BitDepth    int       `json:"bitDepth"`
    Channels    int       `json:"channels"`
    Version     int       `json:"version"`
    Notes       string    `json:"notes"`
    Tags        []string  `json:"tags"`
    Dependencies []string  `json:"dependencies"` // IDs of tracks this depends on
}

// BranchMetadata stores information about a branch mix
type BranchMetadata struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
    CreatedBy   string    `json:"createdBy"`
    Tracks      []string  `json:"tracks"` // Track IDs included in this branch
}

// MainMixMetadata stores information about the main mix
type MainMixMetadata struct {
    PromotedFrom string    `json:"promotedFrom"` // Branch name
    PromotedAt   time.Time `json:"promotedAt"`
    PromotedBy   string    `json:"promotedBy"`
    Description  string    `json:"description"`
    Tracks       []string  `json:"tracks"` // Track IDs included in this mix
    Version      int       `json:"version"`
}