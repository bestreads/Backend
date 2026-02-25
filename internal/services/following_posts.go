package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetFollowingPostsFeed(ctx context.Context, userId uint, offset int) ([]dtos.PostResponse, error) {
	posts, getFollowingPostsErr := repositories.GetFollowingPostsFeed(ctx, userId, offset)
	if getFollowingPostsErr != nil {
		return nil, getFollowingPostsErr
	}

	postResponses, translatePostsErr := translatePostsToPostResponses(posts)
	if translatePostsErr != nil {
		return nil, translatePostsErr
	}

	return postResponses, nil
}

func translatePostsToPostResponses(posts []database.Post) ([]dtos.PostResponse, error) {
	res := make([]dtos.PostResponse, len(posts))
	for i, post := range posts {
		res[i] = dtos.PostResponse{
			ProfilePicture: post.User.ProfilePicture,
			Username:       post.User.Username,
			Uid:            post.User.ID,
			Book:           post.Book,
			Content:        post.Content,
			CreatedAt:      post.CreatedAt,
			State:          post.State,
			Rating:         post.Rating,
		}
	}

	return res, nil
}
