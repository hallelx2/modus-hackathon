package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// ReviewType represents different types of content structure
type ReviewType string

const (
	QuickReview      ReviewType = "quick"
	DetailedReport   ReviewType = "report"
	SystematicReview ReviewType = "systematic"
	TechnicalGuide   ReviewType = "technical"
	Tutorial         ReviewType = "tutorial"
)

// ResponseSchema defines the structure for generated content
type ResponseSchema struct {
	SectionTitle   string `json:"section_title"`
	SectionContent string `json:"section_content"`
}

// GenerateContent generates and writes content for all sections sequentially
func GenerateContent(topic string, reviewType ReviewType, description string) ([]ResponseSchema, error) {
	// Step 1: Generate sections sequentially
	sections, err := GenerateContentSections(topic, reviewType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content sections: %w", err)
	}

	// Step 2: Initialize a slice to hold section content
	sectionContents := []ResponseSchema{}

	// Step 3: Process each section sequentially
	for _, section := range sections {
		// Add a delay to handle rate-limiting
		time.Sleep(5 * time.Second)

		// Generate content for each section
		content, err := GenerateSectionContent(topic, section, reviewType)
		if err != nil {
			// Log the error but continue processing other sections
			fmt.Printf("Error generating content for section '%s': %v\n", section, err)
			continue
		}

		// Append the result to the list
		sectionContents = append(sectionContents, ResponseSchema{
			SectionTitle:   section,
			SectionContent: content,
		})
	}

	return sectionContents, nil
}

// GenerateContentSections generates structured section points based on the review type
func GenerateContentSections(topic string, reviewType ReviewType) ([]string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	instruction := getInstructionForType(reviewType)

	prompt := fmt.Sprintf(`
You are a highly intelligent and structured assistant specializing in generating content outlines. Your task is to create key discussion points for a topic.

#### Instructions:
1. Focus on generating 2-3 **comprehensive and actionable** sections for the topic.
2. The sections should:
   - Cover diverse subtopics within the main topic.
   - Be concise yet descriptive enough to guide detailed content generation later.
3. Keep the tone and complexity aligned with the specified review type.
4. Ensure the sections are ordered logically for maximum readability and coherence.

#### Topic:
"%s"

#### Review Type:
"%s"

#### Example Output:
1. Introduction to [Topic]: Cover the background, importance, and current state.
2. Key Challenges or Considerations in [Topic]: Dive into specific challenges or nuances.
3. Solutions or Future Directions: Propose actionable insights or areas for innovation.

Now, generate the key discussion points for the topic in a similar structure.
`, topic, reviewType)

	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create model input: %w", err)
	}

	input.Temperature = getTemperatureForType(reviewType)
	input.MaxTokens = 1024

	output, err := model.Invoke(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	return parseSections(output.Choices[0].Message.Content), nil
}

// GenerateSectionContent generates content for a specific section
func GenerateSectionContent(topic, segment string, reviewType ReviewType) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", fmt.Errorf("failed to get model: %w", err)
	}

	instruction := getInstructionForType(reviewType)

	prompt := fmt.Sprintf(`
You are a top-tier assistant specializing in generating high-quality, structured content. Your task is to expand on a given section title for a topic.

#### Instructions:
1. Provide **detailed and structured content** for the section.
   - Include an introduction, subpoints, and a conclusion (if applicable).
   - Use concise paragraphs for readability.
2. Ensure the content aligns with the review type and adheres to the following:
   - Depth and tone: %s
   - Format: Formal, logical, and easy to understand.
3. Use **RAG (Retrieval-Augmented Generation)** if context is available; otherwise, rely on internal knowledge.
4. Avoid unnecessary repetition and focus on providing actionable, informative content.

#### Topic:
"%s"

#### Section Title:
"%s"

#### Example Output:
Title: Key Challenges in [Topic]
- **Challenge 1**: Description of the first challenge.
- **Challenge 2**: Explanation of another critical issue.
- **Possible Solutions**: Insights into how these challenges can be mitigated.

Now, generate content for the section.
`, reviewType, topic, segment)

	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create model input: %w", err)
	}

	input.Temperature = getTemperatureForType(reviewType)
	input.MaxTokens = 1024

	output, err := model.Invoke(input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

// Helper Functions

func getInstructionForType(reviewType ReviewType) string {
	instructions := map[ReviewType]string{
		QuickReview:      "Your task is to generate high-level overviews...",
		DetailedReport:   "Your task is to outline detailed report sections...",
		SystematicReview: "Your task is to create sections following PRISMA guidelines...",
		TechnicalGuide:   "Your task is to draft a comprehensive technical guide...",
		Tutorial:         "Your task is to structure content for an educational tutorial...",
	}
	return instructions[reviewType]
}

func getTemperatureForType(reviewType ReviewType) float64 {
	temperatures := map[ReviewType]float64{
		QuickReview:      0.3,
		DetailedReport:   0.2,
		SystematicReview: 0.1,
		TechnicalGuide:   0.2,
		Tutorial:         0.3,
	}
	return temperatures[reviewType]
}

func parseSections(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	sections := make([]string, 0, 6)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if idx := strings.Index(line, "."); idx != -1 {
			line = strings.TrimSpace(line[idx+1:])
		}
		sections = append(sections, line)
	}
	return sections
}
