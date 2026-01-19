package graph

import "context"

// Node represents a node in the property graph.
type Node struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

// Edge represents a relationship between two nodes.
type Edge struct {
	FromID     string                 `json:"from_id"`
	ToID       string                 `json:"to_id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Store defines the interface for interacting with the graph database.
type Store interface {
	// AddNode adds or updates a node in the graph.
	AddNode(ctx context.Context, node *Node) error
	// AddEdge adds or updates an edge between two nodes.
	AddEdge(ctx context.Context, edge *Edge) error
}
