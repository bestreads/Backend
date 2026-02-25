package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
)

func GetFollowingPostsFeed(ctx context.Context, userId uint, offset int) ([]dtos.PostResponse, error) {
	db := middlewares.DB(ctx)
	cfg := middlewares.Config(ctx)

	var posts []dtos.PostResponse

	query := db.Model(&database.Post{}).
		Select("users.profile_picture AS profile_picture, users.username AS username, users.id AS uid, posts.book_id AS book_id, posts.content AS content, posts.created_at AS created_at, libraries.state AS state, libraries.rating AS rating, books.*").
		Joins("INNER JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN libraries ON libraries.user_id = posts.user_id AND libraries.book_id = posts.book_id").
		Joins("LEFT JOIN books ON books.id = posts.book_id").
		Joins("INNER JOIN follow_rels ON follow_rels.following_id = posts.user_id").
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
