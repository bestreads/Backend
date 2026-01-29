package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Health is a simple ping endpoint to check if the server is running properly.
// @Summary      Returns a status ok response
// @Description  Simple ping endpoint to check if the server is running properly.
// @Tags         HealthCheck
// @Produce      json
// @Success 200  {object} dtos.GenericRestResponse "ok"
// @Router       /v1/health [get]
func Health(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	log.Info().Msg("Health check called")

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "ok",
		})
}
