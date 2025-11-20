package handlers

import (
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	log.Info().Msg("Health check called")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
