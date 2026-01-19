package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jStore struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStore(uri, username, password string) (*Neo4jStore, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	return &Neo4jStore{driver: driver}, nil
}

func (s *Neo4jStore) Close(ctx context.Context) error {
	return s.driver.Close(ctx)
}

func (s *Neo4jStore) VerifyConnectivity(ctx context.Context) error {
	return s.driver.VerifyConnectivity(ctx)
}

func (s *Neo4jStore) AddNode(ctx context.Context, node *Node) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := fmt.Sprintf("MERGE (n:%s {id: $id}) SET n += $props", node.Label)
		params := map[string]interface{}{
			"id":    node.ID,
			"props": node.Properties,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	if err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}
	return nil
}

func (s *Neo4jStore) AddEdge(ctx context.Context, edge *Edge) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// Example: MATCH (a {id: $from}), (b {id: $to}) MERGE (a)-[r:TYPE]->(b) SET r += $props
		// Note: We need to know the labels of nodes a and b ideally, or just match by ID if ID is globally unique.
		// Assuming ID is globally unique for simplicity or we do an untyped match (slower but flexible).
		query := fmt.Sprintf(`
			MATCH (a {id: $from})
			MATCH (b {id: $to})
			MERGE (a)-[r:%s]->(b)
			SET r += $props
		`, edge.Type)

		params := map[string]interface{}{
			"from":  edge.FromID,
			"to":    edge.ToID,
			"props": edge.Properties,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	if err != nil {
		return fmt.Errorf("failed to add edge: %w", err)
	}
	return nil
}
