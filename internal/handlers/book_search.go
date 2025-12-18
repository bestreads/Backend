package handlers

import (
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func BookSearch(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	db := middlewares.DB(ctx)
	httpClient := middlewares.HttpClient(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Book search called without query parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	log.Info().Str("query", query).Msg("Searching for books")

	// Search in the database
	books, err := repositories.SearchBooks(db, ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Error searching in database")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error searching in database: " + err.Error(),
		})
	}

	// If no results, search in Open Library
	if len(books) == 0 {
		log.Info().Msg("No local results found, searching Open Library API")
		books, err = services.SearchOpenLibrary(httpClient, ctx, query)
		if err != nil {
			log.Error().Err(err).Msg("Error searching in Open Library")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error searching in Open Library: " + err.Error(),
			})
		}
	}

	log.Info().Int("results", len(books)).Msg("Book search completed")

	return c.Status(fiber.StatusOK).JSON(books)
}
