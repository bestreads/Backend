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

// CheckUserUniqueness prüft, ob ein Benutzername oder eine E-Mail bereits von einem anderen Benutzer verwendet wird.
// Gibt true zurück, wenn Username/Email bereits existiert (Konflikt), false wenn verfügbar.
func CheckUserUniqueness(ctx context.Context, excludeUserID uint, username, email string) (bool, error) {
	if username == "" && email == "" {
		return false, nil
	}

	var count int64
	db := middlewares.DB(ctx)
	query := db.Model(&database.User{}).Where("id <> ?", excludeUserID)

	if username != "" && email != "" {
		query = query.Where("username = ? OR email = ?", username, email)
	} else if username != "" {
		query = query.Where("username = ?", username)
	} else {
		query = query.Where("email = ?", email)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// SaveUser speichert die Änderungen eines Users in der Datenbank.
func SaveUser(ctx context.Context, user *database.User) error {
	return middlewares.DB(ctx).Save(user).Error
}
