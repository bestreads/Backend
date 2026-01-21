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
	cfg := middlewares.Config(ctx)
	httpClient := middlewares.HttpClient(ctx)

	// Get offset from optional query param
	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}

	// Parse offset param
	nOffset, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error parsing int")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "Bad offset",
		})
	}

	// Get query param
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

	log.Debug().Str("query", query).Str("offset", offset).Str("author", authorstr).Msg("Searching for books")

	// Search in the database
	books, err := repositories.SearchBooks(ctx, query, int(nOffset), author)
	if err != nil {
		log.Error().Err(err).Msg("Error searching in database")
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Error searching in database",
			})
	}

	// If fewer results than limit, search in Open Library and then re-query DB
	if len(books) < cfg.PaginationSteps {
		log.Debug().Int("localResults", len(books)).Int("limit", cfg.PaginationSteps).Msg("Not enough local results, searching Open Library API")

		// Re-query DB to get all books including newly added ones
		if err := services.SearchOpenLibrary(httpClient, ctx, query, cfg.PaginationSteps, author); err != nil {
			log.Error().Err(err).Msg("Error searching in Open Library")
			return c.Status(fiber.StatusInternalServerError).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Error searching in Open Library",
				})
		}

		books, err = repositories.SearchBooks(ctx, query, int(nOffset), author)
		if err != nil {
			log.Error().Err(err).Msg("Error searching in database after Open Library")
			return c.Status(fiber.StatusInternalServerError).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Error searching in database",
				})
		}
	}

	log.Debug().Int("results", len(books)).Msg("Book search completed")

	return c.Status(fiber.StatusOK).
		JSON(books)
}
