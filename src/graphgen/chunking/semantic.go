package chunking

import (
	"errors"
	"fmt"
	"my-modus-app/graphgen/models"
	"strings"
	"time"

	"github.com/jdkato/prose/v2"

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
// SentenceGroup represents a collection of sentences that stays within our constraints
type SentenceGroup struct {
    Sentences []string
    TokenCount int
}

func splitIntoSentences(text string) ([]SentenceGroup, error) {
    // Create a new document from the input text
    doc, err := prose.NewDocument(text)
    if err != nil {
        return nil, fmt.Errorf("failed to create document: %w", err)
    }

    // Initialize our result slice to store groups of sentences
    var result []SentenceGroup

    // Initialize variables for the current group we're building
    currentGroup := SentenceGroup{
        Sentences: make([]string, 0),
        TokenCount: 0,
    }

    // Process each sentence
    for _, sent := range doc.Sentences() {
        // Get approximate token count for this sentence
        // A simple approximation: count words and punctuation
        tokens := len(strings.Fields(sent.Text))

        // Check if adding this sentence would exceed our constraints
        if len(currentGroup.Sentences) >= 4 ||
           currentGroup.TokenCount + tokens > 2048 {
            // If current group is not empty, add it to results
            if len(currentGroup.Sentences) > 0 {
                result = append(result, currentGroup)
                // Start a new group
                currentGroup = SentenceGroup{
                    Sentences: make([]string, 0),
                    TokenCount: 0,
                }
            }
        }

        // Add the sentence to the current group
        currentGroup.Sentences = append(currentGroup.Sentences, sent.Text)
        currentGroup.TokenCount += tokens
    }

    // Add the last group if it contains any sentences
    if len(currentGroup.Sentences) > 0 {
        result = append(result, currentGroup)
    }

    // Check if we found any sentences
    if len(result) == 0 {
        return nil, errors.New("no sentences found in the input text")
    }

    return result, nil
}

