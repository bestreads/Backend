package handlers

import (
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBook(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	var book database.Book

	bid, err := strconv.ParseUint(c.Params("bid"), 10, 32)
	if err != nil {
		log.Error().Err(err).Str("bid", c.Params("bid")).Msg("invalid book ID format")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book ID must be a valid positive number",
			})
	}

	if bid == 0 {
		log.Error().Msg("book ID cannot be zero")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book ID must be greater than 0",
			})
	}

	err = middlewares.DB(ctx).Where("id = ?", bid).First(&book).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
