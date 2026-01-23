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
func GetPost(ctx context.Context, uid uint, limit int) ([]database.Post, error) {
	query := gorm.G[database.Post](middlewares.DB(ctx)).
		Preload("User", nil).
		Preload("Book", nil).
		Where("user_id = ?", uid)
	if limit > 0 {
		query = query.Limit(limit)
	}
	posts, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}
	return enrichPostsWithLibrary(ctx, posts)
}

func GetGlobalPosts(ctx context.Context, limit int) ([]database.Post, error) {
	query := gorm.G[database.Post](middlewares.DB(ctx)).
		Preload("User", nil).
		Preload("Book", nil)
	if limit > 0 {
		query = query.Limit(limit)
	}
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
