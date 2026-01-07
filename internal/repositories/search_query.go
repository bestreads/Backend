package repositories

import (
	"context"
	"strings"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
)

func SearchBooks(ctx context.Context, query string, limit int) ([]database.Book, error) {
	var books []database.Book

	// Split query in einzelne Wörter
	words := strings.Fields(strings.TrimSpace(query))
	if len(words) == 0 {
		return books, nil
	}

	// jedes Wort wird in title, author und description gesucht
	dbQuery := middlewares.DB(ctx)

	for i, word := range words {
		pattern := "%" + strings.ToLower(word) + "%"
		condition := "LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(description) LIKE ?"

		if i == 0 {
			dbQuery = dbQuery.Where(condition, pattern, pattern, pattern)
		} else {
			dbQuery = dbQuery.Or(condition, pattern, pattern, pattern)
		}
	}

	err := dbQuery.Limit(limit).Find(&books).Error
	return books, err
}
