package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
	"gorm.io/gorm"
)

func GetPost(c context.Context, userId uint, offset int) ([]dtos.PostResponse, error) {
	posts, err := repositories.GetPosts(c, userId, offset)
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
			ProfilePicture: post.User.ProfilePicture,
			Username:       post.User.Username,
			Uid:            post.User.ID,
			Book:           post.Book,
			Content:        post.Content,
			CreatedAt:      post.UpdatedAt,
			State:          post.State,
			Rating:         post.Rating,
		}
	}

	return res, nil
}

func CreatePost(c context.Context, id uint, bid uint, content string) error {
	// leeres bild wird einfach "0", das ist okay glaube ich
	post := database.Post{
		UserID:    id,
		BookID:    bid,
		Content:   content,
		DeletedAt: gorm.DeletedAt{}, // das deletedAt feld explizit zurücksetzen
	}

	return repositories.CreateDbPost(c, post)

}

func DeletePost(ctx context.Context, uid uint, bid uint) error {
	// idk vllt hiermit noch was machen
	_, err := repositories.DeleteDbPost(ctx, uid, bid)
	return err
}
