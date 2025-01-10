// File: processing/pipeline.go
package processing

import (
	"context"
	"fmt"
	"log"
	"time"

	"my-modus-app/graphgen/chunking"
	"my-modus-app/graphgen/embedding"
	"my-modus-app/graphgen/models"
)

// PipelineConfig holds all configuration needed for the processing pipeline
type PipelineConfig struct {
	// Chunking configuration
	UseAIChunking bool

	// Embedding configuration
	EmbeddingConfig embedding.EmbeddingConfig

	// Graph generation configuration
	GraphConfig Config
}

// DefaultPipelineConfig returns a default configuration for the pipeline
func DefaultPipelineConfig() PipelineConfig {
	return PipelineConfig{
		UseAIChunking: false,
		EmbeddingConfig: embedding.EmbeddingConfig{
			ModelName:    "embeddings", // Default OpenAI embedding model
			BatchSize:    10,
			MaxRetries:   3,
			CacheEnabled: true,
			CacheTTL:     time.Hour * 24,
		},
		GraphConfig: DefaultConfig(),
	}
}

// Pipeline represents the text processing pipeline
type Pipeline struct {
	config     PipelineConfig
	embedder   *embedding.Embedder
	comparator *ChunkComparator
}

// NewPipeline creates a new processing pipeline with the given configuration
func NewPipeline(config PipelineConfig, llmClient LLMClient) (*Pipeline, error) {
	// Initialize embedder
	embedder, err := embedding.NewEmbedder(config.EmbeddingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize embedder: %w", err)
	}

	// Initialize comparator
	comparator := NewChunkComparator(llmClient, config.GraphConfig)

	return &Pipeline{
		config:     config,
		embedder:   embedder,
		comparator: comparator,
	}, nil
}

// ProcessText runs the complete pipeline on the input text
func (p *Pipeline) ProcessText(ctx context.Context, text string) (*models.Graph, error) {
	// Step 1: Chunk the text
	log.Printf("Starting text chunking...")
	chunks, err := p.chunkText(text)
	if err != nil {
		return nil, fmt.Errorf("chunking failed: %w", err)
	}
	log.Printf("Chunking completed. Generated %d chunks", len(chunks))

	// Step 2: Embed the chunks
	log.Printf("Starting chunk embedding...")
	embeddedChunks, err := p.embedChunks(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	log.Printf("Embedding completed")

	// Step 3: Generate graph
	log.Printf("Starting graph generation...")
	graph, err := p.generateGraph(ctx, embeddedChunks)
	if err != nil {
		return nil, fmt.Errorf("graph generation failed: %w", err)
	}
	log.Printf("Graph generation completed. Generated graph with %d nodes and %d edges",
		len(graph.Nodes), len(graph.Edges))

	return graph, nil
}

// chunkText performs text chunking
func (p *Pipeline) chunkText(text string) ([]models.TextChunk, error) {
	return chunking.ChoiceChunker(text, p.config.UseAIChunking)
}

// embedChunks embeds the text chunks
func (p *Pipeline) embedChunks(ctx context.Context, chunks []models.TextChunk) ([]models.TextChunk, error) {
	return p.embedder.EmbedTextChunks(ctx, chunks)
}

// generateGraph generates the knowledge graph from embedded chunks
func (p *Pipeline) generateGraph(ctx context.Context, chunks []models.TextChunk) (*models.Graph, error) {
	return p.comparator.CompareChunks(ctx, chunks)
}

// SimpleProcessText provides a simplified interface for processing text with default configuration
func SimpleProcessText(ctx context.Context, text string, llmClient LLMClient) (*models.Graph, error) {
	config := DefaultPipelineConfig()
	pipeline, err := NewPipeline(config, llmClient)
	if err != nil {
		return nil, err
	}

	return pipeline.ProcessText(ctx, text)
}
