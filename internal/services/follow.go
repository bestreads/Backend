package services

import (
	"context"

	"github.com/bestreads/Backend/internal/repositories"
)

func SetFollow(ctx context.Context, this_id uint, other_id uint, unfollow bool) error {
	if unfollow {
		return repositories.Stopfollow(ctx, this_id, other_id)
	} else {
		return repositories.Startfollow(ctx, this_id, other_id)
	}

}
