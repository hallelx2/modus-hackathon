package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"my-modus-app/graphgen/models"
)

// ComparisonResult represents the relationships found between two chunks
type ComparisonResult struct {
	ChunkIDs    [2]string           `json:"chunk_ids"`
	Relations   []models.Relationship `json:"relations"`
	SharedNodes []models.Node        `json:"shared_nodes"`
	Confidence  float64              `json:"confidence"`
}

// LLMClient interface defines the contract for LLM interactions
type LLMClient interface {
	GenerateCompletion(ctx context.Context, prompt string) (string, error)
}

// Config holds configuration options for the comparator
type Config struct {
	MaxParallel   int
	Timeout       time.Duration
	BatchSize     int
	MinConfidence float64
}

// ChunkComparator handles relationship detection between text chunks
type ChunkComparator struct {
	model       LLMClient
	config      Config
	maxParallel int
	timeout     time.Duration
}

// NewChunkComparator creates a new instance with the provided configuration
func NewChunkComparator(model LLMClient, config Config) *ChunkComparator {
	return &ChunkComparator{
		model:       model,
		config:      config,
		maxParallel: config.MaxParallel,
		timeout:     config.Timeout,
	}
}

// DefaultConfig returns default configuration values
func DefaultConfig() Config {
	return Config{
		MaxParallel:   5,
		Timeout:       2 * time.Minute,
		BatchSize:     10,
		MinConfidence: 0.7,
	}
}

// CompareChunks performs pairwise comparison of text chunks
func (cc *ChunkComparator) CompareChunks(ctx context.Context, chunks []models.TextChunk) (*models.Graph, error) {
	pairs := generateChunkPairs(chunks)
	results := make(chan ComparisonResult, len(pairs))
	errChan := make(chan error, len(pairs))

	var wg sync.WaitGroup
	sem := make(chan struct{}, cc.maxParallel)

	for _, pair := range pairs {
		wg.Add(1)
		go func(pair [2]models.TextChunk) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			result, err := cc.comparePair(ctx, pair[0], pair[1])
			if err != nil {
				errChan <- fmt.Errorf("comparison failed for chunks %s and %s: %w",
					pair[0].ID, pair[1].ID, err)
				return
			}
			results <- result
		}(pair)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf("multiple comparison errors occurred: %v", errors)
	}

	return cc.buildGraphFromResults(results), nil
}

// comparePair compares two chunks using LLM
func (cc *ChunkComparator) comparePair(ctx context.Context, chunk1, chunk2 models.TextChunk) (ComparisonResult, error) {
	prompt := fmt.Sprintf(`Compare these two text chunks and identify relationships:

Chunk 1 (%s):
"%s"

Chunk 2 (%s):
"%s"

Identify:
1. Shared entities/concepts
2. Relationships between chunks
3. Supporting evidence
4. Confidence scores

Output as JSON:
{
  "relations": [...],
  "shared_nodes": [...],
  "confidence": 0.9
}`,
		chunk1.ID, chunk1.Content,
		chunk2.ID, chunk2.Content,
	)

	completion, err := cc.model.GenerateCompletion(ctx, prompt)
	if err != nil {
		return ComparisonResult{}, fmt.Errorf("LLM completion failed: %w", err)
	}

	var result ComparisonResult
	if err := json.Unmarshal([]byte(completion), &result); err != nil {
		return ComparisonResult{}, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	result.ChunkIDs = [2]string{chunk1.ID, chunk2.ID}
	return result, nil
}

// buildGraphFromResults combines comparison results into a graph
func (cc *ChunkComparator) buildGraphFromResults(results chan ComparisonResult) *models.Graph {
	graph := &models.Graph{
		Nodes: make([]models.Node, 0),
		Edges: make([]models.Edge, 0),
	}

	nodeMap := make(map[string]models.Node)
	edgeMap := make(map[string]models.Edge)

	for result := range results {
		// Add shared nodes
		for _, node := range result.SharedNodes {
			if existing, exists := nodeMap[node.ID]; exists {
				existing.Sources = uniqueStrings(append(existing.Sources, node.Sources...))
				for k, v := range node.Properties {
					existing.Properties[k] = v
				}
				nodeMap[node.ID] = existing
			} else {
				nodeMap[node.ID] = node
			}
		}

		// Add relationships as edges
		for _, relation := range result.Relations {
			edge := models.Edge{
				ID:         uuid.New().String(),
				Source:     relation.SourceID,
				Target:     relation.TargetID,
				Type:       string(relation.Type),
				Properties: relation.Properties,
				Evidence:   relation.Evidence,
			}

			key := fmt.Sprintf("%s-%s-%s", edge.Source, edge.Type, edge.Target)
			if existing, exists := edgeMap[key]; exists {
				existing.Evidence = uniqueStrings(append(existing.Evidence, edge.Evidence...))
				for k, v := range edge.Properties {
					existing.Properties[k] = v
				}
				edgeMap[key] = existing
			} else {
				edgeMap[key] = edge
			}
		}
	}

	// Convert maps to slices
	for _, node := range nodeMap {
		graph.Nodes = append(graph.Nodes, node)
	}
	for _, edge := range edgeMap {
		graph.Edges = append(graph.Edges, edge)
	}

	return graph
}

// generateChunkPairs creates all possible pairs of chunks
func generateChunkPairs(chunks []models.TextChunk) [][2]models.TextChunk {
	var pairs [][2]models.TextChunk
	for i := 0; i < len(chunks); i++ {
		for j := i + 1; j < len(chunks); j++ {
			pairs = append(pairs, [2]models.TextChunk{chunks[i], chunks[j]})
		}
	}
	return pairs
}

// uniqueStrings removes duplicates from a string slice
func uniqueStrings(strings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, str := range strings {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}
