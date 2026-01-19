package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetPost(c context.Context, uid uint, limit int) ([]dtos.PostResponse, error) {
	posts, err := repositories.GetPost(c, uid, limit)
	if err != nil {
		return []dtos.PostResponse{}, err
	}

	return convert(posts)
}

func GetGlobalPosts(ctx context.Context, limit int) ([]dtos.PostResponse, error) {
	posts, err := repositories.GetGlobalPosts(ctx, limit)
	if err != nil {
		return []dtos.PostResponse{}, err
	}

	return convert(posts)
}

// man müsste das hier eig auch wegmachen und durch verschachtelte structs
// (sachen mit FKs in der db) machen
func convert(p []database.Post) ([]dtos.PostResponse, error) {
	res := make([]dtos.PostResponse, len(p))
	for i, post := range p {
		res[i] = dtos.PostResponse{
			Pfp:      post.User.ProfilePicture,
			Username: post.User.Username,
			Uid:      post.User.ID,
			Book:     post.Book,
			Content:  post.Content,
		}
	}

	return res, nil
}

func CreatePost(c context.Context, id uint, bid uint, content string) error {
	// leeres bild wird einfach "0", das ist okay glaube ich
	post := database.Post{
		UserID:  id,
		BookID:  bid,
		Content: content,
	}

	return repositories.CreateDbPost(c, post)

}
