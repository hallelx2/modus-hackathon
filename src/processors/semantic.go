package processors

import (
	"fmt"
	models "my-modus-app/src/schemas"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SemanticChunker struct {
    config ChunkingConfig
}

func NewSemanticChunker(config ChunkingConfig) *SemanticChunker {
    return &SemanticChunker{config: config}
}

func (sc *SemanticChunker) ChunkSection(section Section) ([]models.TextChunk, error) {
    // Ensure section content is not empty
    if strings.TrimSpace(section.Content) == "" {
        return nil, fmt.Errorf("section content is empty for section: %s", section.Title)
    }

    // Split the section content into sentences
    sentences, err := splitIntoSentences(section.Content) // Utility function to split text into sentences
    if err != nil {
        return nil, fmt.Errorf("failed to split section content into sentences for section: %s, error: %w", section.Title, err)
    }

    var chunks []models.TextChunk
    var currentChunk strings.Builder
    startIdx := 0

    for _, sentence := range sentences {
        if currentChunk.Len()+len(sentence) > sc.config.MaxChunkSize {
            chunkContent := currentChunk.String()
            if strings.TrimSpace(chunkContent) == "" {
                return nil, fmt.Errorf("chunk content is empty while processing section: %s", section.Title)
            }

            chunks = append(chunks, sc.createChunk(chunkContent, startIdx, section))
            startIdx += len(chunkContent)
            currentChunk.Reset()
        }
        currentChunk.WriteString(sentence + " ")
    }

    // Add the last chunk if any content remains
    if currentChunk.Len() > 0 {
        chunkContent := currentChunk.String()
        if strings.TrimSpace(chunkContent) == "" {
            return nil, fmt.Errorf("last chunk content is empty for section: %s", section.Title)
        }

        chunks = append(chunks, sc.createChunk(chunkContent, startIdx, section))
    }

    // Check if any chunks were created
    if len(chunks) == 0 {
        return nil, fmt.Errorf("no chunks were created for section: %s", section.Title)
    }

    return chunks, nil
}


func (sc *SemanticChunker) createChunk(content string, startIdx int, section Section) models.TextChunk {
    return models.TextChunk{
        ID:      uuid.NewString(),
        Content: strings.TrimSpace(content),
        Metadata: models.ChunkMetadata{
            StartIndex: startIdx,
            EndIndex:   startIdx + len(content),
            Section:    section.Title,
            Timestamp:  time.Now(),
        },
    }
}

// Utility function to split text into sentences
func splitIntoSentences(text string) ([]string, error){
    // Simplistic sentence splitter; consider using a library for production use
    return strings.Split(text, ". "), nil
}
