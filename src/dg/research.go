package dg

import (
	"encoding/json"
	"fmt"
	"my-modus-app/src/schemas"

	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
)



func AddResearchAsJSON(research *schemas.Research) (map[string]string, error) {
	// Marshal user struct into JSON
	data, err := json.Marshal(research)
	if err != nil {
		return nil, fmt.Errorf("error marshaling user to JSON: %w", err)
	}

	// Execute the Dgraph mutation
	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Mutations: []*dgraph.Mutation{
			{
				SetJson: string(data),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error executing Dgraph mutation: %w", err)
	}

	return response.Uids, nil
}
