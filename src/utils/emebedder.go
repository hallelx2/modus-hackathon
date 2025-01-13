package utils

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

func GetEmbeddingsForTextsWithOpenAI(texts ...string) ([][]float32, error) {
	model, err := models.GetModel[openai.EmbeddingsModel]("embeddings")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(texts)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	results := make([][]float32, len(output.Data))
	for i, d := range output.Data {
		results[i] = d.Embedding
	}

	return results, nil
}

func GetEmbeddingsForTextWithOpenAI(text string) ([]float32, error) {
	results, err := GetEmbeddingsForTextsWithOpenAI(text)
	if err != nil {
		return nil, err
	}

	return results[0], nil
}
