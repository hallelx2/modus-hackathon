package dg

import (
	"encoding/json"
	"fmt"

	"my-modus-app/src/schemas"

	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
	"github.com/hypermodeinc/modus/sdk/go/pkg/postgresql"
)

const hostName = "dg"
const databaseName = "postgres"

func AddUserAsJSON(user *schemas.User) (map[string]string, error) {
	// Marshal user struct into JSON
	data, err := json.Marshal(user)
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

func AddUserToDatabase(user schemas.User) (*schemas.User, error) {
	const query = `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Execute the query and retrieve the generated ID
	id, _, err := postgresql.QueryScalar[string](
		databaseName,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to add user to database: %v", err)
	}

	// Assign the ID to the user object and return it
	user.ID = id
	return &user, nil
}
