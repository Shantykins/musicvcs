// internal/p2p/network.go
package p2p

// Node represents a peer in the network
type Node struct {
    ID      string
    Address string
}

// Network defines the P2P communication interface
type Network interface {
    // Connect establishes connection to a peer
    Connect(address string) (*Node, error)
    
    // Disconnect closes connection to a peer
    Disconnect(nodeID string) error
    
    // ListPeers returns all connected peers
    ListPeers() ([]*Node, error)
    
    // FetchMetadata retrieves a peer's branch information
    FetchMetadata(nodeID string) (map[string]interface{}, error)
    
    // PullChanges gets updates from a peer's branch
    PullChanges(nodeID string, branchName string) error
    
    // PushChanges sends updates to a peer
    PushChanges(nodeID string, branchName string) error
}

// This would be expanded in future implementations