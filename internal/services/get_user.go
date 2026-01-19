package services

import (
	"context"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func GetUserById(ctx context.Context, userId uint) (*dtos.OwnProfileResponse, error) {
	// Get user object
	userObj, getUserErr := repositories.GetUserByID(ctx, userId)
	if getUserErr != nil {
		return nil, getUserErr
	}

	// Get library stats
	countBooks, getLibStatsErr := repositories.CountUserLibraryBooks(ctx, userId)
	if getLibStatsErr != nil {
		return nil, getLibStatsErr
	}

	// Get posts count
	countPosts, countPostsErr := repositories.CountUserPosts(ctx, userId)
	if countPostsErr != nil {
		return nil, countPostsErr
	}

	user := dtos.OwnProfileResponse{
		ProfileResponse: dtos.ProfileResponse{
			UserID:               userId,
			Username:             userObj.Username,
			ProfilePicture:       userObj.Pfp,
			AccountCreatedAtYear: uint(userObj.CreatedAt.Year()),
			BooksInLibrary:       uint(countBooks),
			Posts:                uint(countPosts),
		},
		Email: userObj.Email,
	}

	return &user, nil
}
