package embed

import "context"

type NoOpEmbedder struct {
	Dimensions int
}

func NewNoOpEmbedder(dims int) *NoOpEmbedder {
	if dims <= 0 {
		dims = 1536 // Default to Ada-002 size
	}
	return &NoOpEmbedder{Dimensions: dims}
}

func (e *NoOpEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	// Return a zero vector of the specified dimension
	return make([]float32, e.Dimensions), nil
}
