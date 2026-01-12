package services

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func CheckCredentials(ctx context.Context, credentials dtos.LoginRequest) (bool, error) {
	// Get user credentials from db
	user, getUserErr := repositories.GetUserByEmail(ctx, credentials.Email)
	if getUserErr != nil {
		return false, getUserErr
	}

	// Validate password using hash
	match, _, hashValidationErr := argon2id.CheckHash(credentials.Password, user.Password_hash)

	return match, hashValidationErr
}
