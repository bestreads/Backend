package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func DeleteUser(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Get user id
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Delete user from db
	if err := services.DeleteUser(ctx, userId); err != nil {
		msg := "Failed to delete the user"
		log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "User deleted successfully",
		})
}
