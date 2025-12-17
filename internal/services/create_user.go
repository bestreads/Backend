package services

import (
	"context"
	"runtime"

	"github.com/pkg/errors"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

var (
	hashingParams = argon2id.Params{
		Memory:      256 * 1024,
		Iterations:  2,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  32,
		KeyLength:   32,
	}
)

func CreateUser(ctx context.Context, user dtos.CreateUserRequest) (*uint, error) {
	// Hash password
	passwordHash, hashingErr := argon2id.CreateHash(user.Password, &hashingParams)
	if hashingErr != nil {
		err := errors.Wrap(hashingErr, "failed to hash the password")
		return nil, err
	}

	// Create user entry in DB
	userId, createUserErr := repositories.CreateUser(ctx, user.Email, passwordHash)
	if createUserErr != nil {
		err := errors.Wrap(createUserErr, "failed to insert user into db")
		return nil, err
	}

	return userId, nil
}
