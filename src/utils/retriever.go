package utils


import (
	"encoding/json"
	"fmt"
	"strings"

	"my-modus-app/src/schemas"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
)


func GetPubMedDetails(meshTerms string) ([]*schemas.MedlineArticle, error) {
	baseSearchURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"
	baseFetchURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi"

	// Step 1: Search for IDs using esearch
	searchURL := fmt.Sprintf("%s?db=pubmed&term=%s&retmode=json&retmax=5", baseSearchURL, meshTerms)
	searchResponse, err := http.Fetch(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search PubMed: %w", err)
	}

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

	medlineResponse, err := ParseMedlineResponse(fetchResponse.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse MEDLINE response: %w", err)
	}

	return medlineResponse.Articles, nil
}


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
