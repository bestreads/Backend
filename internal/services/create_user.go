package services

import (
	"context"
	"errors"
	"runtime"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

var (
	ErrHashingPassword = errors.New("failed to hash the password")

	hashingParams = argon2id.Params{
		Memory:      256 * 1024,
		Iterations:  2,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  32,
		KeyLength:   32,
	}
)

func CreateUser(ctx context.Context, user dtos.CreateUserRequest) (uint, error) {
	// Hash password
	passwordHash, hashingErr := argon2id.CreateHash(user.Password, &hashingParams)
	if hashingErr != nil {
		err := errors.Join(ErrHashingPassword, hashingErr)
		return 0, err
	}

	// Create user entry in DB
	userId, createUserErr := repositories.CreateUser(ctx, user.Email, passwordHash)

	return *userId, createUserErr
}
