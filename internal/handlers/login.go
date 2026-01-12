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

func Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	validate := middlewares.Validator(ctx)

	// Parse request payload
	requestPayload := new(dtos.LoginRequest)
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

	// Check credentials
	match, checkCredentialsErr := services.CheckCredentials(ctx, *requestPayload)
	invalidCredentialsMsg := "Invalid email or password"
	if errors.Is(checkCredentialsErr, gorm.ErrRecordNotFound) {
		log.Debug().Err(checkCredentialsErr).Msg("Login failed: User not found")
		return c.Status(fiber.StatusUnauthorized).
			JSON(dtos.GenericRestErrorResponse{
				Description: invalidCredentialsMsg,
			})
	} else if checkCredentialsErr != nil {
		msg := "Failed to validate user credentials"
		log.Error().Err(checkCredentialsErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	if !match {
		log.Debug().Msg("Login failed: Password mismatch")
		return c.Status(fiber.StatusUnauthorized).
			JSON(dtos.GenericRestErrorResponse{
				Description: invalidCredentialsMsg,
			})
	}

	// ToDo: Generate access JWT and set as Cookie

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "Login successful",
		})
}
