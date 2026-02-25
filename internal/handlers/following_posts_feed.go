package handlers

import (
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func GetFollowingPostsFeed(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Get user id from token
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Get offset from optional query param
	offsetString := c.Query("offset")
	if offsetString == "" {
		offsetString = "0"
	}

	// Parse offset param
	offset, parseOffsetErr := strconv.ParseInt(offsetString, 10, 32)
	if parseOffsetErr != nil {
		msg := "Failed to parse offset"
		log.Error().Err(parseOffsetErr).Msg(msg)
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Get posts for the feed
	posts, getFollowingPostsErr := services.GetFollowingPostsFeed(ctx, userId, int(offset))
	if getFollowingPostsErr != nil {
		msg := "Failed to get following posts feed"
		log.Error().Err(getFollowingPostsErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(posts)
}
