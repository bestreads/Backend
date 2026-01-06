package services

import (
	"context"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

const CONTENT_LIMIT uint = 10

func GetActivity(ctx context.Context, uids []uint) (dtos.ActivityResponse, error) {
	posts, err := repositories.GetPostActivity(ctx, uids, CONTENT_LIMIT)
	if err != nil {
		return dtos.ActivityResponse{}, err
	}

	converted, err := ConvertPost(posts)
	if err != nil {
		return dtos.ActivityResponse{}, err
	}

	activity, err := repositories.GetBookActivity(ctx, uids, CONTENT_LIMIT)
	if err != nil {
		return dtos.ActivityResponse{}, err
	}

	return dtos.ActivityResponse{Posts: converted, Activity: activity}, nil

}
