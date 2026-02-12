package repositories

import (
	"context"
	"errors"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

func DeleteUser(ctx context.Context, userId uint) error {
	db := middlewares.DB(ctx)

	rowsAffected, err := gorm.G[database.User](db).Where("id = ?", userId).Delete(ctx)
	if err == nil && rowsAffected == 0 {
		err = errors.New("User could not be found")
	}

	return err
}
