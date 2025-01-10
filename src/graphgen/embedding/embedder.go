package embedding

import (
	"context"
	"fmt"
	"time"

	functionsModels "github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	graphgenModels "my-modus-app/graphgen/models"
)

type EmbeddingConfig struct {
	ModelName    string        `json:"model_name"`
	BatchSize    int           `json:"batch_size"`
	MaxRetries   int           `json:"max_retries"`
	CacheEnabled bool          `json:"cache_enabled"`
	CacheTTL     time.Duration `json:"cache_ttl"` // Time-to-live for cached embeddings
}

type Embedder struct {
	config EmbeddingConfig
	model  *openai.EmbeddingsModel
	cache  *EmbeddingCache
}

func NewEmbedder(config EmbeddingConfig) (*Embedder, error) {
	model, err := functionsModels.GetModel[openai.EmbeddingsModel](config.ModelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI embeddings model: %w", err)
	}

	cache := NewEmbeddingCache()

	return &Embedder{
		config: config,
		model:  model,
		cache:  cache,
	}, nil
}

// EmbedText handles embedding for a single text and applies caching
func (e *Embedder) EmbedText(ctx context.Context, text string) ([]float32, error) {
	if e.config.CacheEnabled {
		if embedding, found := e.cache.Get(text); found {
			return embedding, nil
		}
	}

	embeddings, err := e.GetEmbeddings(ctx, text)
	if err != nil {
		return nil, err
	}

	if e.config.CacheEnabled {
		e.cache.Set(text, embeddings, time.Now().Add(e.config.CacheTTL))
	}

	return embeddings, nil
}

// EmbedChunks embeds a list of chunks, with optional caching
func (e *Embedder) EmbedChunks(ctx context.Context, chunks []graphgenModels.TextChunk) error {
	for _, chunk := range chunks {
		_, err := e.EmbedText(ctx, chunk.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetEmbeddings retrieves embeddings from the OpenAI model
func (e *Embedder) GetEmbeddings(ctx context.Context, texts ...string) ([]float32, error) {
	input, err := e.model.CreateInput(texts)
	if err != nil {
		return nil, fmt.Errorf("failed to create input for embeddings: %w", err)
	}

	output, err := e.model.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke OpenAI model: %w", err)
	}

	if len(output.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return output.Data[0].Embedding, nil
}



// EmbedTextChunks processes a list of TextChunks and embeds their content
func (e *Embedder) EmbedTextChunks(ctx context.Context, chunks []graphgenModels.TextChunk) ([]graphgenModels.TextChunk, error) {
	for i, chunk := range chunks {
		// Embed the text content of the chunk
		embedding, err := e.EmbedText(ctx, chunk.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to embed text for chunk ID %s: %w", chunk.ID, err)
		}

		// Save the embedding into the chunk
		chunks[i].Embedding = embedding
	}

	return chunks, nil
}
