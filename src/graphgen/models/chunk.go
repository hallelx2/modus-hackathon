package models

import (
	"time"
)

// ChunkMetadata contains metadata about a text chunk, including information
// about its position in the document, associated citations, keywords, entity types,
// timestamp of extraction, and confidence score. This metadata helps in identifying
// the context of a chunk and assessing the quality of the chunking process.
type ChunkMetadata struct {
	// StartIndex is the starting index of the chunk in the original text.
	// This helps in identifying the exact position of the chunk in the text.
	StartIndex int `json:"start_index"`

	// EndIndex is the ending index of the chunk in the original text.
	// It marks the end of the chunk and helps in determining the length of the chunk.
	EndIndex int `json:"end_index"`

	// Section represents the section name (e.g., Introduction, Methods, Results)
	// where the chunk was found within the document.
	Section string `json:"section"`

	// Citations contains a list of references or citations associated with the chunk.
	// This can be used for citation tracking and research.
	Citations []string `json:"citations"`

	// Keywords is a list of important keywords associated with the chunk.
	// Keywords are typically domain-specific and help with content categorization.
	Keywords []string `json:"keywords"`

	// EntityTypes defines the types of entities present in the chunk (e.g., gene, protein, disease).
	// This is useful for knowledge graph generation and relationship extraction.
	EntityTypes []string `json:"entity_types"`

	// Timestamp represents the time when the chunk was processed or created.
	// This can be useful for tracking data freshness and workflow timing.
	Timestamp time.Time `json:"timestamp"`

	// Confidence is the model's confidence score in the accuracy of the chunk's extraction.
	// This score can be used to assess the reliability of the extracted chunk.
	Confidence float64 `json:"confidence"`
}

// TextChunk represents a chunk of text extracted from a larger document or corpus.
// It includes the chunk content, its embedding (vector representation), metadata,
// and any identified relationships with other entities in the text.
type TextChunk struct {
	// ID is a unique identifier for the text chunk. This helps in referencing the chunk
	// within the system or for tracking purposes.
	ID string `json:"id"`

	// Content contains the actual text of the chunk. This is the portion of the document
	// that was extracted and may contain important information or entities.
	Content string `json:"content"`

	// Embedding is the vector representation of the text chunk. It is used for machine learning
	// tasks such as searching, similarity comparison, and clustering.
	Embedding []float32 `json:"embedding"`

	// Metadata contains additional information about the chunk, such as its position in the text,
	// associated citations, keywords, and more.
	Metadata ChunkMetadata `json:"metadata"`

	// Score represents a numeric score that can be used to rank or evaluate the chunk.
	// This might be based on relevance, confidence, or other scoring criteria.
	Score float64 `json:"score"`

	// Relations is a list of relationships identified within the chunk. Each relationship
	// consists of two entities and a relation type. This is useful for building knowledge graphs
	// and understanding the interactions between entities.
	Relations []Relationship `json:"relations"`
}

