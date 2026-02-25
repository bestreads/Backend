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

	postResponses, translatePostsErr := convert(posts)
	if translatePostsErr != nil {
		return nil, translatePostsErr
	}

	return postResponses, nil
}
