package services

import (
	"context"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetFollowingPostsFeed(ctx context.Context, userId uint, offset int) ([]dtos.PostResponse, error) {
	return repositories.GetFollowingPostsFeed(ctx, userId, offset)
}
