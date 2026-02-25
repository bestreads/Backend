package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
)

func GetFollowingPostsFeed(ctx context.Context, userId uint, offset int) ([]database.Post, error) {
	db := middlewares.DB(ctx)
	cfg := middlewares.Config(ctx)

	var posts []database.Post

	query := db.Model(&database.Post{}).
		Select("posts.*, libraries.state AS state, libraries.rating AS rating").
		Joins("LEFT JOIN libraries ON libraries.user_id = posts.user_id AND libraries.book_id = posts.book_id").
		Joins("INNER JOIN follow_rels ON follow_rels.following_id = posts.user_id").
		Preload("User").
		Preload("Book").
		Where("follow_rels.user_id = ?", userId).
		Order("posts.updated_at DESC").
		Limit(cfg.PaginationSteps)

	// Set offset when given
	if offset != 0 {
		query = query.Offset(offset)
	}

	// Get entries
	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
