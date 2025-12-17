package handlers

import (
	"errors"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	validate := middlewares.Validator(ctx)

	// Parse request payload
	requestPayload := new(dtos.CreateUserRequest)
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
				// Check if format err is from email field
				if e.Field() == "Email" {
					switch e.Tag() {
					case "required":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "Email address is missing",
							})
					case "email":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "Invalid email format",
							})
					}
				}

				// Check if format err is from password field
				if e.Field() == "Password" {
					switch e.Tag() {
					case "required":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "Password is missing",
							})
					case "min":
						return c.Status(fiber.StatusBadRequest).
							JSON(dtos.GenericRestErrorResponse{
								Description: "Password must be at least 12 characters long",
							})
					}
				}
			}
		}

		// Fallback for other errs
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Validation failed",
			})
	}

	// Create user and retrieve user id
	userId, createUserErr := services.CreateUser(ctx, *requestPayload)
	if createUserErr != nil {
		if errors.Is(createUserErr, gorm.ErrDuplicatedKey) {
			msg := "User already exists"
			log.Debug().Err(createUserErr).Msg(msg)
			return c.Status(fiber.StatusConflict).
				JSON(dtos.GenericRestErrorResponse{
					Description: msg,
				})
		}

		msg := "Failed to create user"
		log.Error().Err(createUserErr).Msg(msg)

		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"userId": userId,
		})
}
