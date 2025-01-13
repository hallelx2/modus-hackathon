package schemas

import (
	"time"
)

type ChunkMetadata struct {
	StartIndex  int                    `json:"ChunkMetadata.start_index"`
	EndIndex    int                    `json:"ChunkMetadata.end_index"`
	Section     string                 `json:"ChunkMetadata.section"`
	Citations   []string               `json:"ChunkMetadata.citations"`
	Keywords    []string               `json:"ChunkMetadata.keywords"`
	EntityTypes []string               `json:"ChunkMetadata.entity_types"`
	Timestamp   time.Time              `json:"ChunkMetadata.timestamp"`
	Confidence  float64                `json:"ChunkMetadata.confidence"`
	MedlineData MedlineArticleMetadata `json:"ChunkMetadata.medline_data"`
}

type TextChunk struct {
	ID        string         `json:"TextChunk.id"`
	UserID    string         `json:"TextChunk.user_id"`
	Content   string         `json:"TextChunk.content"`
	Embedding []float32      `json:"TextChunk.embedding"`
	Metadata  ChunkMetadata  `json:"TextChunk.metadata"`
	Score     float64        `json:"TextChunk.score"`
	Relations []Relationship `json:"TextChunk.relations"`
}

type TextChunkList struct {
	Chunks []TextChunk `json:"chunks"`
}
