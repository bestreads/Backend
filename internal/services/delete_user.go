package services

import (
	"context"

	"github.com/bestreads/Backend/internal/repositories"
)

func DeleteUser(ctx context.Context, userId uint) error {
	return repositories.DeleteUser(ctx, userId)
}
