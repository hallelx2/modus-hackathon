package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hypermodeinc/modus/sdk/go/pkg/http"

	"my-modus-app/src/dg"
	"my-modus-app/src/graph"
	"my-modus-app/src/processors"
	"my-modus-app/src/tools"
	"my-modus-app/src/user"

	"my-modus-app/src/schemas"

	// "my-modus-app/src/schemas"
	"my-modus-app/src/utils"
)

// const modelName = "section-generator"

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

func GetPubMedDetails(meshTerms string) ([]*schemas.MedlineArticle, error) {
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
	var searchResult schemas.SearchResult
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
	medlineResponse, err := utils.ParseMedlineResponse(fetchResponse.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse MEDLINE response: %w", err)
	}

	return medlineResponse.Articles, nil
}

// 	return string(chunksJSON), nil
// }

func ChunkBasedOnChoice(text string, useAI bool) (string, error) {
	// Calling the function with the parameters
	chunks, err := processors.ChoiceChunker(text, useAI)

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
	var reviewType tools.ReviewType
	switch contentType {
	case "quick":
		reviewType = tools.QuickReview
	case "report":
		reviewType = tools.DetailedReport
	case "systematic":
		reviewType = tools.SystematicReview
	case "technical":
		reviewType = tools.TechnicalGuide
	case "tutorial":
		reviewType = tools.Tutorial
	default:
		reviewType = tools.QuickReview // Default to quick review if type not specified
	}

	// Generate sections using the llmtools package
	sections, err := tools.GenerateContentSections(topic, reviewType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sections: %w", err)
	}

	return sections, nil
}

// RetrieveAndChunk retrieves PubMed details, chunks the articles, and returns a list of JSON strings.
func RetrieveAndChunk(title string, useAi bool) ([]string, error) {
	meshText, err := tools.GenerateAdvancedMeSHKeywords(title)
	if err != nil {
		return nil, fmt.Errorf("error generating advanced mesh keywords: %w", err)
	}
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

// Adds a user to the database
func Signup(email, name, password string) (*schemas.User, error) {
	// Hash the password
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing the password: %w", err)
	}

	// Create a new user object
	user := schemas.User{
		ID:        uuid.NewString(),
		Name:      name,
		Password:  hashedPassword,
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Add user to the database
	addedUser, err := dg.AddUserToDatabase(user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user to database: %w", err)
	}

	return addedUser, nil
}

func SignupWithDgraph(email, name, password string) (map[string]string, error) {
	// Hash the password
	uids, err := user.Signup(email, name, password)

	if err != nil {
		return nil, fmt.Errorf("error creating the user: %w", err)
	}

	return uids, nil
}

func LoginWithDgraph(email, password string) (*schemas.LoginUser, error) {
	user, err := user.Login(email, password)
	if err != nil {
		return nil, fmt.Errorf("could not log user in: %w", err)
	}

	return user, nil
}

// ContentGeneratorFunction generates content for a specific topic and returns itimport (

func ContentGeneratorFunction(topic string, reviewType string, description string) (string, error) {
	// Call the GenerateContent function from the llmtools package

	reviewtype := tools.ReviewType(reviewType)
	content, err := tools.GenerateContent(topic, reviewtype, description)
	if err != nil {
		// Return the error so the caller can handle it
		return "", err
	}

	// Marshal the content into a JSON string
	contentJSON, err := json.Marshal(content)
	if err != nil {
		// Return the error if JSON marshalling fails
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	// Return the JSON string
	return string(contentJSON), nil
}

// func GenerateContent(topic string, reviewType llmtools.ReviewType, description string) ([]*llmtools.ResponseSchema, error) {
// 	// Step 1: Generate sections sequentially
// 	sections, err := llmtools.GenerateContentSections(topic, reviewType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate content sections: %w", err)
// 	}

// 	// Step 2: Initialize a slice to hold section content
// 	sectionContents := []llmtools.ResponseSchema{}

// 	// Step 3: Process each section sequentially
// 	for _, section := range sections {
// 		// Add a delay to handle rate-limiting
// 		time.Sleep(2 * time.Second)

// 		// Generate content for each section
// 		content, err := llmtools.GenerateSectionContent(topic, section, reviewType)
// 		if err != nil {
// 			// Log the error but continue processing other sections
// 			fmt.Printf("Error generating content for section '%s': %v\n", section, err)
// 			continue
// 		}

// 		// Append the result to the list
// 		sectionContents = append(sectionContents, llmtools.ResponseSchema{
// 			SectionTitle:   section,
// 			SectionContent: content,
// 		})
// 	}

// 	return sectionContents, nil
// }
