package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func Startfollow(ctx context.Context, this_id uint, other_id uint) error {
	return gorm.G[database.FollowRel](middlewares.DB(ctx)).
		Create(ctx, &database.FollowRel{UserID: this_id, FollowingID: other_id})
}

func Stopfollow(ctx context.Context, this_id uint, other_id uint) error {
	_, err := gorm.G[database.FollowRel](middlewares.DB(ctx)).
		Where("user_id = ? AND following_id = ?", this_id, other_id).
		Delete(ctx)
	return err
}
