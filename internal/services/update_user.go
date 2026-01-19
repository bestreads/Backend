package services

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
)

func UpdateUser(ctx context.Context, userId uint, req dtos.UpdateUserRequest) error {
	// Hole aktuellen User aus DB
	user, err := repositories.GetUserByID(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update nur die Felder, die nicht leer sind
	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Pfp != "" {
		user.Pfp = req.Pfp
	}

	if req.Password != "" {
		// Hash das neue Passwort
		passwordHash, hashErr := argon2id.CreateHash(req.Password, &hashingParams)
		if hashErr != nil {
			return fmt.Errorf("failed to hash password: %w", hashErr)
		}
		user.Password_hash = passwordHash
	}

	// Speichere die Änderungen
	if err := middlewares.DB(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}

	return nil
}
