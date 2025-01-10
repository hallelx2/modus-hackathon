package schemas



//--------------------------------------
// SCHEMA: Relationships between chunks
//--------------------------------------
type RelationType string

type Relationship struct {
	SourceID string                 `json:"source_id"`
	TargetID string                 `json:"target_id"`
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


//--------------------------------------
// SCHEMA: Graph representation
//--------------------------------------
type Node struct {
    ID    string                 `json:"id"`
    Label string                 `json:"label"`
    Data  map[string]interface{} `json:"data"`
}
type Edge struct {
	ID         string                 `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Evidence []string `json:"evidence"`
}


type Graph struct {
    Nodes []Node `json:"nodes"`
    Edges []Edge `json:"edges"`
}
