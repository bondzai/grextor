package vector

import "context"

// Point represents a data point in the vector store.
type Point struct {
	ID       string                 `json:"id"`
	Vector   []float32              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ScoredPoint represents a search result with a similarity score.
type ScoredPoint struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Store defines the interface for interacting with the vector database.
type Store interface {
	// Upsert stores or updates points in the vector database.
	Upsert(ctx context.Context, points []*Point) error
	// Search finds the nearest neighbors for the given vector.
	Search(ctx context.Context, vector []float32, limit int) ([]*ScoredPoint, error)
}
