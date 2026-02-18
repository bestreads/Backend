package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const MAX_UPLOAD_SZ = 5 << 20 // 5mb

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
	payload.Description = c.FormValue("description")

	// Handle Profilbild-Upload
	file, err := c.FormFile("profile_picture")
	if err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) && err.Error() != "there is no uploaded file associated with the given key" {
		return retErr(log, c, fiber.StatusBadRequest, err, "Failed to parse uploaded file")
	}
	if file != nil {
		fileData, err := getMediaBytes(file)
		if err != nil {
			return retErr(log, c, fiber.StatusBadRequest, err, "error processing file")
		}
		payload.ProfilePicture = fileData
	}

	// // Prüfe ob mindestens ein Feld gesetzt ist
	// if payload.IsEmpty() {
	// 	return retErr(log, c, fiber.StatusBadRequest, fmt.Errorf("payload be empty"), "At least one field must be provided")
	// }

	// Validierung der Felder
	// if validationErr := validate.Struct(&payload); validationErr != nil {
	// 	var valErrs validator.ValidationErrors
	// 	if errors.As(validationErr, &valErrs) {
	// 		for _, e := range valErrs {
	// 			if e.Field() == "Email" && e.Tag() == "email" {
	// 				// return c.Status(fiber.StatusBadRequest).
	// 				// 	JSON(dtos.GenericRestErrorResponse{
	// 				// 		Description: "Invalid email format",
	// 				// 	})
	// 			}

	// 			if e.Field() == "Password" && e.Tag() == "min" {
	// 				// return c.Status(fiber.StatusBadRequest).
	// 				// 	JSON(dtos.GenericRestErrorResponse{
	// 				// 		Description: "Password must be at least 12 characters long",
	// 				// 	})
	// 			}
	// 		}
	// 	}

	// 	// Fallback
	// 	return retErr(log, c, fiber.StatusBadRequest, validationErr, "Validation Failed")
	// }
	if c.FormValue("email") != "" {
		panic("handle email here")
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

func getMediaBytes(f *multipart.FileHeader) ([]byte, error) {
	if f.Size > MAX_UPLOAD_SZ {
		return []byte{}, fmt.Errorf("file too large")
	}

	open, err := f.Open()
	bytes, err := io.ReadAll(open)
	if err != nil {
		return []byte{}, err
	} else {
		return bytes, nil
	}
}

func retErr(log zerolog.Logger, c *fiber.Ctx, t int, err error, msg string) error {
	log.Error().Err(err).Msg(msg)
	return c.Status(t).JSON(dtos.GenericRestErrorResponse{
		Description: msg,
	})
}
