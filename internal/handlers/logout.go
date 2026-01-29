package handlers

import (
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/services"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

// Logout handles the user logout process.
// @Summary Logout of a user account
// @Description Logs the user out by invalidating cookies from the client’s browser.
// @Tags User Management
// @Produce json
// @Success 200 {object} dtos.GenericRestResponse "Logout successful"
// @Router /v1/auth/logout [post]
func Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Generate expired cookies with empty values to override the existing ones, and thus remove them
	accessTokenCookie := services.CreateCookie(ctx, types.AccessToken, "", true, true)
	refreshTokenCookie := services.CreateCookie(ctx, types.RefreshToken, "", true, true)

	// Set the cookies to the response to override the existing ones
	c.Cookie(accessTokenCookie)
	c.Cookie(refreshTokenCookie)

	return c.Status(fiber.StatusOK).
		JSON(dtos.GenericRestResponse{
			Message: "Logout successful",
		})
}
