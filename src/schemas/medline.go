package schemas

// SchemaType is a type representing different schema types
type MedlineArticleMetadata struct {
	PMID             string      `json:"MedlineArticleMetadata.pmid"`
	Title            string      `json:"MedlineArticleMetadata.title"`
	Authors          []Author    `json:"MedlineArticleMetadata.authors"`
	MeshTerms        []string    `json:"MedlineArticleMetadata.mesh_terms"`
	JournalInfo      JournalInfo `json:"MedlineArticleMetadata.journal_info"`
	PublicationTypes []string    `json:"MedlineArticleMetadata.publication_types"`
	Language         string      `json:"MedlineArticleMetadata.language"`
	DateAdded        string      `json:"MedlineArticleMetadata.date_added"`
	DOI              string      `json:"MedlineArticleMetadata.doi"`
	PubMedURL        string      `json:"MedlineArticleMetadata.pubmed_url"`
}

type Author struct {
	FullName   string `json:"Author.full_name"`
	LastName   string `json:"Author.last_name"`
	Afiliation string `json:"Author.affiliation"`
}

type JournalInfo struct {
	Abbreviation string `json:"JournalInfo.abbreviation"`
	FullTitle    string `json:"JournalInfo.full_title"`
	Volume       string `json:"JournalInfo.volume"`
	Issue        string `json:"JournalInfo.issue"`
	Pages        string `json:"JournalInfo.pages"`
	Date         string `json:"JournalInfo.date"`
}

// SearchResult represents the esearch response structure
type SearchResult struct {
	ESearchResult struct {
		Count    string   `json:"count"`
		RetMax   string   `json:"retmax"`
		RetStart string   `json:"retstart"`
		IdList   []string `json:"idlist"`
	} `json:"esearchresult"`
}

// MedlineResponse represents multiple articles
type MedlineResponse struct {
	Articles []*MedlineArticle `json:"Articles"`
}

// MedlineArticle represents a single article in MEDLINE format
type MedlineArticle struct {
	PMID             string      `json:"PMID"`
	Title            string      `json:"Title"`
	Abstract         string      `json:"Abstract"`
	Authors          []Author    `json:"Authors"`
	MeshTerms        []string    `json:"MeshTerms"`
	JournalInfo      JournalInfo `json:"JournalInfo"`
	PublicationTypes []string    `json:"PublicationTypes"`
	Language         string      `json:"Language"`
	DateAdded        string      `json:"DateAdded"`
	DOI              string      `json:"DOI"`
	PubMedURL        string      `json:"PubMedURL"` // Added for convenience
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
