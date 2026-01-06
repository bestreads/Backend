package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func GetPostActivity(ctx context.Context, uids []uint, limit uint) ([]database.Post, error) {
	return gorm.G[database.Post](middlewares.DB(ctx)).
		Limit(int(limit)).
		Preload("User", nil).
		Preload("Book", nil).
		Where("user_id IN ?", uids).
		Find(ctx)
}

func GetBookActivity(ctx context.Context, uids []uint, limit uint) ([]database.ReadState, error) {
	return gorm.G[database.ReadState](middlewares.DB(ctx)).
		Limit(int(limit)).
		Where("user_id IN ?", uids).
		Find(ctx)
}
