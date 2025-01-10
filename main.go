package main

import (
	"encoding/json"
	"fmt"
	"log"
	"my-modus-app/graphgen/chunking"

	// "my-modus-app/graphgen/knowledge"
	funcModels "my-modus-app/graphgen/models"
	"my-modus-app/graphgen/processing"

	// graphgenModels "my-modus-app/graphgen/models"
	"strings"

	"my-modus-app/graphgen/knowledge"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"

	"my-modus-app/src/graph"
	llmtools "my-modus-app/src/llmTools"

	// "my-modus-app/src/schemas"

	// "my-modus-app/src/schemas"
	"my-modus-app/src/utils"

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

func GenerateAdvancedMeSHKeywords(articleText string) (string, error) {
	// Retrieve the OpenAI chat model
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", fmt.Errorf("failed to get model: %w", err)
	}

	// System instruction: Define the behavior of the model for MeSH keyword generation
	instruction := `
You are a medical librarian with expert knowledge of MeSH (Medical Subject Headings) and PubMed search strategies.
Your task is to generate advanced MeSH terms, incorporating Boolean operators (AND, OR, NOT) to create a search query that retrieves as many relevant articles as possible from PubMed.

**Output format**:
1. Group related MeSH terms with OR for inclusivity.
2. Combine broader and narrower terms logically using AND for relevance.
3. Use NOT only to exclude irrelevant topics.
4. Output the terms in PubMed-ready syntax. Do not include explanations or any other text. Provide only the search query.

For example:
("Diabetes Mellitus, Type 2"[MeSH] OR "Insulin Resistance"[MeSH]) AND ("Metformin"[MeSH] OR "Hypoglycemic Agents"[MeSH]) NOT "Type 1 Diabetes Mellitus"[MeSH]
`

	// User prompt: Article text
	prompt := fmt.Sprintf(`
Analyze the following text and generate an advanced PubMed search query using MeSH terms, Boolean operators, and the format described above:

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

	// Configure optional parameters for the OpenAI chat model
	input.Temperature = 0.2 // Low temperature for deterministic and precise output
	input.MaxTokens = 1024  // Allow for longer responses with complex queries

	// Invoke the model
	output, err := model.Invoke(input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Return the PubMed search query
	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

func FetchPubMedAccessions(meshTerms string) (string, error) {
	// Construct the PubMed API URL
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"
	url := fmt.Sprintf("%s?db=pubmed&term=%s&retmode=json&retmax=100", baseURL, meshTerms)

	// Perform the GET request
	response, err := http.Fetch(url)
	if err != nil {
		return "", fmt.Errorf("error fetching PubMed data: %w", err)
	}

	// Check if the response is successful
	if !response.Ok() {
		return "", fmt.Errorf("failed to fetch PubMed data. Status: %d %s", response.Status, response.StatusText)
	}

	// Return the response body as a string
	return response.Text(), nil
}

// GetPubMedAccessions queries the PubMed API with MeSH terms and returns a list of PMIDs
func GetPubMedAccessions(meshTerms string) ([]string, error) {
	type PubMedSearchResult struct {
		ESearchResult struct {
			Count    string   `json:"Count"`
			RetMax   string   `json:"RetMax"`
			RetStart string   `json:"RetStart"`
			IdList   []string `json:"IdList"`
		} `json:"esearchresult"`
	}
	// PubMed E-utilities URL for esearch
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"

	requestURL := fmt.Sprintf("%s?db=pubmed&term=%s&retmode=json&retmax=100", baseURL, meshTerms)

	// Make the HTTP GET request using Modus
	response, err := http.Fetch(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from PubMed API: %w", err)
	}

	// Parse the JSON response
	var searchResult PubMedSearchResult
	response.JSON(&searchResult)

	// Return the list of PMIDs
	return searchResult.ESearchResult.IdList, nil
}

func GetPubMedDetails(meshTerms string) ([]*processing.MedlineArticle, error) {
	// PubMed E-utilities URLs
	baseSearchURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"
	baseFetchURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi"

	// Step 1: Search for IDs using esearch
	searchURL := fmt.Sprintf("%s?db=pubmed&term=%s&retmode=json&retmax=100", baseSearchURL, meshTerms)
	searchResponse, err := http.Fetch(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search PubMed: %w", err)
	}

	// Parse search response to get IDs
	var searchResult funcModels.SearchResult
	if err := json.Unmarshal([]byte(searchResponse.Text()), &searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	if len(searchResult.ESearchResult.IdList) == 0 {
		return nil, fmt.Errorf("no results found for the query")
	}

	// Step 2: Retrieve detailed metadata using efetch
	ids := strings.Join(searchResult.ESearchResult.IdList, ",")
	fetchURL := fmt.Sprintf("%s?db=pubmed&id=%s&rettype=medline&retmode=text", baseFetchURL, ids)
	fetchResponse, err := http.Fetch(fetchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PubMed details: %w", err)
	}

	// Parse the MEDLINE format response
	medlineResponse, err := processing.ParseMedlineResponse(fetchResponse.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse MEDLINE response: %w", err)
	}

	return medlineResponse.Articles, nil
}

// ChunkingText chunks the text into different categories and returns the chunks
func ChunkingText(
	text string,
) (string, error) {
	// Hardcoded parameter values
	maxChunkSize := 500
	minChunkSize := 100
	chunkOverlap := 50
	preserveParagraphs := true
	preserveSentences := true

	// Validate parameters
	if maxChunkSize <= 0 || minChunkSize <= 0 {
		return "", fmt.Errorf("chunk size values must be greater than 0")
	}
	if chunkOverlap < 0 {
		return "", fmt.Errorf("chunk overlap must be non-negative")
	}

	// Create config for the chunker
	config := chunking.ChunkingConfig{
		MaxChunkSize:       maxChunkSize,
		MinChunkSize:       minChunkSize,
		ChunkOverlap:       chunkOverlap,
		PreserveParagraphs: preserveParagraphs,
		PreserveSentences:  preserveSentences,
	}

	// Initialize chunker with config
	chunker := chunking.NewChunker(config)

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
		return "", fmt.Errorf("error chunking text: %w", err)
	}

	// Convert to JSON
	chunksJSON, err := json.Marshal(chunks)
	if err != nil {
		return "", fmt.Errorf("error serializing chunks to JSON: %w", err)
	}

	return string(chunksJSON), nil
}

func DocumentProcessorModel(
	text string,
) (string, error) {
	// Initialize parameters directly within the function
	modelName := "section-generator" // Set model name to 'text-generator'
	maxChunkSize := 1000             // Set max chunk size
	minChunkSize := 500              // Set min chunk size
	chunkOverlap := 50               // Set chunk overlap
	preserveParagraphs := true       // Set preserve paragraphs flag
	preserveSentences := true        // Set preserve sentences flag

	// Call FallbackToLLMChunking to process the document and get the chunked JSON response
	allChunks, err := chunking.FallbackToLLMChunking(
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
		return "", fmt.Errorf("failed to process document: %w", err)
	}

	// Convert chunks to JSON format
	chunksJSON, err := json.Marshal(allChunks)
	if err != nil {
		return "", fmt.Errorf("error serializing chunks to JSON: %w", err)
	}

	// Return the chunked JSON string
	return string(chunksJSON), nil
}

func ChunkBasedOnChoice(text string, useAI bool) (string, error) {
	// Calling the function with the parameters
	chunks, err := chunking.ChoiceChunker(text, useAI)

	// Error handling
	if err != nil {
		return "", fmt.Errorf("failed to chunk the text: %w", err)
	}

	// Convert chunks to JSON format
	chunksJSON, err := json.Marshal(chunks)
	if err != nil {
		return "", fmt.Errorf("error serializing chunks to JSON: %w", err)
	}

	return string(chunksJSON), nil
}

// GetContentSections takes a topic string and content type string, returns a list of sections to cover
func GetContentSections(topic, contentType string) ([]string, error) {
	// Clean inputs
	topic = strings.TrimSpace(topic)
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	// Map content type to review type
	var reviewType llmtools.ReviewType
	switch contentType {
	case "quick":
		reviewType = llmtools.QuickReview
	case "report":
		reviewType = llmtools.DetailedReport
	case "systematic":
		reviewType = llmtools.SystematicReview
	case "technical":
		reviewType = llmtools.TechnicalGuide
	case "tutorial":
		reviewType = llmtools.Tutorial
	default:
		reviewType = llmtools.QuickReview // Default to quick review if type not specified
	}

	// Generate sections using the llmtools package
	sections, err := llmtools.GenerateContentSections(topic, reviewType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sections: %w", err)
	}

	return sections, nil
}

func GraphGenerator(meshKeyword string) (*knowledge.KnowledgeGraph, error) {
	graph, err := knowledge.GenerateKnowledgeGraph(meshKeyword)

	if err != nil {
		return nil, fmt.Errorf("error generating the knowledge graph: %v", err)
	}
	return graph, nil
}

// RetrieveAndChunk retrieves PubMed details, chunks the articles, and returns a list of JSON strings.
func RetrieveAndChunk(meshText string, useAi bool) ([]string, error) {
	articles, err := utils.GetPubMedDetails(meshText)
	if err != nil {
		return nil, fmt.Errorf("error retrieving articles: %s", err)
	}

	chunks, err := graph.ChunkAndEmbedManyMedlineRetrievals(articles, useAi)
	if err != nil {
		return nil, fmt.Errorf("error chunking the multiple entries: %w", err)
	}

	// Prepare a list of JSON strings
	var jsonStrings []string
	for _, chunk := range chunks {
		jsonData, err := json.Marshal(chunk)
		if err != nil {
			return nil, fmt.Errorf("error marshaling chunk to JSON: %w", err)
		}
		jsonStrings = append(jsonStrings, string(jsonData))
	}

	return jsonStrings, nil
}

// func GraphRetrieveAndChunk(meshText string, useAi bool) ([]*schemas.TextChunk, error) {
// 	articles, err := utils.GetPubMedDetails(meshText)
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving articles: %s", err)
// 	}
// 	var chunks []*schemas.TextChunk
// 	for _, article := range articles {
// 		chunk, err := graph.ChunkAndEmbedOneMedlineRetrieval(*article, useAi)
// 		if err != nil {
// 			return nil, fmt.Errorf("error genenerating the chunks and embedding")
// 		}
// 		chunks = append(chunks, chunk...)
// 	}
// 	return chunks, nil
// }
