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
	c.Accepts("image/png")
	c.Accepts("image/webp")
	c.Accepts("image/jpg")

	log := middlewares.Logger(c.UserContext())

	itype := c.Query("type")
	if itype == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "no image type",
		})
	}

	nitype, err := strconv.ParseInt(itype, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad type",
		})
	}

	if len(c.Body()) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "no data",
		})
	}

	// ...warum
	t, exists := database.ImageTypeMap[int(nitype)]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad image type",
		})

	}

	hash, err := services.SaveFile(c.Body(), t)
	log.Info().Msg(fmt.Sprintf("using type %d and path %d", t, hash))

	url := fmt.Sprintf("%s/api/v1/media/%d", middlewares.Config(c.UserContext()).ApiBaseURL, hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}
