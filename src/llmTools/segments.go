package llmtools

import (
    "fmt"
    "strings"
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

// GenerateContentSections generates structured section points based on the review type
func GenerateContentSections(topic string, reviewType ReviewType) ([]string, error) {
    // Retrieve the OpenAI chat model
    model, err := models.GetModel[openai.ChatModel](modelName)
    if err != nil {
        return nil, fmt.Errorf("failed to get model: %w", err)
    }

    // Get the appropriate instruction based on review type
    instruction := getInstructionForType(reviewType)

    // User prompt with the topic
    prompt := fmt.Sprintf(`
Generate 5-6 key discussion points for the following topic, following the format and depth specified:

Topic: "%s"
`, topic)

    // Prepare the model input
    input, err := model.CreateInput(
        openai.NewSystemMessage(instruction),
        openai.NewUserMessage(prompt),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create model input: %w", err)
    }

    // Configure parameters based on review type
    input.Temperature = getTemperatureForType(reviewType)
    input.MaxTokens = 1024

    // Invoke the model
    output, err := model.Invoke(input)
    if err != nil {
        return nil, fmt.Errorf("failed to invoke model: %w", err)
    }

    // Parse the output into a slice of strings
    sections := parseSections(output.Choices[0].Message.Content)
    return sections, nil
}

// getInstructionForType returns the appropriate system instruction based on review type
func getInstructionForType(reviewType ReviewType) string {
    instructions := map[ReviewType]string{
        QuickReview: `
You are an expert content organizer specialized in creating concise overviews.
Your task is to generate 5-6 essential points that would cover a topic at a high level.

Guidelines:
- Focus on fundamental concepts and key takeaways
- Keep points broad but informative
- Ensure logical progression from basics to conclusion
- Avoid technical jargon unless absolutely necessary
- Each point should be self-contained and clear

Output format:
- Return only numbered points (1-6)
- Each point should be a clear section heading
- No explanations or additional text
- No subcategories or nested points`,

        DetailedReport: `
You are a professional report structure specialist.
Your task is to outline 5-6 comprehensive sections for a detailed report.

Guidelines:
- Include both theoretical and practical aspects
- Ensure coverage of background, methodology, and implications
- Balance depth with breadth
- Include data analysis or evaluation sections where relevant
- Consider stakeholder perspectives

Output format:
- Return only numbered points (1-6)
- Each point should be a clear section heading
- No explanations or additional text
- No subcategories or nested points`,

        SystematicReview: `
You are an academic research methodologist.
Your task is to outline 5-6 sections for a systematic review following academic standards.

Guidelines:
- Follow PRISMA guidelines where applicable
- Include methodology and quality assessment sections
- Ensure comprehensive coverage of evidence synthesis
- Focus on reproducibility and rigor
- Include meta-analysis considerations where relevant

Output format:
- Return only numbered points (1-6)
- Each point should be a clear section heading
- No explanations or additional text
- No subcategories or nested points`,

        TechnicalGuide: `
You are a technical documentation specialist.
Your task is to outline 5-6 sections for a comprehensive technical guide.

Guidelines:
- Progress from setup to advanced implementation
- Include prerequisites and requirements
- Focus on practical implementation steps
- Cover troubleshooting and best practices
- Consider scalability and optimization

Output format:
- Return only numbered points (1-6)
- Each point should be a clear section heading
- No explanations or additional text
- No subcategories or nested points`,

        Tutorial: `
You are an educational content designer.
Your task is to outline 5-6 sections for a learning-focused tutorial.

Guidelines:
- Structure content from basic to advanced concepts
- Include hands-on exercises or examples
- Focus on practical application
- Include assessment or practice opportunities
- Consider different learning styles

Output format:
- Return only numbered points (1-6)
- Each point should be a clear section heading
- No explanations or additional text
- No subcategories or nested points`,
    }

    return instructions[reviewType]
}

// getTemperatureForType returns the appropriate temperature setting based on review type
func getTemperatureForType(reviewType ReviewType) float64 {
    temperatures := map[ReviewType]float64{
        QuickReview:      0.3,  // More creative for high-level overview
        DetailedReport:   0.2,  // More structured for reports
        SystematicReview: 0.1,  // Most deterministic for academic content
        TechnicalGuide:   0.2,  // Structured for technical content
        Tutorial:         0.3,  // More creative for educational content
    }

    return temperatures[reviewType]
}

// parseSections converts the model output into a clean slice of sections
func parseSections(output string) []string {
    // Split the output into lines
    lines := strings.Split(strings.TrimSpace(output), "\n")

    // Clean and collect valid sections
    sections := make([]string, 0, 6)
    for _, line := range lines {
        line = strings.TrimSpace(line)
        // Skip empty lines or lines without content
        if line == "" {
            continue
        }
        // Remove numbering and clean up
        if idx := strings.Index(line, "."); idx != -1 {
            line = strings.TrimSpace(line[idx+1:])
        }
        sections = append(sections, line)
    }

    return sections
}
