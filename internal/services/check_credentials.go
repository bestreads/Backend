package services

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

// CheckCredentials validates the user's password against the stored hash.
// It returns true and the user's ID if the credentials are valid.
// If the user is not found or the password is incorrect, it returns false.
func CheckCredentials(ctx context.Context, credentials dtos.LoginRequest) (bool, uint, error) {
	// Get user credentials from db
	user, getUserErr := repositories.GetUserByEmail(ctx, credentials.Email)
	if getUserErr != nil {
		return false, 0, getUserErr
	}

	// Validate password using hash
	match, _, hashValidationErr := argon2id.CheckHash(credentials.Password, user.Password_hash)

	return match, user.ID, hashValidationErr
}
