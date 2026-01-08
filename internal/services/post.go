package services

import (
	"context"
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetPost(c context.Context, uid uint, bid uint, limit int) ([]dtos.PostResponse, error) {
	posts, err := repositories.GetPost(c, uid, bid, limit)
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
		// das hier brauchen wir bald nicht mehr
		// imageData, err := database.FileRetrieve(post.ImageHash, database.PostImage)
		// if err != nil {
		// 	return make([]dtos.PostResponse, 0), err
		// }

		res[i] = dtos.PostResponse{
			Pfp:      post.User.Pfp,
			Username: post.User.Username,
			Uid:      post.User.ID,
			Book:     post.Book,
			Content:  post.Content,
			ImageUrl: post.ImageUrl,
		}
	}

	return res, nil
}

func CreatePost(c context.Context, id uint, bid uint, content string, b64i string) error {
	// leeres bild wird einfach "0", das ist okay glaube ich
	hash, err := database.FileStoreB64(b64i, database.PostImage)
	if err != nil {
		return err
	}

	post := database.Post{
		UserID:   id,
		BookID:   bid,
		Content:  content,
		ImageUrl: strconv.Itoa(hash),
	}

	return repositories.CreateDbPost(c, post)

}
