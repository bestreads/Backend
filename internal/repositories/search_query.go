package repositories

import (
	"context"
	"strings"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
)

func SearchBooks(ctx context.Context, query string, offset int, searchAuthor bool) ([]database.Book, error) {
	cfg := middlewares.Config(ctx)
	db := middlewares.DB(ctx)

	var books []database.Book

	query = strings.TrimSpace(query)
	if query == "" {
		return books, nil
	}

	pattern := "%" + strings.ToLower(query) + "%"

	// Build db query
	dbQuery := db.
		Offset(offset).
		Limit(cfg.PaginationSteps)

	// Add relevant filter
	if searchAuthor {
		dbQuery = dbQuery.Where("LOWER(author) LIKE ?", pattern)
	} else {
		dbQuery = dbQuery.Where("LOWER(title) LIKE ?", pattern)
	}

	// Query books
	bookQueryErr := dbQuery.Find(&books).Error

	return books, bookQueryErr
}
