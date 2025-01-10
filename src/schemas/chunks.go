package schemas

import (
	"time"
)

type ChunkMetadata struct {
	StartIndex  int                    `json:"start_index"`
	EndIndex    int                    `json:"end_index"`
	Section     string                 `json:"section"`
	Citations   []string               `json:"citations"`
	Keywords    []string               `json:"keywords"`
	EntityTypes []string               `json:"entity_types"`
	Timestamp   time.Time              `json:"timestamp"`
	Confidence  float64                `json:"confidence"`
	MedlineData MedlineArticleMetadata `json:"medline_data"`
}

type TextChunk struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Content   string         `json:"content"`
	Embedding []float32      `json:"embedding"`
	Metadata  ChunkMetadata  `json:"metadata"`
	Score     float64        `json:"score"`
	Relations []Relationship `json:"relations"`
}

type TextChunkList struct {
	Chunks []TextChunk `json:"chunks"`
}
