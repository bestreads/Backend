package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// wichtige informationen:
//  1. den namen der preload dinger *RICHTIG* zu schreiben (auch groß- und kleinschreibung)
//  2. gorm wird (sporadisch) andere sachen preloaden, keine ahnung wieso?
//  3. eigentlich wollte ich eine api machen, mit der man nur die metadaten von einem user lädt.
//     würde ein bisschen hübscher in json aussehen, das habe ich aber mal nicht gemacht
func GetPosts(ctx context.Context, userId uint, offset int) ([]database.Post, error) {
	db := middlewares.DB(ctx)
	cfg := middlewares.Config(ctx)

	// Build query
	query := gorm.G[database.Post](db).
		Preload("User", nil).
		Preload("Book", nil).
		Order("updated_at DESC").
		Limit(cfg.PaginationSteps)

	// Set offset when given
	if offset != 0 {
		query = query.Offset(offset)
	}

	// Filter for user id when given
	if userId != 0 {
		query = query.Where("user_id = ?", userId)
	}

	// Get entries
	posts, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}

	return enrichPostsWithLibrary(ctx, posts)
}

// enrichPostsWithLibrary lädt für jeden Post die Library-Daten (State, Rating)
func enrichPostsWithLibrary(ctx context.Context, posts []database.Post) ([]database.Post, error) {
	for i := range posts {
		libraries, err := gorm.G[database.Library](middlewares.DB(ctx)).
			Where("user_id = ? AND book_id = ?", posts[i].UserID, posts[i].BookID).
			Find(ctx)
		if err != nil {
			return nil, err
		}
		if len(libraries) > 0 {
			posts[i].State = libraries[0].State
			posts[i].Rating = libraries[0].Rating
		}
	}
	return posts, nil
}

func CreateDbPost(ctx context.Context, post database.Post) error {
	err := middlewares.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "book_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"content"}),
	}).Create(&post)

	return err.Error
}

func DeleteDbPost(ctx context.Context, uid uint, bid uint) (int, error) {
	result := middlewares.DB(ctx).
		Unscoped().
		Where("user_id = ? AND book_id = ?", uid, bid).
		Delete(&database.Post{})

	return int(result.RowsAffected), result.Error
}
