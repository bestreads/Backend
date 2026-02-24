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

	// Get followers count
	countFollowers, countFollowersErr := repositories.CountFollowers(ctx, userId)
	if countFollowersErr != nil {
		return nil, countFollowersErr
	}

	// Get following count
	countFollowing, countFollowingErr := repositories.CountFollowing(ctx, userId)
	if countFollowingErr != nil {
		return nil, countFollowingErr
	}

	user := dtos.OwnProfileResponse{
		ProfileResponse: dtos.ProfileResponse{
			UserID:               userId,
			Username:             userObj.Username,
			ProfilePicture:       userObj.ProfilePicture,
			AccountCreatedAtYear: uint(userObj.CreatedAt.Year()),
			BooksInLibrary:       uint(countBooks),
			Posts:                uint(countPosts),
			FollowersCount:       uint(countFollowers),
			FollowingCount:       uint(countFollowing),
		},
		Email: userObj.Email,
	}

	return &user, nil
}
