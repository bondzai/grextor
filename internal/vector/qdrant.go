package vector

import (
	"context"
	"fmt"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QdrantStore struct {
	conn           *grpc.ClientConn
	pointsClient   pb.PointsClient
	collectionName string
	vectorSize     uint64
}

func NewQdrantStore(addr string, collectionName string, vectorSize uint64) (*QdrantStore, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %w", err)
	}

	client := pb.NewPointsClient(conn)

	// Note: We assume the collection might need to be created elsewhere or we can add an EnsureCollection method.
	// For simplicity, we just return the store.

	return &QdrantStore{
		conn:           conn,
		pointsClient:   client,
		collectionName: collectionName,
		vectorSize:     vectorSize,
	}, nil
}

func (s *QdrantStore) Close() error {
	return s.conn.Close()
}

// EnsureCollection creates the collection if it doesn't exist.
func (s *QdrantStore) EnsureCollection(ctx context.Context) error {
	collectionsClient := pb.NewCollectionsClient(s.conn)

	// Check if exists
	exists, err := collectionsClient.Get(ctx, &pb.GetCollectionInfoRequest{CollectionName: s.collectionName})
	if err == nil && exists != nil {
		return nil
	}

	// Create
	_, err = collectionsClient.Create(ctx, &pb.CreateCollection{
		CollectionName: s.collectionName,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     s.vectorSize,
					Distance: pb.Distance_Cosine,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}

func (s *QdrantStore) Upsert(ctx context.Context, points []*Point) error {
	qPoints := make([]*pb.PointStruct, len(points))
	for i, p := range points {
		// Convert metadata map to Qdrant payload
		payload := make(map[string]*pb.Value)
		for k, v := range p.Metadata {
			payload[k] = toPbValue(v)
		}

		qPoints[i] = &pb.PointStruct{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Uuid{Uuid: p.ID},
			},
			Vectors: &pb.Vectors{
				VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: p.Vector}},
			},
			Payload: payload,
		}
	}

	_, err := s.pointsClient.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: s.collectionName,
		Points:         qPoints,
	})
	return err
}

func (s *QdrantStore) Search(ctx context.Context, vector []float32, limit int) ([]*ScoredPoint, error) {
	res, err := s.pointsClient.Search(ctx, &pb.SearchPoints{
		CollectionName: s.collectionName,
		Vector:         vector,
		Limit:          uint64(limit),
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}

	results := make([]*ScoredPoint, len(res.Result))
	for i, r := range res.Result {
		meta := make(map[string]interface{})
		for k, v := range r.Payload {
			meta[k] = fromPbValue(v)
		}

		var id string
		if r.Id != nil {
			if u := r.Id.GetUuid(); u != "" {
				id = u
			} else {
				id = fmt.Sprintf("%d", r.Id.GetNum())
			}
		}

		results[i] = &ScoredPoint{
			ID:       id,
			Score:    r.Score,
			Metadata: meta,
		}
	}
	return results, nil
}

// Helper to convert Go interface{} to Qdrant Value
func toPbValue(v interface{}) *pb.Value {
	switch val := v.(type) {
	case string:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: val}}
	case int:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: int64(val)}}
	case int64:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: val}}
	case float64:
		return &pb.Value{Kind: &pb.Value_DoubleValue{DoubleValue: val}}
	case bool:
		return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: val}}
	default:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: fmt.Sprintf("%v", val)}}
	}
}

// Helper to convert Qdrant Value to Go interface{}
func fromPbValue(v *pb.Value) interface{} {
	switch k := v.Kind.(type) {
	case *pb.Value_StringValue:
		return k.StringValue
	case *pb.Value_IntegerValue:
		return k.IntegerValue
	case *pb.Value_DoubleValue:
		return k.DoubleValue
	case *pb.Value_BoolValue:
		return k.BoolValue
	default:
		return nil
	}
}
