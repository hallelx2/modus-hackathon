package chunking

import (
	"fmt"
	"my-modus-app/graphgen/models"
)

// ChunkingStrategy interface for different chunking strategies
type ChunkingStrategy interface {
	Chunk(text string) ([]models.TextChunk, error)
}

// ChunkingConfig holds configuration for chunking
type ChunkingConfig struct {
	MaxChunkSize       int      `json:"max_chunk_size"`
	MinChunkSize       int      `json:"min_chunk_size"`
	ChunkOverlap       int      `json:"chunk_overlap"`
	PreserveParagraphs bool     `json:"preserve_paragraphs"`
	PreserveSentences  bool     `json:"preserve_sentences"`
	SectionHeaders     []string `json:"section_headers"`
}

// Chunker is the main struct for text chunking
type Chunker struct {
	config           ChunkingConfig
	sectionExtractor *SectionExtractor
	semanticChunker  *SemanticChunker
}

// NewChunker initializes a new Chunker with the provided configuration
func NewChunker(config ChunkingConfig) *Chunker {
	return &Chunker{
		config:           config,
		sectionExtractor: NewSectionExtractor(), // Initialize with default headers
		semanticChunker:  NewSemanticChunker(config),
	}
}

// ProcessText processes the input text by splitting it into sections and chunking each section
func (c *Chunker) ProcessText(text string) ([]models.TextChunk, error) {
	// Extract sections using the SectionExtractor
	sections, err := c.sectionExtractor.ExtractSections(text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sections: %w", err)
	}

	var chunks []models.TextChunk

	// Process each section
	for _, section := range sections {
		// Chunk the section using the SemanticChunker
		sectionChunks, err := c.semanticChunker.ChunkSection(section)
		if err != nil {
			return nil, fmt.Errorf("failed to chunk section '%s': %w", section.Title, err)
		}

		// Apply overlap to chunks if needed
		processedChunks := c.applyOverlap(sectionChunks)
		chunks = append(chunks, processedChunks...)
	}

	return chunks, nil
}

