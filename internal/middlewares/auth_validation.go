package middlewares

import (
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/types"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	TokenKey = "user_token"
)

// AccessProtected is for regular api routes (access token)
func AccessProtected(cfg *config.Config, log zerolog.Logger) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(cfg.AccessTokenSecretKey)},
		TokenLookup: fmt.Sprintf("cookie:%s", types.AccessToken),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			msg := "Invalid or expired access token"
			log.Debug().Err(err).Msg(msg)
			return c.Status(fiber.StatusUnauthorized).
				JSON(dtos.GenericRestErrorResponse{
					Description: msg,
				})
		},
	})
}

// RefreshProtected is just for /refresh endpoint (refresh token)
func RefreshProtected(cfg *config.Config, log zerolog.Logger) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(cfg.RefreshTokenSecretKey)},
		TokenLookup: fmt.Sprintf("cookie:%s", types.RefreshToken),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			msg := "Invalid or expired refresh token"
			log.Debug().Err(err).Msg(msg)
			return c.Status(fiber.StatusUnauthorized).
				JSON(dtos.GenericRestErrorResponse{
					Description: msg,
				})
		},
	})
}
