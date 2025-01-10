package schemas

// SchemaType is a type representing different schema types
type MedlineArticleMetadata struct {
	PMID string
	Title string
	Authors []Author
	MeshTerms []string
	JournalInfo	JournalInfo
	PublicationTypes []string
	Language string
	DateAdded string
	DOI string
	PubMedURL string
}

type Author struct {
	FullName string
	LastName string
	Afiliation string
}

type JournalInfo struct {
	Abbreviation string
	FullTitle string
	Volume string
	Issue string
	Pages string
	Date string

}

// SearchResult represents the esearch response structure
type SearchResult struct {
	ESearchResult struct {
		Count   string   `json:"count"`
		RetMax  string   `json:"retmax"`
		RetStart string  `json:"retstart"`
		IdList  []string `json:"idlist"`
	} `json:"esearchresult"`
}

// MedlineResponse represents multiple articles
type MedlineResponse struct {
	Articles []*MedlineArticle
}

// MedlineArticle represents a single article in MEDLINE format
type MedlineArticle struct {
	PMID        string
	Title       string
	Abstract    string
	Authors     []Author
	MeshTerms   []string
	JournalInfo JournalInfo
	PublicationTypes []string
	Language    string
	DateAdded   string
	DOI         string
	PubMedURL   string // Added for convenience
}

func ConvertToMetadata(article MedlineArticle) MedlineArticleMetadata {
	return MedlineArticleMetadata{
		PMID:             article.PMID,
		Title:            article.Title,
		Authors:          article.Authors,
		MeshTerms:        article.MeshTerms,
		JournalInfo:      article.JournalInfo,
		PublicationTypes: article.PublicationTypes,
		Language:         article.Language,
		DateAdded:        article.DateAdded,
		DOI:              article.DOI,
		PubMedURL:        article.PubMedURL,
	}
}
