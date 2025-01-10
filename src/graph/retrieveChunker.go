package graph

import (
	"fmt"
	// "strings"
	"my-modus-app/src/schemas"
	"my-modus-app/src/processors"
	"my-modus-app/src/utils"
)


func ChunkAndEmbedOneMedlineRetrieval(article schemas.MedlineArticle, ai bool) ([]schemas.TextChunk, error) {
	useAI := ai
	text := article.Abstract

	// Convert the article to metadata
	metadata := schemas.ConvertToMetadata(article)

	// Chunk the text using the processor
	chunks, err := processors.ChoiceChunker(text, useAI)
	if err != nil {
		return nil, fmt.Errorf("error chunking the abstract: %s", err)
	}

	// Update the metadata for each chunk
	for i := range chunks {
		chunks[i].Metadata.MedlineData = metadata
		embedding, err := utils.GetEmbeddingsForTextWithOpenAI(chunks[i].Content)

		if err !=nil {
			return nil, fmt.Errorf("error generating the embedding of the chunk: %v", err)
		}
		chunks[i].Embedding = embedding
	}

	return chunks, nil
}

// Modify the function signature to accept pointer slice
func ChunkAndEmbedManyMedlineRetrievals(articles []*schemas.MedlineArticle, ai bool) ([]schemas.TextChunk, error) {
    var allChunks []schemas.TextChunk  // Now just a single slice of TextChunk

    for _, article := range articles {
        // Chunk a single article
        chunks, err := ChunkAndEmbedOneMedlineRetrieval(*article, ai)
        if err != nil {
            return nil, fmt.Errorf("error processing article with PMID %s: %s", article.PMID, err)
        }
        // Append the chunks to the overall slice
        allChunks = append(allChunks, chunks...)
    }

    return allChunks, nil
}
