package services

import (
	"context"
	"fmt"
	"runtime"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

var (
	hashingParams = argon2id.Params{
		Memory:      256 * 1024,
		Iterations:  3,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  32,
		KeyLength:   32,
	}
)

func CreateUser(ctx context.Context, user dtos.CreateUserRequest) (*uint, error) {
	// Hash password
	passwordHash, hashingErr := argon2id.CreateHash(user.Password, &hashingParams)
	if hashingErr != nil {
		err := fmt.Errorf("failed to hash the password: %w", hashingErr)
		return nil, err
	}

	// Create user entry in DB
	userId, createUserErr := repositories.CreateUser(ctx, user.Email, user.Username, passwordHash)
	if createUserErr != nil {
		err := fmt.Errorf("failed to insert user into db: %w", createUserErr)
		return nil, err
	}

	return userId, nil
}
