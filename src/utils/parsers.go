package utils

import (
	"bufio"
	"fmt"
	"my-modus-app/src/schemas"
	"strings"
)

// ParseMedlineResponse parses multiple MEDLINE format articles
func ParseMedlineResponse(content string) (*schemas.MedlineResponse, error) {
	response := &schemas.MedlineResponse{
		Articles: make([]*schemas.MedlineArticle, 0),
	}

	// Split content into individual articles (they're typically separated by blank lines)
	articles := strings.Split(content, "\n\n")

	for _, articleContent := range articles {
		if strings.TrimSpace(articleContent) == "" {
			continue
		}

		article, err := ParseMedline(articleContent)
		if err != nil {
			return nil, err
		}

		// Add PubMed URL
		article.PubMedURL = fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov/%s", article.PMID)
		response.Articles = append(response.Articles, article)
	}

	return response, nil
}

// ParseMedline parses a single MEDLINE format article
func ParseMedline(content string) (*schemas.MedlineArticle, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	article := &schemas.MedlineArticle{
		Authors:          make([]schemas.Author, 0),
		MeshTerms:        make([]string, 0),
		PublicationTypes: make([]string, 0),
	}

	var currentField string
	var currentValue strings.Builder
	var currentAuthor schemas.Author

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Check if this is a continuation line (starts with 6 spaces)
		if strings.HasPrefix(line, "      ") {
			currentValue.WriteString(" " + strings.TrimSpace(line))
			continue
		}

		// Process the previous field before starting a new one
		if currentField != "" {
			processField(article, currentField, currentValue.String(), &currentAuthor)
		}

		// Parse new field
		if len(line) < 6 {
			continue
		}

		parts := strings.SplitN(line, "- ", 2)
		if len(parts) != 2 {
			continue
		}

		currentField = strings.TrimSpace(parts[0])
		currentValue.Reset()
		currentValue.WriteString(strings.TrimSpace(parts[1]))
	}

	// Process the last field
	if currentField != "" {
		processField(article, currentField, currentValue.String(), &currentAuthor)
	}

	return article, nil
}

func processField(article *schemas.MedlineArticle, field, value string, currentAuthor *schemas.Author) {
	switch field {
	case "PMID":
		article.PMID = value
	case "TI":
		article.Title = value
	case "AB":
		article.Abstract = value
	case "FAU":
		currentAuthor.FullName = value
	case "AU":
		currentAuthor.LastName = value
		article.Authors = append(article.Authors, *currentAuthor)
		*currentAuthor = schemas.Author{} // Reset for next author
	case "AD":
		if len(article.Authors) > 0 {
			lastIdx := len(article.Authors) - 1
			article.Authors[lastIdx].Afiliation = value
		}
	case "MH":
		article.MeshTerms = append(article.MeshTerms, value)
	case "PT":
		article.PublicationTypes = append(article.PublicationTypes, value)
	case "LA":
		article.Language = value
	case "DP":
		article.JournalInfo.Date = value
	case "TA":
		article.JournalInfo.Abbreviation = value
	case "JT":
		article.JournalInfo.FullTitle = value
	case "VI":
		article.JournalInfo.Volume = value
	case "IP":
		article.JournalInfo.Issue = value
	case "PG":
		article.JournalInfo.Pages = value
	case "EDAT":
		article.DateAdded = value
	case "AID":
		if strings.Contains(value, "[doi]") {
			article.DOI = strings.TrimSuffix(value, " [doi]")
		}
	}
}
