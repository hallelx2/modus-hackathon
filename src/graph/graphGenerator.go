package graph

// import (
// 	"encoding/json"
// 	"fmt"
// 	"my-modus-app/src/schemas"
// 	"my-modus-app/src/utils"
// 	"time"

// 	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
// )

// const dgraphConn = "dg-conn"

// // AddToDgraph processes articles and adds them to Dgraph
// func AddToDgraph(meshText string, useAi bool) (map[string]string, error) {
// 	articles, err := utils.GetPubMedDetails(meshText)
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving articles: %s", err)
// 	}

// 	chunks, err := ChunkAndEmbedManyMedlineRetrievals(articles, useAi)
// 	if err != nil {
// 		return nil, fmt.Errorf("error chunking the multiple entries: %w", err)
// 	}

// 	uids := make(map[string]string)
// 	for _, chunk := range chunks {
// 		// Create the Dgraph mutation
// 		mutation, err := createDgraphMutation(chunk)
// 		if err != nil {
// 			return nil, fmt.Errorf("error creating mutation for chunk: %w", err)
// 		}

// 		// Execute the mutation
// 		chunkUids, err := AddChunkWithMutation(mutation)
// 		if err != nil {
// 			return nil, fmt.Errorf("error adding chunk to Dgraph: %w", err)
// 		}

// 		// Merge the UIDs
// 		for k, v := range chunkUids {
// 			uids[k] = v
// 		}
// 	}

// 	return uids, nil
// }

// // createDgraphMutation creates a properly formatted mutation for Dgraph
// func createDgraphMutation(chunk schemas.TextChunk) (string, error) {
// 	// Create the Dgraph-specific structure
// 	dgraphChunk := struct {
// 		Uid       string         `json:"uid,omitempty"`
// 		Type      string         `json:"dgraph.type,omitempty"`
// 		schemas.TextChunk
// 	}{
// 		Type:      "TextChunk",
// 		TextChunk: chunk,
// 	}

// 	// Convert timestamp to RFC3339 format
// 	if !chunk.Metadata.Timestamp.IsZero() {
// 		dgraphChunk.Metadata.Timestamp = chunk.Metadata.Timestamp.Format(time.RFC3339)
// 	}

// 	// Marshal to JSON
// 	data, err := json.Marshal(dgraphChunk)
// 	if err != nil {
// 		return "", fmt.Errorf("error marshaling chunk: %w", err)
// 	}

// 	return string(data), nil
// }

// // AddChunkWithMutation executes a Dgraph mutation
// func AddChunkWithMutation(mutation string) (map[string]string, error) {
// 	// Create the mutation request
// 	request := &dgraph.Request{
// 		Mutations: []*dgraph.Mutation{
// 			{
// 				SetJson: mutation,
// 			},
// 		},
// 		CommitNow: true,
// 	}

// 	// Execute the mutation
// 	response, err := dgraph.Execute(dgraphConn, request)
// 	if err != nil {
// 		return nil, fmt.Errorf("error executing Dgraph mutation: %w", err)
// 	}

// 	return response.Uids, nil
// }

// // AddChunksInBatch adds multiple chunks in a single transaction
// func AddChunksInBatch(chunks []schemas.TextChunk) (map[string]string, error) {
// 	var mutations []string
// 	for _, chunk := range chunks {
// 		mutation, err := createDgraphMutation(chunk)
// 		if err != nil {
// 			return nil, fmt.Errorf("error creating mutation for chunk: %w", err)
// 		}
// 		mutations = append(mutations, mutation)
// 	}

// 	// Create batch mutation request
// 	var dgraphMutations []*dgraph.Mutation
// 	for _, mut := range mutations {
// 		dgraphMutations = append(dgraphMutations, &dgraph.Mutation{
// 			SetJson: mut,
// 		})
// 	}

// 	request := &dgraph.Request{
// 		Mutations: dgraphMutations,
// 		CommitNow: true,
// 	}

// 	// Execute batch mutation
// 	response, err := dgraph.Execute(dgraphConn, request)
// 	if err != nil {
// 		return nil, fmt.Errorf("error executing batch mutation: %w", err)
// 	}

// 	return response.Uids, nil
// }

// // QueryChunks retrieves chunks based on criteria
// func QueryChunks(userId string, keywords []string) ([]schemas.TextChunk, error) {
// 	// Construct the DQL query
// 	query := `
// 		query chunks($userId: string, $keywords: string) {
// 			chunks(func: type(TextChunk)) @filter(
// 				eq(userId, $userId) AND anyofterms(keywords, $keywords)
// 			) {
// 				id
// 				userId
// 				content
// 				score
// 				metadata {
// 					expand(_all_)
// 					medlineData {
// 						expand(_all_)
// 						authors {
// 							expand(_all_)
// 						}
// 					}
// 				}
// 				relations {
// 					expand(_all_)
// 				}
// 			}
// 		}
// 	`

// 	// Create variables map
// 	vars := map[string]string{
// 		"$userId":   userId,
// 		"$keywords": strings.Join(keywords, " "),
// 	}

// 	// Execute query
// 	resp, err := dgraph.Execute(dgraphConn, &dgraph.Request{
// 		Query: query,
// 		Vars:  vars,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error executing query: %w", err)
// 	}

// 	// Parse response
// 	var result struct {
// 		Chunks []schemas.TextChunk `json:"chunks"`
// 	}
// 	if err := json.Unmarshal(resp.Json, &result); err != nil {
// 		return nil, fmt.Errorf("error unmarshaling response: %w", err)
// 	}

// 	return result.Chunks, nil
// }
