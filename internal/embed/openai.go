package embed

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIEmbedder struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

func NewOpenAIEmbedder(apiKey string, model openai.EmbeddingModel) *OpenAIEmbedder {
	if model == "" {
		model = openai.AdaEmbeddingV2
	}
	return &OpenAIEmbedder{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	text = strings.ReplaceAll(text, "\n", " ")
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: e.model,
	}

	resp, err := e.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("creating embeddings: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
