package processors

import (
	"my-modus-app/src/schemas"
	"regexp"
	"strings"
)

type TextCleaner struct {
	patterns []string
}

func NewTextCleaner() *TextCleaner {
	return &TextCleaner{
		patterns: []string{
			`\s+`,                      // Multiple whitespaces
			`\[[\d,\s]+\]`,             // Citation brackets
			`\((?:[^()]|\([^)]*\))*\)`, // nested parentheses
		},
	}
}

func (tc *TextCleaner) Clean(text string) (string, error) {
	cleaned := text
	for _, pattern := range tc.patterns {
		re := regexp.MustCompile(pattern)
		cleaned = re.ReplaceAllString(cleaned, " ")
	}
	return strings.TrimSpace(cleaned), nil
}

func (c *Chunker) applyOverlap(chunks []schemas.TextChunk) []schemas.TextChunk {
	if len(chunks) == 0 || c.config.ChunkOverlap <= 0 {
		return chunks
	}
	var result []schemas.TextChunk
	for i, chunk := range chunks {
		// Add the current chunk to the result
		if i == 0 {
			result = append(result, chunk)
		} else {
			// Add the overlap from the previous chunk
			prevChunk := result[len(result)-1]
			overlapEnd := min(len(prevChunk.Content), c.config.ChunkOverlap)
			overlapContent := prevChunk.Content[len(prevChunk.Content)-overlapEnd:]

			// Create a new chunk with overlap and current content
			newChunk := schemas.TextChunk{
				ID:      chunk.ID,
				Content: overlapContent + chunk.Content,
				Metadata: schemas.ChunkMetadata{
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
