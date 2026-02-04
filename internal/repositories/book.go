package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func GetBookFromDB(ctx context.Context, log zerolog.Logger, bid uint64) (database.Book, error) {
	var book database.Book

	err := middlewares.DB(ctx).Where("id = ?", bid).First(&book).Error
	if err != nil {
		log.Debug().Err(err).Uint64("bid", bid).Msg("database query error")
		return book, err
	}

	return book, nil
}

func UpdateBookAvgRating(ctx context.Context, bookId uint, avgRating float32, count uint) error {
	db := middlewares.DB(ctx)

	_, updateBookAvgRatingErr := gorm.G[database.Book](db).
		Where("id = ?", bookId).
		Updates(ctx, database.Book{Rating: database.Rating{Avg: avgRating, Count: count}})

	return updateBookAvgRatingErr
}
