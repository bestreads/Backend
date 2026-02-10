package handlers

import (
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func StartFollow(c *fiber.Ctx) error {
	return followinner(c, false)
}

func StopFollow(c *fiber.Ctx) error {
	return followinner(c, true)
}

func followinner(c *fiber.Ctx, unfollow bool) error {
	log := middlewares.Logger(c.UserContext())
	this_id, err := middlewares.User(c).GetId()
	if err != nil {
		return err
	}

	id_str := c.Params("id")
	if id_str == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "empty user id",
		})

	}

	other_id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		log.Err(err).Msg("failed to parse user id")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "failed to parse uid",
		})

	}

	if err := services.SetFollow(c.UserContext(), this_id, uint(other_id), unfollow); err != nil {
		log.Err(err).Msg("failed to update followers")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error updating followers",
		})

	}

	return c.Status(fiber.StatusOK).JSON(dtos.GenericRestResponse{
		Message: "ok",
	})

}
