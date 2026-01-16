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

// Protected secures api routes
func Protected(cfg *config.Config, log zerolog.Logger, tokenType types.TokenType) fiber.Handler {
	// Get correct secret
	var tokenSecretKey []byte
	if tokenType == types.RefreshToken {
		tokenSecretKey = []byte(cfg.RefreshTokenSecretKey)
	} else {
		tokenSecretKey = []byte(cfg.AccessTokenSecretKey)
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    tokenSecretKey,
		},
		TokenLookup: fmt.Sprintf("cookie:%s", tokenType),

		// ContextKey stores the data in c.Locals(TokenKey)
		ContextKey: TokenKey,

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			msg := "Invalid or expired token"
			log.Debug().Err(err).Msg(msg)
			return c.Status(fiber.StatusUnauthorized).
				JSON(dtos.GenericRestErrorResponse{
					Description: msg,
				})
		},
	})
}
