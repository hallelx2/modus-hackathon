package processors

import (
    "fmt"
    "regexp"
    "strings"
    "unicode"
)

// Section struct remains the same as it's well-defined
type Section struct {
    Title   string
    Content string
    Type    string
}

// SectionExtractor enhanced with additional detection capabilities
type SectionExtractor struct {
    headers        map[string][]string
    supportedTypes []string
    // Adding new fields to enhance section detection
    sectionIndicators map[string][]string // Additional indicators for each format
    formatPatterns   []string            // Common formatting patterns
}

// NewSectionExtractor enhanced with more comprehensive section headers and indicators
func NewSectionExtractor() *SectionExtractor {
    headers := map[string][]string{
        "PubMed":  {"Background", "Methods", "Results", "Conclusions", "Objective", "Study Design", "Materials", "Discussion"},
        "Book":    {"Introduction", "Chapter", "Summary", "Appendix", "Preface", "Part", "Section", "Notes", "Bibliography"},
        "Legal":   {"Clause", "Sub-clause", "Definitions", "Conclusion", "Article", "Section", "Subsection", "Amendment"},
        "Generic": {"Abstract", "Introduction", "Discussion", "References", "Overview", "Summary", "Conclusion"},
    }

    // Additional indicators that might suggest section starts
    sectionIndicators := map[string][]string{
        "PubMed":  {"Purpose:", "Methodology:", "Findings:", "Keywords:", "Author contributions:"},
        "Book":    {"Chapter", "Volume", "Part", "Section", "Notes:", "Further reading:"},
        "Legal":   {"§", "Article", "Provision", "Term", "Whereas:", "Hereinafter:"},
        "Generic": {"Overview:", "Summary:", "Key points:", "Discussion points:", "Conclusion:"},
    }

    // Common formatting patterns that might indicate section starts
    formatPatterns := []string{
        `^\s*\d+\.\s+[A-Z]`,      // Numbered sections (e.g., "1. Introduction")
        `^\s*[IVXLCDM]+\.\s+`,    // Roman numerals
        `^\s*[A-Z]\.\s+`,         // Single letters (e.g., "A. Introduction")
        `^\s*§\s*\d+`,            // Section symbol
        `^\s*[-*•]\s+[A-Z]`,      // Bullet points with capitalized text
    }

    return &SectionExtractor{
        headers:          headers,
        supportedTypes:   []string{"PubMed", "Book", "Legal", "Generic"},
        sectionIndicators: sectionIndicators,
        formatPatterns:   formatPatterns,
    }
}

// DetectFormat enhanced to consider more indicators and patterns
func (se *SectionExtractor) DetectFormat(text string) string {
    // First check explicit headers as in original implementation
    for format, headers := range se.headers {
        headerCount := 0
        for _, header := range headers {
            if strings.Contains(strings.ToLower(text), strings.ToLower(header)) {
                headerCount++
            }
        }
        // If multiple headers of same format are found, it's more likely to be that format
        if headerCount >= 2 {
            return format
        }
    }

    // Check additional indicators if no clear format was found
    for format, indicators := range se.sectionIndicators {
        for _, indicator := range indicators {
            if strings.Contains(text, indicator) {
                return format
            }
        }
    }

    // Check formatting patterns
    for _, pattern := range se.formatPatterns {
        if regexp.MustCompile(pattern).MatchString(text) {
            return "Generic" // Default to Generic for pattern-based matches
        }
    }

    return "Unknown"
}

// ChunkBasedOnFormat enhanced with more sophisticated section detection
func (se *SectionExtractor) ChunkBasedOnFormat(text, format string) ([]Section, error) {
    headers, ok := se.headers[format]
    if !ok {
        return nil, fmt.Errorf("unsupported format: %s", format)
    }

    // Combine format-specific headers with their indicators
    allPatterns := headers
    if indicators, exists := se.sectionIndicators[format]; exists {
        allPatterns = append(allPatterns, indicators...)
    }

    // Create an enhanced regex pattern that includes formatting patterns
    sectionPattern := "(?i)(?m)^\\s*("
    sectionPattern += strings.Join(allPatterns, "|") + "|"
    sectionPattern += strings.Join(se.formatPatterns, "|") + ")"
    sectionPattern += "(:|\\s*$)"

    re := regexp.MustCompile(sectionPattern)
    matches := re.FindAllStringIndex(text, -1)

    if len(matches) == 0 {
        return nil, fmt.Errorf("no sections found for format: %s", format)
    }

    var sections []Section
    start := 0

    for i, match := range matches {
        headerStart, headerEnd := match[0], match[1]

        // Enhanced section boundary detection
        if i > 0 {
            content := text[start:headerStart]
            header := text[matches[i-1][0]:matches[i-1][1]]

            // Clean up header and content
            header = strings.TrimSpace(strings.TrimRight(header, ":"))
            content = strings.TrimSpace(content)

            // Additional validation of content boundaries
            if isValidSection(header, content) {
                sections = append(sections, Section{
                    Title:   header,
                    Content: content,
                    Type:    format,
                })
            }
        }
        start = headerEnd
    }

    // Handle the last section with enhanced boundary detection
    lastHeader := strings.TrimSpace(text[matches[len(matches)-1][0]:matches[len(matches)-1][1]])
    lastContent := strings.TrimSpace(text[start:])
    lastHeader = strings.TrimRight(lastHeader, ":")

    if isValidSection(lastHeader, lastContent) {
        sections = append(sections, Section{
            Title:   lastHeader,
            Content: lastContent,
            Type:    format,
        })
    }

    return sections, nil
}

// New helper function to validate section boundaries
func isValidSection(header, content string) bool {
    // Validate header
    if len(header) == 0 || len(content) == 0 {
        return false
    }

    // Check if header has proper capitalization
    if len(header) > 0 && !unicode.IsUpper(rune(header[0])) {
        return false
    }

    // Check reasonable length ratios
    if len(header) > len(content) {
        return false
    }

    return true
}

// ExtractSections extracts sections from the text based on the detected format
func (se *SectionExtractor) ExtractSections(text string) ([]Section, error) {
	// Detect the format
	format := se.DetectFormat(text)

	// Attempt chunking based on the detected format
	if format != "Unknown" {
		sections, err := se.ChunkBasedOnFormat(text, format)
		if err == nil && len(sections) > 0 {
			return sections, nil
		}
	}

	// If no format is detected, fallback to generic splitting logic
	return se.GenericChunking(text)
}

// GenericChunking provides a fallback chunking mechanism when no format is detected
func (se *SectionExtractor) GenericChunking(text string) ([]Section, error) {
	// Use a simple split by paragraphs
	paragraphs := strings.Split(text, "\n\n")
	var sections []Section
	for i, paragraph := range paragraphs {
		sections = append(sections, Section{
			Title:   fmt.Sprintf("Section-%d", i+1),
			Content: strings.TrimSpace(paragraph),
			Type:    "Miscellaneous",
		})
	}
	return sections, nil
}
