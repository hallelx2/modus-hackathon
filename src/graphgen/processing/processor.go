package processing


import (
	"fmt"
	"time"
	"my-modus-app/graphgen/models"
	"my-modus-app/graphgen/embedding"
	"my-modus-app/graphgen/chunking"
	"context"
)

// Initialize embedding configuration
func NewEmbeddingConfig() embedding.EmbeddingConfig {
	return embedding.EmbeddingConfig{
		ModelName:    "embeddings", // Adjust as needed
		BatchSize:    10,
		MaxRetries:   3,
		CacheEnabled: true,
		CacheTTL:     24 * time.Hour, // Set TTL for cached embeddings (24 hours)
	}
}

func ChunkAndEmbed(text string, useAI bool) ([]models.TextChunk, error) {
	// Initialize the embedding configuration
	embeddingConfig := NewEmbeddingConfig()

	// Initialize the Embedder
	embedder, err := embedding.NewEmbedder(embeddingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize embedder: %w", err)
	}

	// Chunk the text using the chunker
	chunks, err := chunking.ChoiceChunker(text, useAI)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk the text: %w", err)
	}

	// Embed each chunk using the embedder
	for i, chunk := range chunks {
		embedding, err := embedder.EmbedText(context.Background(), chunk.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to embed text for chunk ID %s: %w", chunk.ID, err)
		}
		// Assign the embedding to the chunk
		chunks[i].Embedding = embedding
	}

	// Return the chunks with embeddings
	return chunks, nil
}
