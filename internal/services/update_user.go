package services

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
)

func UpdateUser(ctx context.Context, userId uint, req dtos.UpdateUserRequest) error {
	cfg := middlewares.Config(ctx)
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

	if len(req.Pfp) > 0 {
		// Speichere das Bild im Filestore
		hash, err := database.FileStoreRaw(req.Pfp)
		if err != nil {
			return fmt.Errorf("failed to store profile picture: %w", err)
		}
		// Generiere den Link zum Bild
		url := fmt.Sprintf("%s://%s:%s%s/v1/media/%d", cfg.ApiProtocol, cfg.ApiDomain, cfg.ApiPort, cfg.ApiBasePath, hash)
		user.Pfp = url
	}

	if req.Password != "" {
		// Hash das neue Passwort
		passwordHash, hashErr := argon2id.CreateHash(req.Password, &hashingParams)
		if hashErr != nil {
			return fmt.Errorf("failed to hash password: %w", hashErr)
		}
		user.Password_hash = passwordHash
	}

	db := middlewares.DB(ctx)

	// Prüfe, ob die neue E-Mail oder der neue Benutzername bereits von einem anderen Benutzer verwendet wird
	if req.Username != "" || req.Email != "" {
		var count int64
		if err := db.Model(&user).
			Where("id <> ? AND (username = ? OR email = ?)", user.ID, user.Username, user.Email).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check user uniqueness: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("username or email already in use")
		}
	}

	// Speichere die Änderungen
	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
