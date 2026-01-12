package handlers

import (
	"errors"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBook(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	bid, err := strconv.ParseUint(c.Params("bid"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("bid", c.Params("bid")).Msg("invalid book ID format")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book ID must be a valid positive number",
			})
	}

	if bid == 0 {
		log.Warn().Msg("book ID cannot be zero")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book ID must be greater than 0",
			})
	}

	book, err := repositories.GetBookFromDB(ctx, log, bid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Uint64("bid", bid).Msg("book not found")
			return c.Status(fiber.StatusNotFound).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Book not found",
				})
		}
		log.Error().Err(err).Uint64("bid", bid).Msg("database error while fetching book")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Internal server error",
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(book)
}
