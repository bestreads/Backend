package handlers

import (
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
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

	limit := c.Query("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt <= 0 {
			log.Warn().Msg("Book search called with wrong limit")
			return c.Status(fiber.StatusBadRequest).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Query parameter 'limit' has to be a number > 0",
				})
		}
	}
	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Book search called without query parameter")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Query parameter 'q' is required",
			})
	}

	log.Info().Str("query", query).Msg("Searching for books")

	// Search in the database
	books, err := repositories.SearchBooks(db, ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Error searching in database")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Error searching in database",
			})
	}

	// If no results, search in Open Library
	if len(books) == 0 {
		log.Info().Msg("No local results found, searching Open Library API")
		books, err = services.SearchOpenLibrary(httpClient, ctx, query, limit)
		if err != nil {
			log.Error().Err(err).Msg("Error searching in Open Library")
			return c.Status(fiber.StatusInternalServerError).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Error searching in Open Library",
				})
		}
	}

	log.Info().Int("results", len(books)).Msg("Book search completed")

	return c.Status(fiber.StatusOK).
		JSON(books)
}
