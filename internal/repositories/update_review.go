package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func UpdateReview(ctx context.Context, userId, bookId, rating uint) (int, error) {
	db := middlewares.DB(ctx)
	rowsAffected, updateReviewErr := gorm.G[database.Library](db).
		Where("user_id = ? AND book_id = ?", userId, bookId).
		Update(ctx, "Rating", rating)
	return rowsAffected, updateReviewErr
}
