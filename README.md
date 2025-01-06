I want you to enhance the schema selection algorithm

but apart from that, i even want this to be more sophisticated... here is how an llm graph builder works

The LLM Graph Builder follows the process you learned earlier in the course:
Gather the data
Chunk the data, creating something like langchain's CharacterTextSPlit or recursiveCHaracter text split, but I wnat you to omplement something more sophisticated for this use case
Vectorize the data, this is so that RAG can be performed later on on the chunked piece of data
Pass the data to an LLM to extract nodes and relationships
Use the output to generate the graph
here is a function to perform text embedding with the modus sdk

package main
import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)
func GetEmbeddingsForTextWithOpenAI(text string) ([]float32, error) {
	results, err := GetEmbeddingsForTextsWithOpenAI(text)
	if err != nil {
		return nil, err
	}
	return results[0], nil
}
func GetEmbeddingsForTextsWithOpenAI(texts ...string) ([][]float32, error) {
	model, err := models.GetModel[openai.EmbeddingsModel]("openai-embeddings")
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
And this one is for vector operations

/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */
package main
import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/vectors"
)
func Add(a, b []float64) []float64 {
	return vectors.Add(a, b)
}
func AddInPlace(a, b []float64) []float64 {
	vectors.AddInPlace(a, b)
	return a
}
func Subtract(a, b []float64) []float64 {
	return vectors.Subtract(a, b)
}
func SubtractInPlace(a, b []float64) []float64 {
	vectors.SubtractInPlace(a, b)
	return a
}
func AddNumber(a []float64, b float64) []float64 {
	return vectors.AddNumber(a, b)
}
func AddNumberInPlace(a []float64, b float64) []float64 {
	vectors.AddNumberInPlace(a, b)
	return a
}
func SubtractNumber(a []float64, b float64) []float64 {
	return vectors.SubtractNumber(a, b)
}
func SubtractNumberInPlace(a []float64, b float64) []float64 {
	vectors.SubtractNumberInPlace(a, b)
	return a
}
func MultiplyNumber(a []float64, b float64) []float64 {
	return vectors.MultiplyNumber(a, b)
}
func MultiplyNumberInPlace(a []float64, b float64) []float64 {
	vectors.MultiplyNumberInPlace(a, b)
	return a
}
func DivideNumber(a []float64, b float64) []float64 {
	return vectors.DivideNumber(a, b)
}
func DivideNumberInPlace(a []float64, b float64) []float64 {
	vectors.DivideNumberInPlace(a, b)
	return a
}
func Dot(a, b []float64) float64 {
	return vectors.Dot(a, b)
}
func Magnitude(a []float64) float64 {
	return vectors.Magnitude(a)
}
func Normalize(a []float64) []float64 {
	return vectors.Normalize(a)
}
func Sum(a []float64) float64 {
	return vectors.Sum(a)
}
func Product(a []float64) float64 {
	return vectors.Product(a)
}
func Mean(a []float64) float64 {
	return vectors.Mean(a)
}
func Min(a []float64) float64 {
	return vectors.Min(a)
}
func Max(a []float64) float64 {
	return vectors.Max(a)
}
func Abs(a []float64) []float64 {
	return vectors.Abs(a)
}
func AbsInPlace(a []float64) []float64 {
	vectors.AbsInPlace(a)
	return a
}
func EuclidianDistance(a, b []float64) float64 {
	return vectors.EuclidianDistance(a, b)
}

can you help me with this?




package main import ( "strings" "github.com/hypermodeinc/modus/sdk/go/pkg/models" "github.com/hypermodeinc/modus/sdk/go/pkg/models/openai" ) // this model name should match the one defined in the modus.json manifest file const modelName = "text-generator" func GenerateText(instruction, prompt string) (string, error) { model, err := models.GetModel[openai.ChatModel](modelName) if err != nil { return "", err } input, err := model.CreateInput( openai.NewSystemMessage(instruction), openai.NewUserMessage(prompt), ) if err != nil { return "", err } // this is one of many optional parameters available for the OpenAI chat interface input.Temperature = 0.7 output, err := model.Invoke(input) if err != nil { return "", err } return strings.TrimSpace(output.Choices[0].Message.Content), nil } Can you do this in go, make the system instruction something to be an advanced dtat sciientist who can model information into graphs and ask it to model a piece of information i want into a schema for a graph... I want you to optimised the way you write the prompt in that we are dealing with scientific articles here such as those that we find on piubmed and we wnat it to extract relationships and entities from there... can you go ahead and hel- me with this?

Add functionality for processing multiple papers in batch?
Include specific parsing for different types of biomedical relationships?
Add validation for the output format?
Create a more specific schema for certain types of studies?
Yes do all these things, but do not include a main function, all i want in the end is a dunction like GenerateGraphRelationship with all these things that you have said... i will also appreciate if you have like a multiple schemas for the graph that will fit most use cases of scientific paper publications which the model will check for to see whcih schema he wants to go for to generate the text for and then go on to generate this based on the chosen schema... make this schema tobust as much as possible so that we can  extract as much relationship as possible from here
