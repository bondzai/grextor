package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bondzai/grextor/internal/embed"
	"github.com/bondzai/grextor/internal/engine"
	"github.com/bondzai/grextor/internal/graph"
	"github.com/bondzai/grextor/internal/vector"
	"github.com/google/uuid"
)

func main() {
	var (
		qdrantAddr = flag.String("qdrant-addr", "localhost:6334", "Qdrant gRPC address")
		neo4jURI   = flag.String("neo4j-uri", "bolt://localhost:7687", "Neo4j URI")
		neo4jUser  = flag.String("neo4j-user", "neo4j", "Neo4j username")
		neo4jPass  = flag.String("neo4j-pass", "grextor123", "Neo4j password")
		collection = flag.String("collection", "grextor_docs", "Qdrant collection name")
		content    = flag.String("content", "", "Content to ingest")
		docID      = flag.String("id", "", "Document ID (optional, generated if empty)")
	)
	flag.Parse()

	if *content == "" {
		log.Fatal("Please provide content to ingest using --content or pass a file (not yet impl)")
	}
	if *docID == "" {
		*docID = uuid.New().String()
	}

	ctx := context.Background()

	// 1. Setup Embedder
	var embedder embed.Embedder
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		log.Println("Using OpenAI Embedder")
		embedder = embed.NewOpenAIEmbedder(apiKey, "")
	} else {
		log.Println("Using NoOp Embedder (Dummy)")
		embedder = embed.NewNoOpEmbedder(1536)
	}

	// 2. Setup Vector Store (Qdrant)
	vStore, err := vector.NewQdrantStore(*qdrantAddr, *collection, 1536)
	if err != nil {
		log.Fatalf("Failed to connect to Qdrant: %v", err)
	}
	defer vStore.Close()

	if err := vStore.EnsureCollection(ctx); err != nil {
		log.Fatalf("Failed to ensure collection: %v", err)
	}

	// 3. Setup Graph Store (Neo4j)
	gStore, err := graph.NewNeo4jStore(*neo4jURI, *neo4jUser, *neo4jPass)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer gStore.Close(ctx)

	// Verify Neo4j
	if err := gStore.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("Failed to verify Neo4j connectivity: %v. Make sure Docker is running.", err)
	}

	// 4. Initialize Engine
	eng := engine.NewEngine(embedder, vStore, gStore)

	// 5. Ingest
	start := time.Now()
	err = eng.IngestDocument(ctx, *docID, *content, map[string]interface{}{
		"source": "cli",
		"time":   time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Fatalf("Ingestion failed: %v", err)
	}

	fmt.Printf("Ingestion successful! ID: %s (took %v)\n", *docID, time.Since(start))
}
