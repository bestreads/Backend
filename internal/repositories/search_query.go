package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"gorm.io/gorm"
)

func SearchBooks(db *gorm.DB, ctx context.Context, query string) ([]database.Book, error) {
	var books []database.Book
	err := db.WithContext(ctx).Where("LOWER(title) LIKE LOWER(?) OR LOWER(author) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books).Error
	return books, err
}
