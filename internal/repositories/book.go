package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/rs/zerolog"
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
