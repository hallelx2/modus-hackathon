package schema

type SchemaRegistry struct {
	Schemas map[SchemaType]Schema
}

func NewSchemaRegistry() *SchemaRegistry {
	registry := &SchemaRegistry{
		Schemas: make(map[SchemaType]Schema),
	}

	// Initialize with default schemas
	registry.registerDefaultSchemas()

	return registry
}

func (sr *SchemaRegistry) registerDefaultSchemas() {
	// Register drug interaction schema
	sr.Schemas[DrugInteraction] = Schema{
		Type:       DrugInteraction,
		Entities:   []string{"Drug", "Protein", "Pathway"},
		Relations:  []string{"INHIBITS", "ACTIVATES", "METABOLIZES"},
		Properties: []string{"mechanism", "effect_strength", "clinical_significance"},
		Constraints: map[string]interface{}{
			"required_evidence":    "clinical_trial|in_vitro_study",
			"confidence_threshold": 0.7,
		},
	}

	// Register gene disease schema
	sr.Schemas[GeneDisease] = Schema{
		Type:       GeneDisease,
		Entities:   []string{"Gene", "Disease", "Phenotype", "Variant"},
		Relations:  []string{"ASSOCIATED_WITH", "CAUSES", "CONTRIBUTES_TO"},
		Properties: []string{"association_type", "effect_size", "population_background"},
		Constraints: map[string]interface{}{
			"required_evidence": "gwas|clinical_study",
			"min_sample_size":   100,
		},
	}

	// Add more schemas...
}
