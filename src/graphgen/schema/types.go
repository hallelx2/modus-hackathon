package schema

type SchemaType string

const (
	DrugInteraction    SchemaType = "drug-interaction"
	GeneDisease        SchemaType = "gene-disease"
	ProteinInteraction SchemaType = "protein-interaction"
	ClinicalOutcome    SchemaType = "clinical-outcome"
	MolecularPathway   SchemaType = "molecular-pathway"
	// Add more schema types as needed
)

type Schema struct {
	Type        SchemaType             `json:"type"`
	Entities    []string               `json:"entities"`
	Relations   []string               `json:"relations"`
	Properties  []string               `json:"properties"`
	Constraints map[string]interface{} `json:"constraints"`
	Embedding   []float32              `json:"embedding"`
}
