package handlers

import (
	"fmt"
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SaveFile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	cfg := middlewares.Config(ctx)

	c.Accepts("image/png")
	c.Accepts("image/webp")
	c.Accepts("image/jpg")

	if len(c.Body()) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "no data",
		})
	}

	hash, err := services.SaveFile(c.Body())
	if err != nil {
		log.Error().Err(err).Msg("file store error")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "it exploded",
		})
	}

	log.Info().Msg(fmt.Sprintf("using path %d", hash))

	url := fmt.Sprintf("%s://%s%s/v1/media/%d", cfg.ApiProtocol, cfg.ApiDomain, cfg.ApiBasePath, hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}

func GetFile(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())

	key, err := strconv.ParseUint(c.Params("KEY"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad key format",
		})
	}

	// 400iq path sanitizing: doppeltes casting sichert den pfad bestimmt :3
	data, err := database.FileRetrieve(strconv.FormatUint(key, 10))
	if err != nil {
		log.Error().Err(err).Msg("file retrieve error")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "it did not wor :(",
		})
	}

	return c.Status(fiber.StatusOK).Send(data)
}
