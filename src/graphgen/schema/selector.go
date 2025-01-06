package schema

import (
	"context"
	"my-modus-app/graphgen/embedding"
	"my-modus-app/graphgen/models"
)

type SchemaSelector struct {
	registry  *SchemaRegistry
	embedder  *embedding.Embedder
	vectorOps *embedding.VectorOperations
}

func NewSchemaSelector(registry *SchemaRegistry, embedder *embedding.Embedder) *SchemaSelector {
	return &SchemaSelector{
		registry:  registry,
		embedder:  embedder,
		vectorOps: embedding.NewVectorOperations(),
	}
}

func (ss *SchemaSelector) SelectSchema(ctx context.Context, chunks []models.TextChunk) (SchemaType, float64, error) {
	scores := make(map[SchemaType]float64)

	// Calculate schema scores based on multiple factors
	for _, chunk := range chunks {
		// Calculate embedding similarity
		embedScore := ss.calculateEmbeddingScore(chunk)

		// Calculate entity type match
		entityScore := ss.calculateEntityScore(chunk)

		// Calculate relationship pattern match
		relationScore := ss.calculateRelationScore(chunk)

		// Combine scores with weights
		for schemaType := range ss.registry.Schemas {
			scores[schemaType] += ss.combineScores(
				embedScore[schemaType],
				entityScore[schemaType],
				relationScore[schemaType],
			)
		}
	}

	// Select best matching schema
	bestSchema, bestScore := ss.selectBestSchema(scores)

	return bestSchema, bestScore, nil
}
