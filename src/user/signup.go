package user

import (
	"errors"
	"fmt"
	"my-modus-app/src/dg"
	"my-modus-app/src/schemas"
	"time"

	"github.com/google/uuid"
)

var ErrUserExists = errors.New("user already exists with this email")

func Signup(email, name, password string) (map[string]string, error) {
	// Check if user already exists
	exists, err := CheckUserExists(email)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing the password: %w", err)
	}

	// Create a new user object
	user := schemas.User{
		ID:        uuid.NewString(),
		Name:      name,
		Password:  hashedPassword,
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DType:     []string{"User"}, // Important: Set the type for Dgraph
	}

	// Add user to the database
	uids, err := dg.AddUserAsJSON(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user to database: %w", err)
	}

	return uids, nil
}

func SignupWithDatabase(email, name, password string) (*schemas.User, error) {
	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing the password: %w", err)
	}

	// Create a new user object
	user := schemas.User{
		ID:        uuid.NewString(),
		Name:      name,
		Password:  hashedPassword,
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Add user to the database
	addedUser, err := dg.AddUserToDatabase(user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user to database: %w", err)
	}

	return addedUser, nil
}
