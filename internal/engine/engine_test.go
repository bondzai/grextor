package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/bondzai/grextor/internal/graph"
	"github.com/bondzai/grextor/internal/vector"
)

func TestEngine_IngestDocument(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{}
		mockVectorStore := &MockVectorStore{
			UpsertFunc: func(ctx context.Context, points []*vector.Point) error {
				if len(points) != 1 {
					t.Errorf("expected 1 point, got %d", len(points))
				}
				if points[0].ID != "123" {
					t.Errorf("expected ID '123', got %s", points[0].ID)
				}
				return nil
			},
		}
		mockGraphStore := &MockGraphStore{
			AddNodeFunc: func(ctx context.Context, node *graph.Node) error {
				if node.ID != "123" {
					t.Errorf("expected node ID '123', got %s", node.ID)
				}
				return nil
			},
		}

		eng := NewEngine(mockEmbedder, mockVectorStore, mockGraphStore)
		err := eng.IngestDocument(ctx, "123", "test content", nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("EmbeddingError", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{
			EmbedFunc: func(ctx context.Context, text string) ([]float32, error) {
				return nil, errors.New("embed error")
			},
		}
		eng := NewEngine(mockEmbedder, &MockVectorStore{}, &MockGraphStore{})
		err := eng.IngestDocument(ctx, "123", "test content", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("VectorStoreError", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{}
		mockVectorStore := &MockVectorStore{
			UpsertFunc: func(ctx context.Context, points []*vector.Point) error {
				return errors.New("vector store error")
			},
		}
		eng := NewEngine(mockEmbedder, mockVectorStore, &MockGraphStore{})
		err := eng.IngestDocument(ctx, "123", "test content", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("GraphStoreError", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{}
		mockVectorStore := &MockVectorStore{}
		mockGraphStore := &MockGraphStore{
			AddNodeFunc: func(ctx context.Context, node *graph.Node) error {
				return errors.New("graph store error")
			},
		}
		eng := NewEngine(mockEmbedder, mockVectorStore, mockGraphStore)
		err := eng.IngestDocument(ctx, "123", "test content", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestEngine_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{}
		mockVectorStore := &MockVectorStore{
			SearchFunc: func(ctx context.Context, vec []float32, limit int) ([]*vector.ScoredPoint, error) {
				return []*vector.ScoredPoint{
					{
						ID:    "123",
						Score: 0.9,
						Metadata: map[string]interface{}{
							"content": "test content",
						},
					},
				}, nil
			},
		}
		eng := NewEngine(mockEmbedder, mockVectorStore, &MockGraphStore{})

		results, err := eng.Search(ctx, "query", 10)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].ID != "123" {
			t.Errorf("expected result ID '123', got %s", results[0].ID)
		}
	})

	t.Run("EmbeddingError", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{
			EmbedFunc: func(ctx context.Context, text string) ([]float32, error) {
				return nil, errors.New("embed error")
			},
		}
		eng := NewEngine(mockEmbedder, &MockVectorStore{}, &MockGraphStore{})
		_, err := eng.Search(ctx, "query", 10)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("VectorSearchError", func(t *testing.T) {
		mockEmbedder := &MockEmbedder{}
		mockVectorStore := &MockVectorStore{
			SearchFunc: func(ctx context.Context, vec []float32, limit int) ([]*vector.ScoredPoint, error) {
				return nil, errors.New("search error")
			},
		}
		eng := NewEngine(mockEmbedder, mockVectorStore, &MockGraphStore{})
		_, err := eng.Search(ctx, "query", 10)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
