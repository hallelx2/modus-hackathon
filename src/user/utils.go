package user

import (
	"encoding/json"
	"fmt"
	"my-modus-app/src/schemas"

	"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
	"golang.org/x/crypto/bcrypt"
)

const hostName = "dg"

func HashPassword(password string) (string, error) {
	// Generate a hashed password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Return the hashed password as a string
	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, plainPassword string) error {
	// Compare the hashed password with the plain text password
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}


// QueryUserByEmail searches for a user by email address
// Returns the user if found, nil if not found, or error if query fails
func QueryUserByEmail(email string) (*schemas.User, error) {
    statement := `
    query queryUser($email: string) {
        users(func: eq(User.email, $email)) {
            uid
            User.id
            User.name
            User.email
            User.password
            User.created_at
            User.updated_at
            dgraph.type
        }
    }
    `
    variables := map[string]string{
        "$email": email,
    }

    response, err := dgraph.Execute(hostName, &dgraph.Request{
        Query: &dgraph.Query{
            Query:     statement,
            Variables: variables,
        },
    })
    if err != nil {
        return nil, fmt.Errorf("error executing Dgraph query: %w", err)
    }

    var result struct {
        Users []schemas.User `json:"users"`
    }

    if err := json.Unmarshal([]byte(response.Json), &result); err != nil {
        return nil, fmt.Errorf("error unmarshaling response: %w", err)
    }

    if len(result.Users) == 0 {
        return nil, nil // User not found
    }

    return &result.Users[0], nil
}

// CheckUserExists checks if a user with the given email already exists
// Returns true if user exists, false if not, or error if query fails
func CheckUserExists(email string) (bool, error) {
    user, err := QueryUserByEmail(email)
    if err != nil {
        return false, err
    }
    return user != nil, nil
}
