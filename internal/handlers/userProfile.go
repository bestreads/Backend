package handlers

import (
	"errors"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserProfile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// get and validate userID
	uid, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Params("id")).Msg("invalid user ID format")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "User ID must be a valid positive number",
			})
	}

	if uid == 0 {
		log.Warn().Msg("user ID cannot be zero")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "User ID must be greater than 0",
			})
	}

	// get user
	user, err := repositories.GetUserByID(ctx, uint(uid))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Uint64("uid", uid).Msg("user not found")
			return c.Status(fiber.StatusNotFound).
				JSON(dtos.GenericRestErrorResponse{
					Description: "User not found",
				})
		}
		log.Error().Err(err).Uint64("uid", uid).Msg("database error while fetching user")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Internal server error",
			})
	}

	// get library stats
	countBooks, err := repositories.CountUserLibraryBooks(ctx, uint(uid))
	if err != nil {
		log.Error().Err(err).Uint64("uid", uid).Msg("database error while counting library books")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Internal server error",
			})
	}

	// get posts count
	countPosts, err := repositories.CountUserPosts(ctx, uint(uid))
	if err != nil {
		log.Error().Err(err).Uint64("uid", uid).Msg("database error while counting posts")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Internal server error",
			})
	}

	// return userProfile
	return c.JSON(dtos.ProfileResponse{
		UserID:               uid,
		Username:             user.Username,
		ProfilePicture:       user.Pfp,
		AccountCreatedAtYear: uint(user.CreatedAt.Year()),
		BooksInLibrary:       uint(countBooks),
		Posts:                uint(countPosts),
		Follower:             0,
		Following:            0,
	})
}
