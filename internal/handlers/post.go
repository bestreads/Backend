package handlers

import (
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("GET demopost")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"uid":     "1",
		"bid":     "1",
		"content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	})
}
