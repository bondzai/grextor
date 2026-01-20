package engine

import (
	"context"

	"github.com/bondzai/grextor/internal/graph"
	"github.com/bondzai/grextor/internal/vector"
)

// MockEmbedder implements embed.Embedder
type MockEmbedder struct {
	EmbedFunc func(ctx context.Context, text string) ([]float32, error)
}

func (m *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if m.EmbedFunc != nil {
		return m.EmbedFunc(ctx, text)
	}
	return []float32{0.1, 0.2, 0.3}, nil
}

// MockVectorStore implements vector.Store
type MockVectorStore struct {
	UpsertFunc func(ctx context.Context, points []*vector.Point) error
	SearchFunc func(ctx context.Context, vec []float32, limit int) ([]*vector.ScoredPoint, error)
}

func (m *MockVectorStore) Upsert(ctx context.Context, points []*vector.Point) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, points)
	}
	return nil
}

func (m *MockVectorStore) Search(ctx context.Context, vec []float32, limit int) ([]*vector.ScoredPoint, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, vec, limit)
	}
	return nil, nil
}

// MockGraphStore implements graph.Store
type MockGraphStore struct {
	AddNodeFunc func(ctx context.Context, node *graph.Node) error
	AddEdgeFunc func(ctx context.Context, edge *graph.Edge) error
}

func (m *MockGraphStore) AddNode(ctx context.Context, node *graph.Node) error {
	if m.AddNodeFunc != nil {
		return m.AddNodeFunc(ctx, node)
	}
	return nil
}

func (m *MockGraphStore) AddEdge(ctx context.Context, edge *graph.Edge) error {
	if m.AddEdgeFunc != nil {
		return m.AddEdgeFunc(ctx, edge)
	}
	return nil
}
