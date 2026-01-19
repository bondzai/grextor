package embed

import "context"

// Embedder defines the interface for generating vector embeddings from text.
type Embedder interface {
	// Embed generates a vector embedding for the given text.
	Embed(ctx context.Context, text string) ([]float32, error)
}
