package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/services"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	ctx := c.Context()

	// Generate expired cookies with empty values to override the existing ones, and thus remove them
	accessTokenCookie := services.CreateCookie(ctx, types.AccessToken, "", true, false)
	refreshTokenCookie := services.CreateCookie(ctx, types.RefreshToken, "", true, false)

	// Set the cookies to the response to override the existing ones
	c.Cookie(accessTokenCookie)
	c.Cookie(refreshTokenCookie)

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "Logout successful",
		})
}
