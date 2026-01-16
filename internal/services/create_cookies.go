package services

import (
	"context"

	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

// CreateCookie generates a fiber cookie for the given token type.
// If 'rememberMe' is true, the cookie includes a MaxAge, making it persistent for the specified duration.
func CreateCookie(ctx context.Context, tokenType types.TokenType, cookieJwt string, rememberMe bool, expired bool) *fiber.Cookie {
	cfg := middlewares.Config(ctx)

	// Set cookie path according to jwt type
	cookiePath := "/"
	if tokenType == types.RefreshToken {
		cookiePath = cfg.TokenRefreshPath
	}

	// Create cookie base structure
	cookie := &fiber.Cookie{
		Name:     string(tokenType),
		Value:    cookieJwt,
		Path:     cookiePath,
		Domain:   cfg.ApiDomain,
		Secure:   cfg.TokenSecureFlag,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	}

	// Set MaxAge only if 'rememberMe' is true
	if rememberMe {
		switch tokenType {
		case types.AccessToken:
			cookie.MaxAge = int(cfg.AccessTokenDurationMinutes) * 60
		case types.RefreshToken:
			cookie.MaxAge = int(cfg.RefreshTokenDurationDays) * 24 * 60 * 60
		}
	}

	// Set maxAge to -1 if expired is true
	if expired {
		cookie.MaxAge = -1
	}

	return cookie
}
