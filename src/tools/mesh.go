package tools

import (
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

const modelName = "text-generator"

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
