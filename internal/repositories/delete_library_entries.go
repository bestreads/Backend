package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func DeleteLibraryEntries(ctx context.Context, userId uint) error {
	db := middlewares.DB(ctx)

	_, err := gorm.G[database.Library](db).Where("user_id = ?", userId).Delete(ctx)

	return err
}
