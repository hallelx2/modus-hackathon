package main

import (
	"encoding/json"
	"fmt"
	"my-modus-app/graphgen/chunking"

	// graphgenModels "my-modus-app/graphgen/models"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

const modelName = "section-generator"

// GenerateText processes PubMed-style articles to extract entities and relationships
func GenerateText(articleText string) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", fmt.Errorf("failed to get model: %w", err)
	}

	// System instruction: Define the behavior of the model
	instruction := `
You are an advanced data scientist and knowledge graph expert. Your task is to analyze scientific articles, extract key entities (e.g., genes, proteins, drugs, diseases) and their relationships (e.g., interactions, pathways, causative effects), and model them into a graph schema.
Output the results in the following format:
Entity 1 -> Relation -> Entity 2
e.g.,
Gene A -> inhibits -> Protein B
Disease C -> treated by -> Drug D
Use concise language and ensure scientific accuracy in your extraction.
`

	// User prompt: Article text
	prompt := fmt.Sprintf(`
Here is the text of a scientific article. Extract entities and relationships as specified:

"%s"
`, articleText)

	// Prepare the model input
	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create model input: %w", err)
	}

	// Optional parameters for OpenAI chat
	input.Temperature = 0.7
	input.MaxTokens = 512

	// Invoke the model
	output, err := model.Invoke(input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

// ChunkingText chunks the text into different categories and returns the chunks
func ChunkingText(
	text string,
	maxChunkSize int,
	minChunkSize int,
	chunkOverlap int,
	preserveParagraphs bool,
	preserveSentences bool,
	// sectionHeaders []string,
) (string, error) {
	// Validate parameters before passing to the chunker
	if maxChunkSize <= 0 || minChunkSize <= 0 {
		return "", fmt.Errorf("chunk size values must be greater than 0")
	}
	if chunkOverlap < 0 {
		return "", fmt.Errorf("chunk overlap must be non-negative")
	}

	// Create a ChunkingConfig from the provided parameters
	config := chunking.ChunkingConfig{
		MaxChunkSize:       maxChunkSize,
		MinChunkSize:       minChunkSize,
		ChunkOverlap:       chunkOverlap,
		PreserveParagraphs: preserveParagraphs,
		PreserveSentences:  preserveSentences,
		// SectionHeaders:     sectionHeaders,
	}

	// Initialize the Chunker with the provided configuration
	chunker := chunking.NewChunker(config)

	// Process the input text and generate chunks
	chunks, err := chunker.ProcessText(text)
	if err != nil {
		return "", fmt.Errorf("error chunking text: %w", err)
	}

	// Convert chunks to JSON format
	chunksJSON, err := json.Marshal(chunks)
	if err != nil {
		return "", fmt.Errorf("error serializing chunks to JSON: %w", err)
	}

	// Return the JSON string of the chunks
	return string(chunksJSON), nil
}

func FallbackToLLMChunking(text string) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", fmt.Errorf("failed to get model: %w", err)
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
! Do not use code tags like like json as when formatting the string in markdown... Just return the Json

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
		return "", fmt.Errorf("failed to create model input: %w", err)
	}

	// Set model parameters
	input.Temperature = 0.2 // Minimize randomness for strict compliance
	input.MaxTokens = 2048  // Adjust token limit for larger outputs

	// Invoke the model
	output, err := model.Invoke(input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Extract the raw JSON part from the model's response and clean it
	cleanedOutput := output.Choices[0].Message.Content

	// Remove backticks and extra newlines, only keep valid JSON
	cleanedOutput = strings.TrimPrefix(cleanedOutput, "```json\n")
	cleanedOutput = strings.TrimSuffix(cleanedOutput, "\n```")

	// Ensure the cleaned output is a valid JSON array
	if !strings.HasPrefix(cleanedOutput, "[") || !strings.HasSuffix(cleanedOutput, "]") {
		return "", fmt.Errorf("invalid JSON format: %s", cleanedOutput)
	}

	// Return the cleaned JSON string
	return cleanedOutput, nil

}
