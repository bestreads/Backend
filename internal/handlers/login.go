package handlers

import (
	"errors"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/bestreads/Backend/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Login handles the user login process.
// @Summary Login to an existing user account
// @Description Logs the user into the system with the provided credentials and returns a refresh-JWT and an access-JWT.
// @Tags User Management
// @Accept json
// @Produce json
// @Param User body dtos.LoginRequest true "The user's email address and password."
// @Success 200 {object} dtos.GenericRestResponse "Login successful"
// @Failure 400 {object} dtos.GenericRestErrorResponse "Invalid request (-body)"
// @Failure 401 {object} dtos.GenericRestErrorResponse "Wrong credentials"
// @Failure 500 {object} dtos.GenericRestErrorResponse "Internal server error"
// @Router /v1/auth/login [post]
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
	match, userId, checkCredentialsErr := services.CheckCredentials(ctx, *requestPayload)
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

	// Create access JWT
	accessJwt, accessJwtGenerationErr := services.GenerateToken(ctx, strconv.FormatUint(uint64(userId), 10), types.AccessToken)
	if accessJwtGenerationErr != nil {
		msg := "Failed to create access JWT"
		log.Error().Err(accessJwtGenerationErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Create refresh JWT
	refreshJwt, refreshJwtGenerationErr := services.GenerateToken(ctx, strconv.FormatUint(uint64(userId), 10), types.RefreshToken)
	if refreshJwtGenerationErr != nil {
		msg := "Failed to create refresh JWT"
		log.Error().Err(refreshJwtGenerationErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Create cookies
	accessTokenCookie := services.CreateCookie(ctx, types.AccessToken, accessJwt, true, false)
	refreshTokenCookie := services.CreateCookie(ctx, types.RefreshToken, refreshJwt, true, false)

	// Set tokens as cookies
	c.Cookie(accessTokenCookie)
	c.Cookie(refreshTokenCookie)

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "Login successful",
		})
}
