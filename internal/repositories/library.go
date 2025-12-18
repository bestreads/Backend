package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func AddBook(ctx context.Context, uid uint, bid uint, state database.ReadState) error {
	lib := database.Library{
		UserID: uid,
		BookID: bid,
		State:  state,
		Rating: 0,
	}

	return gorm.G[database.Library](middlewares.DB(ctx)).Create(ctx, &lib)
}

func QueryLibraryDb(ctx context.Context, uid uint, limit uint64) ([]database.Library, error) {
	// insane das hier random der user grepreloaded wird???
	return gorm.G[database.Library](middlewares.DB(ctx)).
		Preload("Book", nil).
		Limit(int(limit)).
		Where("user_id = ?", uid).
		Find(ctx)
}
