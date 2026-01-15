package repositories

import (
	"context"
	"strings"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
)

func SearchBooks(ctx context.Context, query string, limit int) ([]database.Book, error) {
	var books []database.Book

	query = strings.TrimSpace(query)
	if query == "" {
		return books, nil
	}

	pattern := "%" + strings.ToLower(query) + "%"

	err := middlewares.DB(ctx).
		Where("LOWER(title) LIKE ? OR LOWER(author) LIKE ?", pattern, pattern).
		Limit(limit).
		Find(&books).Error

	return books, err
}
