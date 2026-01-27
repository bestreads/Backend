package handlers

import (
	"errors"
	"io"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserProfile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// get and validate userID
	userId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Params("id")).Msg("invalid user ID format")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "User ID must be a valid positive number",
			})
	}

	if userId == 0 {
		log.Warn().Msg("user ID cannot be zero")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "User ID must be greater than 0",
			})
	}

	// Get user data for the given user id
	userObj, getUserErr := services.GetUserById(ctx, uint(userId))
	if errors.Is(getUserErr, gorm.ErrRecordNotFound) {
		msg := "User could not be found"
		log.Error().Err(getUserErr).Uint64("user-id", uint64(userId)).Msg(msg)
		return c.Status(fiber.StatusNotFound).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	} else if getUserErr != nil {
		msg := "Failed to get user data"
		log.Error().Err(getUserErr).Uint64("user-id", uint64(userId)).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(userObj.ProfileResponse)
}

func ChangeUserData(c *fiber.Ctx) error {
	userFromMiddleware := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	validate := middlewares.Validator(ctx)

	// Get user id from token
	userId, err := userFromMiddleware.GetId()
	if err != nil {
		msg := "Failed to get user id"
		log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Parse request payload
	var payload dtos.UpdateUserRequest
	payload.Email = c.FormValue("email")
	payload.Username = c.FormValue("username")
	payload.Password = c.FormValue("password")

	// Handle Profilbild-Upload
	file, err := c.FormFile("profile_picture")
	if err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) && err.Error() != "there is no uploaded file associated with the given key" {
		log.Warn().Err(err).Msg("failed to parse uploaded profile picture")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Failed to parse uploaded file",
			})
	}
	if file != nil {
		// Maximale Upload-Größe: 5 MB
		const maxUploadSize = 5 << 20 // 5 MB
		if file.Size > maxUploadSize {
			log.Warn().Int64("size", file.Size).Msg("uploaded file too large")
			return c.Status(fiber.StatusRequestEntityTooLarge).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Profile picture must be smaller than 5 MB",
				})
		}

		// Öffne die Datei
		openedFile, openErr := file.Open()
		if openErr != nil {
			log.Warn().Err(openErr).Msg("failed to open uploaded file")
			return c.Status(fiber.StatusBadRequest).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Failed to process uploaded file",
				})
		}
		defer openedFile.Close()

		fileData, err := io.ReadAll(io.LimitReader(openedFile, maxUploadSize+1))
		if err != nil {
			log.Warn().Err(err).Msg("failed to read uploaded file")
			return c.Status(fiber.StatusBadRequest).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Failed to read uploaded file",
				})
		}
		payload.ProfilePicture = fileData
	}

	// Prüfe ob mindestens ein Feld gesetzt ist
	if payload.IsEmpty() {
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "At least one field must be provided",
			})
	}

	// Validierung der Felder
	if validationErr := validate.Struct(&payload); validationErr != nil {
		var valErrs validator.ValidationErrors
		if errors.As(validationErr, &valErrs) {
			for _, e := range valErrs {
				if e.Field() == "Email" && e.Tag() == "email" {
					return c.Status(fiber.StatusBadRequest).
						JSON(dtos.GenericRestErrorResponse{
							Description: "Invalid email format",
						})
				}

				if e.Field() == "Password" && e.Tag() == "min" {
					return c.Status(fiber.StatusBadRequest).
						JSON(dtos.GenericRestErrorResponse{
							Description: "Password must be at least 12 characters long",
						})
				}
			}
		}

		// Fallback
		log.Warn().Err(validationErr).Msg("validation failed")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Validation failed",
			})
	}

	// Business Logic im Service
	if err := services.UpdateUser(ctx, uint(userId), payload); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, services.ErrUserConflict) {
			return c.Status(fiber.StatusConflict).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Username or email already in use",
				})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Failed to update user",
			})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ResetUserPassword(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)
	validate := middlewares.Validator(ctx)

	var payload dtos.ResetPasswordRequest
	if err := c.BodyParser(&payload); err != nil {
		log.Warn().Err(err).Msg("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Invalid request body",
			})
	}

	// Validate request
	if validationErr := validate.Struct(&payload); validationErr != nil {
		var valErrs validator.ValidationErrors
		if errors.As(validationErr, &valErrs) {
			for _, e := range valErrs {
				if e.Field() == "Email" && e.Tag() == "email" {
					return c.Status(fiber.StatusBadRequest).
						JSON(dtos.GenericRestErrorResponse{
							Description: "Invalid email format",
						})
				}

				if e.Field() == "NewPassword" && e.Tag() == "min" {
					return c.Status(fiber.StatusBadRequest).
						JSON(dtos.GenericRestErrorResponse{
							Description: "Password must be at least 12 characters long",
						})
				}

				if e.Field() == "SecurityAnswer" && e.Tag() == "required" {
					return c.Status(fiber.StatusBadRequest).
						JSON(dtos.GenericRestErrorResponse{
							Description: "Security answer is required",
						})
				}
			}
		}

		// Fallback
		log.Warn().Err(validationErr).Msg("validation failed")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Validation failed",
			})
	}

	// Call service to reset password
	if err := services.ResetPassword(ctx, payload); err != nil {
		log.Warn().Err(err).Str("email", payload.Email).Msg("password reset failed")
		if errors.Is(err, services.ErrInvalidSecurityAnswer) {
			return c.Status(fiber.StatusUnauthorized).
				JSON(dtos.GenericRestErrorResponse{
					Description: "Invalid email or security answer",
				})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: "Failed to reset password",
			})
	}

	return c.SendStatus(fiber.StatusOK)
}
