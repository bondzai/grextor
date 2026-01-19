package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bondzai/grextor/internal/embed"
	"github.com/bondzai/grextor/internal/engine"
	"github.com/bondzai/grextor/internal/graph"
	"github.com/bondzai/grextor/internal/vector"
)

func main() {
	var (
		qdrantAddr = flag.String("qdrant-addr", "localhost:6334", "Qdrant gRPC address")
		neo4jURI   = flag.String("neo4j-uri", "bolt://localhost:7687", "Neo4j URI")
		neo4jUser  = flag.String("neo4j-user", "neo4j", "Neo4j username")
		neo4jPass  = flag.String("neo4j-pass", "grextor123", "Neo4j password")
		collection = flag.String("collection", "grextor_docs", "Qdrant collection name")
		query      = flag.String("q", "", "Query text")
		limit      = flag.Int("limit", 5, "Number of results")
	)
	flag.Parse()

	if *query == "" {
		log.Fatal("Please provide a query using -q")
	}

	ctx := context.Background()

	// 1. Setup Embedder
	var embedder embed.Embedder
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		// log.Println("Using OpenAI Embedder")
		embedder = embed.NewOpenAIEmbedder(apiKey, "")
	} else {
		// log.Println("Using NoOp Embedder (Dummy)")
		embedder = embed.NewNoOpEmbedder(1536)
	}

	// 2. Setup Vector Store (Qdrant)
	vStore, err := vector.NewQdrantStore(*qdrantAddr, *collection, 1536)
	if err != nil {
		log.Fatalf("Failed to connect to Qdrant: %v", err)
	}
	defer vStore.Close()

	// 3. Setup Graph Store (Neo4j)
	gStore, err := graph.NewNeo4jStore(*neo4jURI, *neo4jUser, *neo4jPass)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer gStore.Close(ctx)

	// 4. Initialize Engine
	eng := engine.NewEngine(embedder, vStore, gStore)

	// 5. Search
	results, err := eng.Search(ctx, *query, *limit)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	fmt.Printf("Found %d results for '%s':\n", len(results), *query)
	for i, res := range results {
		fmt.Printf("%d. [Score: %.4f] %s\n   Content: %s\n", i+1, res.Score, res.ID, res.Content)
	}
}
