package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, email, passwordHash string) (*uint, error) {
	db := middlewares.DB(ctx)

	user := database.User{
		Email:         email,
		Password_hash: passwordHash,
	}

	createUserErr := gorm.G[database.User](db).
		Create(ctx, &user)

	if createUserErr != nil {
		return nil, createUserErr
	}

	return &user.ID, nil
}
