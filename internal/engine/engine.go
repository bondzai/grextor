package engine

import (
	"context"
	"fmt"
	"log"

	"github.com/bondzai/grextor/internal/embed"
	"github.com/bondzai/grextor/internal/graph"
	"github.com/bondzai/grextor/internal/vector"
)

type Engine struct {
	embedder    embed.Embedder
	vectorStore vector.Store
	graphStore  graph.Store
}

func NewEngine(e embed.Embedder, v vector.Store, g graph.Store) *Engine {
	return &Engine{
		embedder:    e,
		vectorStore: v,
		graphStore:  g,
	}
}

// IngestDocument processes a document: embeds it, stores in vector DB, and creates a node in graph DB.
func (e *Engine) IngestDocument(ctx context.Context, id, content string, metadata map[string]interface{}) error {
	log.Printf("Ingesting document %s...", id)

	// 1. Generate Embedding
	vec, err := e.embedder.Embed(ctx, content)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}

	// 2. Store in Vector DB
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["content"] = content // Store content in metadata for retrieval

	err = e.vectorStore.Upsert(ctx, []*vector.Point{
		{
			ID:       id,
			Vector:   vec,
			Metadata: metadata,
		},
	})
	if err != nil {
		return fmt.Errorf("vector storage failed: %w", err)
	}

	// 3. Store in Graph DB (Node)
	node := &graph.Node{
		ID:         id,
		Label:      "Document",
		Properties: metadata,
	}
	err = e.graphStore.AddNode(ctx, node)
	if err != nil {
		return fmt.Errorf("graph storage failed: %w", err)
	}

	log.Printf("Successfully ingested document %s", id)
	return nil
}

// SearchResult combines vector score and metadata.
type SearchResult struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (e *Engine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	log.Printf("Searching for: %s", query)

	// 1. Embed Query
	vec, err := e.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query embedding failed: %w", err)
	}

	// 2. Vector Search
	scoredPoints, err := e.vectorStore.Search(ctx, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 3. Map Results
	results := make([]SearchResult, len(scoredPoints))
	for i, sp := range scoredPoints {
		content, _ := sp.Metadata["content"].(string)
		results[i] = SearchResult{
			ID:       sp.ID,
			Score:    sp.Score,
			Content:  content,
			Metadata: sp.Metadata,
		}
	}

	return results, nil
}
