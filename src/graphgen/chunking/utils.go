package chunking

import (
	"my-modus-app/graphgen/models"
	"regexp"
	"strings"
)

// TextCleaner removes unwanted artifacts from text
type TextCleaner struct {
	patterns []string
}

func NewTextCleaner() *TextCleaner {
	return &TextCleaner{
		patterns: []string{
			`\s+`,                      // Multiple whitespace
			`\[[\d,\s]+\]`,             // Citation brackets
			`\((?:[^()]|\([^)]*\))*\)`, // Nested parentheses
		},
	}
}

func (c *TextCleaner) Clean(text string) string {
	cleaned := text
	for _, pattern := range c.patterns {
		re := regexp.MustCompile(pattern)
		cleaned = re.ReplaceAllString(cleaned, " ")
	}
	return strings.TrimSpace(cleaned)
}

func (c *Chunker) applyOverlap(chunks []models.TextChunk) []models.TextChunk {
    if len(chunks) == 0 || c.config.ChunkOverlap <= 0 {
        return chunks // No overlap needed or no chunks to process
    }

    var result []models.TextChunk

    for i, chunk := range chunks {
        // Add the current chunk to the result
        if i == 0 {
            result = append(result, chunk) // First chunk, no overlap
        } else {
            // Add overlap from the previous chunk
            prevChunk := result[len(result)-1]
            overlapEnd := min(len(prevChunk.Content), c.config.ChunkOverlap)
            overlapContent := prevChunk.Content[len(prevChunk.Content)-overlapEnd:]

            // Create a new chunk with overlap and current content
            newChunk := models.TextChunk{
                ID:      chunk.ID,
                Content: overlapContent + chunk.Content,
                Metadata: models.ChunkMetadata{
                    StartIndex: prevChunk.Metadata.EndIndex - overlapEnd,
                    EndIndex:   chunk.Metadata.EndIndex,
                    Section:    chunk.Metadata.Section,
                    Timestamp:  chunk.Metadata.Timestamp,
                },
                Embedding: chunk.Embedding,
                Score:     chunk.Score,
                Relations: chunk.Relations,
            }
            result = append(result, newChunk)
        }
    }

    return result
}

// Utility function to calculate the minimum of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

