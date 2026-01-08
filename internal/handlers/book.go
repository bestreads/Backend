package handlers

import (
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func GetBook(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	var book database.Book

	bid, err := strconv.ParseUint(c.Params("BID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing bookid error")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Invalid book ID format",
			})
	}

	if bid == 0 {
		log.Error().Msg("bid is not valid")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book ID must be greater than 0",
			})
	}

	err = middlewares.DB(ctx).Where("id = ?", bid).First(&book).Error
	if err != nil {
		log.Error().Err(err).Msg("error finding book in db")
		return c.Status(fiber.StatusNotFound).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Book not found",
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(book)
}
