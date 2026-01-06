package services

import (
	"context"
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetPost(c context.Context, uid uint, bid uint) ([]dtos.PostResponse, error) {
	posts, err := repositories.GetDbPost(c, uid, bid)
	if err != nil {
		return []dtos.PostResponse{}, err
	}

	return ConvertPost(posts)
}

// gerade wird das bild noch automatisch aus dem storage wieder zurückgeholt,
// in zukunft vllt durch eine "data" api oder so
//
// man müsste das hier eig auch wegmachen und durch verschachtelte structs
// (sachen mit FKs in der db) machen
func ConvertPost(p []database.Post) ([]dtos.PostResponse, error) {
	res := make([]dtos.PostResponse, len(p))
	for i, post := range p {
		imageData, err := database.FileRetrieve(post.ImageHash, database.PostImage)
		if err != nil {
			return make([]dtos.PostResponse, 0), err
		}

		res[i] = dtos.PostResponse{
			Pfp:      post.User.Pfp,
			Username: post.User.Username,
			Uid:      post.User.ID,
			Book:     post.Book,
			Content:  post.Content,
			Image:    imageData,
		}
	}

	return res, nil
}

func CreatePost(c context.Context, id uint, bid uint, content string, b64i string) error {
	// leeres bild wird einfach "0", das ist okay glaube ich
	hash, err := database.FileStore(b64i, database.PostImage)
	if err != nil {
		return err
	}

	post := database.Post{
		UserID:    id,
		BookID:    bid,
		Content:   content,
		ImageHash: strconv.Itoa(hash),
	}

	return repositories.CreateDbPost(c, post)

}
