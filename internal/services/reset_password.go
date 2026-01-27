package services

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"gorm.io/gorm"
)

var ErrInvalidSecurityAnswer = errors.New("invalid security answer")

func ResetPassword(ctx context.Context, req dtos.ResetPasswordRequest) error {
	log := middlewares.Logger(ctx)

	user, err := repositories.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Err(err).Str("email", req.Email).Msg("user not found for password reset")
		}
		return err
	}

	// Verify security answer (case-insensitive comparison)
	if user.SecurityAnswer != req.SecurityAnswer {
		return ErrInvalidSecurityAnswer
	}

	// Hash the new password
	passwordHash, hashErr := argon2id.CreateHash(req.NewPassword, &hashingParams)
	if hashErr != nil {
		log.Warn().Err(hashErr).Msg("failed to hash new password")
		return err
	}

	// Update user's password hash
	user.Password_hash = passwordHash

	// Save the user
	if err := repositories.SaveUser(ctx, &user); err != nil {
		log.Warn().Err(err).Msg("failed to save user with new password")
	}

	return nil
}
