package services

import (
	"context"
	"errors"

	"github.com/bestreads/Backend/internal/repositories"
)

func DeleteUser(ctx context.Context, userId uint) error {
	deletePostsErr := repositories.DeletePosts(ctx, userId)
	deleteLibraryEntriesErr := repositories.DeleteLibraryEntries(ctx, userId)
	deleteUserErr := repositories.DeleteUser(ctx, userId)

	err := errors.Join(deletePostsErr, deleteLibraryEntriesErr, deleteUserErr)

	return err
}
