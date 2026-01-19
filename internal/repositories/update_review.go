package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func UpdateReview(ctx context.Context, userId, bookId, rating uint) error {
	db := middlewares.DB(ctx)
	_, updateReviewErr := gorm.G[database.Library](db).
		Where("user_id = ? AND book_id = ?", userId, bookId).
		Update(ctx, "Rating", rating)
	return updateReviewErr
}
