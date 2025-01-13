package user

import (
	"errors"
	"fmt"
	"my-modus-app/src/schemas"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)



func Login(email, password string) (*schemas.LoginUser, error) {
    // Query the user
    user, err := QueryUserByEmail(email)
    if err != nil {
        return nil, fmt.Errorf("error querying user: %w", err)
    }
    if user == nil {
        return nil, ErrUserNotFound
    }

    // Verify password using the provided VerifyPassword function
    err = VerifyPassword(user.Password, password)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    return &schemas.LoginUser{
		ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
	}, nil
}
