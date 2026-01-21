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
	httpClient := middlewares.HttpClient(ctx)

	limit := c.Query("limit")
	if limit == "" {
		log.Warn().Msg("Book search called without limit parameter")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Query parameter 'limit' is required",
			})
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		log.Warn().Msg("Book search called with wrong limit")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Query parameter 'limit' has to be a number > 0",
			})
	}

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Book search called without query parameter")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Query parameter 'q' is required",
			})
	}

	var author bool

	authorstr := c.Query("author")
	if authorstr == "" {
		author = false
	} else {
		author = true
	}

	log.Info().Str("query", query).Str("limit", limit).Str("author", authorstr).Msg("Searching for books")

	// Search in the database
	books, err := repositories.SearchBooks(ctx, query, limitInt, author)
	if err != nil {
		log.Error().Err(err).Msg("Error searching in database")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Error searching in database",
			})
	}

	// If fewer results than limit, search in Open Library and then re-query DB
	if len(books) < limitInt {
		log.Info().Int("localResults", len(books)).Int("limit", limitInt).Msg("Not enough local results, searching Open Library API")
		err := services.SearchOpenLibrary(httpClient, ctx, query, limit, author)
		if err != nil {
			log.Error().Err(err).Msg("Error searching in Open Library")
			return c.Status(fiber.StatusInternalServerError).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Error searching in Open Library",
				})
		}

		// Re-query DB to get all books including newly added ones
		books, err = repositories.SearchBooks(ctx, query, limitInt, author)
		if err != nil {
			log.Error().Err(err).Msg("Error searching in database after Open Library")
			return c.Status(fiber.StatusInternalServerError).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Error searching in database",
				})
		}
	}

	log.Info().Int("results", len(books)).Msg("Book search completed")

	return c.Status(fiber.StatusOK).
		JSON(books)
}
