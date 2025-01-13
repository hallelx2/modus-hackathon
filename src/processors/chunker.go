package processors

import (
	"encoding/json"
	"fmt"
	"log"
	models "my-modus-app/src/schemas"
	"strings"

	generativeModel "github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// ChunkingStrategy interface for different chunking strategies
type ChunkingStrategy interface {
	Chunk(text string) ([]models.TextChunk, error)
}

// ChunkingConfig holds configuration for chunking
type ChunkingConfig struct {
	MaxChunkSize       int  `json:"max_chunk_size"`
	MinChunkSize       int  `json:"min_chunk_size"`
	ChunkOverlap       int  `json:"chunk_overlap"`
	PreserveParagraphs bool `json:"preserve_paragraphs"`
	PreserveSentences  bool `json:"preserve_sentences"`
	// SectionHeaders     []string `json:"section_headers"`
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
func (c *Chunker) ProcessText(
	text string,
	maxChunkSize int,
	minChunkSize int,
	chunkOverlap int,
	preserveParagraphs bool,
	preserveSentences bool,
) ([]models.TextChunk, error) {
	// Create a new ChunkingConfig based on input parameters
	config := ChunkingConfig{
		MaxChunkSize:       maxChunkSize,
		MinChunkSize:       minChunkSize,
		ChunkOverlap:       chunkOverlap,
		PreserveParagraphs: preserveParagraphs,
		PreserveSentences:  preserveSentences,
	}

	// Instantiate a new SemanticChunker with the updated config
	semanticChunker := NewSemanticChunker(config)

	// Extract sections using the SectionExtractor
	sections, err := c.sectionExtractor.ExtractSections(text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sections: %w", err)
	}

	var chunks []models.TextChunk

	// Process each section
	for _, section := range sections {
		// Chunk the section using the newly configured SemanticChunker
		sectionChunks, err := semanticChunker.ChunkSection(section)
		if err != nil {
			return nil, fmt.Errorf("failed to chunk section '%s': %w", section.Title, err)
		}

		// Apply overlap to chunks if needed
		processedChunks := c.applyOverlap(sectionChunks)
		chunks = append(chunks, processedChunks...)
	}

	return chunks, nil
}

// FallbackToLLMChunking takes in a text and parameters, processes it through LLM chunking,
// and returns the chunked JSON string output.
func FallbackToLLMChunking(
	text string,
	maxChunkSize int,
	minChunkSize int,
	chunkOverlap int,
	preserveParagraphs bool,
	preserveSentences bool,
	modelName string,
) ([]models.TextChunk, error) {
	// Get the model from the available models
	model, err := generativeModel.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	// Updated Instruction for LLM
	instruction := `
You are a scientific document parsing expert. Your task is to analyze and divide a scientific/non-scientific article into logical sections. Each section must be classified based on its title, content, and type, and then represented in a structured JSON format.

### Requirements for the Task:
1. **Title**: Provide a clear and descriptive title of each section.
2. **Content**: The full text of the section.
3. **Type**: Classify the section based on its structure (e.g., Introduction, Methods, Results, Discussion, Conclusion).

### JSON Representation Format:
Return the output strictly as a JSON array of objects, each representing a section of the document, in the following format:

[
  {
    "Title": "<Title of the section>",
    "Content": "<Content of the section>",
    "Type": "<Type of the section>",
  },
  ...
]

**Note**
! Do not use code tags like json as when formatting the string in markdown... Just return the Json

**Analysis Requirements:**
1. **Context Awareness:**
   - Consider explicit and implicit indicators of document position (e.g., headers, linguistic markers, content focus).
   - Analyze linguistic cues like transitional phrases, topic shifts, and thematic progressions.
   - Pay attention to changes in style, tone, or technical depth.

2. **Section Identification Strategy:**
   - Detect discourse markers (e.g., "In conclusion", "We investigated").
   - Recognize vocabulary indicative of section types (e.g., "methods include" suggests a methodology section).
   - Use relationships between topics and transitions to maintain logical flow.

3. **Document Structure Awareness:**
   - Use contextual clues to determine if content belongs to:
     * **Introduction**: Sets up background, objectives, or hypotheses.
     * **Methods**: Describes methodology and materials.
     * **Results**: Presents findings.
     * **Discussion**: Interprets findings in context.
     * **Conclusion**: Synthesizes ideas, highlights implications.
   - Adapt to variations in structure while maintaining scientific norms.

4. **Section Properties to Capture:**
   - **Title**: Clear and descriptive (e.g., "Introduction", "Results").
   - **Type**: Both structural (e.g., "Introduction") and functional (e.g., "Background").

**Output Instructions:**
- Divide the document into 5–7 meaningful sections.
- Ensure each section is comprehensive without losing logical flow.
- Strictly output the JSON array as described above. Do not include any additional text or explanations.

Here is the document to analyze:
`

	// Append the text to the instruction
	prompt := fmt.Sprintf("Here is the document to analyze\n\n%s", text)

	// Prepare input for the model
	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create model input: %w", err)
	}

	// Set model parameters
	input.Temperature = 0.2 // Minimize randomness for strict compliance
	input.MaxTokens = 2048  // Adjust token limit for larger outputs

	// Invoke the model
	output, err := model.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	// Extract the raw JSON part from the model's response and clean it
	cleanedOutput := output.Choices[0].Message.Content

	// Create a ChunkingConfig from the provided parameters
	config := ChunkingConfig{
		MaxChunkSize:       maxChunkSize,
		MinChunkSize:       minChunkSize,
		ChunkOverlap:       chunkOverlap,
		PreserveParagraphs: preserveParagraphs,
		PreserveSentences:  preserveSentences,
	}

	// Remove backticks and extra newlines, only keep valid JSON
	cleanedOutput = strings.TrimPrefix(cleanedOutput, "```json\n")
	cleanedOutput = strings.TrimSuffix(cleanedOutput, "\n```")

	// Ensure the cleaned output is a valid JSON array
	if !strings.HasPrefix(cleanedOutput, "[") || !strings.HasSuffix(cleanedOutput, "]") {
		return nil, fmt.Errorf("invalid JSON format: %s", cleanedOutput)
	}

	// Unmarshal the cleaned output to sections
	var sections []Section
	err = json.Unmarshal([]byte(cleanedOutput), &sections)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Initialize the Chunker with the provided configuration
	chunker := NewSemanticChunker(config)

	// Process sections and generate chunks
	var allChunks []models.TextChunk
	for _, section := range sections {
		// Pass each section to the chunker
		chunks, err := chunker.ChunkSection(section) // Pass individual section, not a slice
		if err != nil {
			return nil, fmt.Errorf("error chunking section: %w", err)
		}
		// Append the chunks to the overall list
		allChunks = append(allChunks, chunks...)
	}

	// Return the JSON string of the chunks
	return allChunks, nil
}

func ChoiceChunker(text string, use_ai bool) ([]models.TextChunk, error) {
	// Initialize parameters directly within the function
	modelName := "section-generator" // Set model name to 'section-generator as seen in the modus.json'
	maxChunkSize := 1000             // Set max chunk size
	minChunkSize := 500              // Set min chunk size
	chunkOverlap := 50               // Set chunk overlap
	preserveParagraphs := true       // Set preserve paragraphs flag
	preserveSentences := true        // Set preserve sentences flag
	if use_ai {
		// Call FallbackToLLMChunking to process the document and get the chunked JSON response
		allChunks, err := FallbackToLLMChunking(
			text,
			maxChunkSize,
			minChunkSize,
			chunkOverlap,
			preserveParagraphs,
			preserveSentences,
			modelName,
		)
		if err != nil {
			log.Printf("Error in chunking document: %v", err)
			return nil, fmt.Errorf("failed to process document: %w", err)
		}

		return allChunks, nil
	} else {
		// Validate parameters
		if maxChunkSize <= 0 || minChunkSize <= 0 {
			return nil, fmt.Errorf("chunk size values must be greater than 0")
		}
		if chunkOverlap < 0 {
			return nil, fmt.Errorf("chunk overlap must be non-negative")
		}

		// Create config for the chunker
		config := ChunkingConfig{
			MaxChunkSize:       maxChunkSize,
			MinChunkSize:       minChunkSize,
			ChunkOverlap:       chunkOverlap,
			PreserveParagraphs: preserveParagraphs,
			PreserveSentences:  preserveSentences,
		}

		// Initialize chunker with config
		chunker := NewChunker(config)

		// Process the text
		chunks, err := chunker.ProcessText(
			text,
			maxChunkSize,
			minChunkSize,
			chunkOverlap,
			preserveParagraphs,
			preserveSentences,
		)
		if err != nil {
			return nil, fmt.Errorf("error chunking text: %w", err)
		}
		return chunks, nil
	}
}
