package handlers

import (
	"errors"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UpdateReview(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	validate := middlewares.Validator(ctx)

	// Parse request payload
	requestPayload := new(dtos.UpdateReviewRequest)
	if bodyParsingErr := c.BodyParser(requestPayload); bodyParsingErr != nil {
		log.Warn().Err(bodyParsingErr).Msg("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Invalid request body",
			})
	}

	// Validate payload
	if validationErr := validate.Struct(requestPayload); validationErr != nil {
		// Cast format err
		var valErrs validator.ValidationErrors
		if errors.As(validationErr, &valErrs) {
			for _, e := range valErrs {
				// Check if format err is from rating field
				if e.Field() == "Rating" {
					switch e.Tag() {
					case "required":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "Rating is missing",
							})
					case "min":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "The rating must be between 1 and 5",
							})
					case "max":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "The rating must be between 1 and 5",
							})
					}
				}
			}
		}

		// Fallback for other errs
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Validation failed",
			})
	}

	// Get user id from token
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Update review in db
	rowsAffected, updateReviewErr := repositories.UpdateReview(ctx, userId, requestPayload.BookID, requestPayload.Rating)
	if rowsAffected == 0 {
		msg := "No review found for the given book"
		log.Error().Err(updateReviewErr).Msg(msg)
		return c.Status(fiber.StatusNotFound).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}
	if updateReviewErr != nil {
		msg := "Failed to update review"
		log.Error().Err(updateReviewErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "Review updated successfully",
		})
}
