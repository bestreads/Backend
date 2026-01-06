package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func GetActivity(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("GET activities")

	pl := struct {
		Uids []uint `json:"uids"`
	}{}

	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("Json Parser Error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "Json Invalid",
		})
	}

	res, err := services.GetActivity(c.UserContext(), pl.Uids)
	if err != nil {
		log.Error().Err(err).Msg("internal error")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
