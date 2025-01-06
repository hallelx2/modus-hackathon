package models

type RelationType string

type Relationship struct {
	SourceID   string                 `json:"source_id"`
	TargetID   string                 `json:"target_id"`
	Type       RelationType           `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Confidence float64                `json:"confidence"`
	Evidence   []string               `json:"evidence"`
}

type RelationshipGroup struct {
	ID            string                 `json:"id"`
	Relationships []Relationship         `json:"relationships"`
	Context       []TextChunk            `json:"context"`
	Metadata      map[string]interface{} `json:"metadata"`
}
