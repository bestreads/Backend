package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

func TokenRefresh(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Create access JWT
	accessJwt, accessJwtGenerationErr := services.GenerateToken(ctx, user.Subject, types.AccessToken)
	if accessJwtGenerationErr != nil {
		msg := "Failed to create access JWT"
		log.Error().Err(accessJwtGenerationErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	// Create refresh JWT
	refreshJwt, refreshJwtGenerationErr := services.GenerateToken(ctx, user.Subject, types.RefreshToken)
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
			Message: "Refresh successful",
		})
}
