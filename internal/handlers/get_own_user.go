package handlers

import (
	"errors"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOwnUser(c *fiber.Ctx) error {
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

	// Get user data for the given user id
	userObj, getUserErr := services.GetUserById(ctx, userId)
	if errors.Is(getUserErr, gorm.ErrRecordNotFound) {
		msg := "User could not be found"
		log.Error().Err(getUserErr).Uint64("user-id", uint64(userId)).Msg(msg)
		return c.Status(fiber.StatusNotFound).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	} else if getUserErr != nil {
		msg := "Failed to get user data"
		log.Error().Err(getUserErr).Uint64("user-id", uint64(userId)).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(userObj)
}
