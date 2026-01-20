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

func QueryLibraryDb(ctx context.Context, uid uint, limit int64) ([]database.Library, error) {
	// insane das hier random der user grepreloaded wird???
	return gorm.G[database.Library](middlewares.DB(ctx)).
		Preload("Book", nil).
		Limit(int(limit)).
		Where("user_id = ?", uid).
		Find(ctx)
}

func ReadLibrariesForBook(ctx context.Context, bookId uint, filterZeroRatings bool) ([]database.Library, error) {
	db := middlewares.DB(ctx)

	query := gorm.G[database.Library](db).
		Where("book_id = ?", bookId)

	if filterZeroRatings {
		query = query.Where("rating != 0")
	}
	libraries, librariesReadErr := query.Find(ctx)

	return libraries, librariesReadErr
}

func UpdateReadState(ctx context.Context, uid uint, bid uint, state database.ReadState) (int, error) {
	count, err := gorm.G[database.Library](middlewares.DB(ctx)).
		Where("user_id = ? AND book_id = ?", uid, bid).
		Update(ctx, "state", state)
	return count, err
}

func DeleteFromLibrary(ctx context.Context, uid uint, bid uint) (int, error) {
	return gorm.G[database.Library](middlewares.DB(ctx)).Where("user_id = ? AND book_ID = ?", uid, bid).Delete(ctx)
}
