package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories/generated"
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

func GetFollowing(ctx context.Context, uid uint) ([]uint, error) {
	return generated.Query[database.FollowRel](middlewares.DB(ctx)).GetFollowingGen(ctx, uid)
}

func GetFollowers(ctx context.Context, uid uint) ([]uint, error) {
	return generated.Query[database.FollowRel](middlewares.DB(ctx)).GetFollowersGen(ctx, uid)
}

// magic gorm shit
// https://gorm.io/docs/the_generics_way.html#Code-Generator-Workflow
// tl;dr: ~/go/bin/gorm gen -i ./internal/repositories/follow.go -o internal/repositories/generated
type Query[T any] interface {
	// SELECT user_id FROM @@table WHERE following_id=@uid
	GetFollowersGen(ctx context.Context, uid uint) ([]uint, error)

	// SELECT following_id FROM @@table WHERE user_id=@uid
	GetFollowingGen(ctx context.Context, uid uint) ([]uint, error)
}
