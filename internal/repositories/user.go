package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func GetUserByID(ctx context.Context, uid uint) (database.User, error) {
	var user database.User
	err := middlewares.DB(ctx).Where("id = ?", uid).First(&user).Error
	return user, err
}

func GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	db := middlewares.DB(ctx)
	user, getUserErr := gorm.G[database.User](db).Where("email = ?", email).First(ctx)
	return user, getUserErr
}

func CountUserLibraryBooks(ctx context.Context, uid uint) (int64, error) {
	var count int64
	err := middlewares.DB(ctx).Model(&database.Library{}).Where("user_id = ?", uid).Count(&count).Error
	return count, err
}

func CountUserPosts(ctx context.Context, uid uint) (int64, error) {
	var count int64
	err := middlewares.DB(ctx).Model(&database.Post{}).Where("user_id = ?", uid).Count(&count).Error
	return count, err
}
