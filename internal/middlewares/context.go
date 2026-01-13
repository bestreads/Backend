package middlewares

import (
	"context"
	"strconv"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"resty.dev/v3"
)

type ctxKey string

const (
	ConfigKey     ctxKey = "config"
	LoggerKey     ctxKey = "logger"
	DBKey         ctxKey = "db"
	HttpClientKey ctxKey = "http_client"
	ValidatorKey  ctxKey = "validator"
)

func ContextMiddleware(cfg *config.Config, logger zerolog.Logger, db *gorm.DB, httpClient *resty.Client, validator *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Add config to ctx
		ctx = context.WithValue(ctx, ConfigKey, cfg)

		// Add logger to ctx
		ctx = context.WithValue(ctx, LoggerKey, logger.With().
			Str("request_id", strconv.FormatUint(c.Context().ID(), 10)).
			Logger(),
		)

		// Add db to ctx
		ctx = context.WithValue(ctx, DBKey, db)

		// Add http client to ctx
		ctx = context.WithValue(ctx, HttpClientKey, httpClient)

		// Add validator to ctx
		ctx = context.WithValue(ctx, ValidatorKey, validator)

		c.SetUserContext(ctx)
		return c.Next()
	}
}

func Config(ctx context.Context) *config.Config {
	return ctx.Value(ConfigKey).(*config.Config)
}

func Logger(ctx context.Context) zerolog.Logger {
	return ctx.Value(LoggerKey).(zerolog.Logger)
}

func DB(ctx context.Context) *gorm.DB {
	return ctx.Value(DBKey).(*gorm.DB)
}

func HttpClient(ctx context.Context) *resty.Client {
	return ctx.Value(HttpClientKey).(*resty.Client)
}

func Validator(ctx context.Context) *validator.Validate {
	return ctx.Value(ValidatorKey).(*validator.Validate)
}

func User(c *fiber.Ctx) *dtos.CustomTokenClaims {
	// Get token from request context
	accessToken := c.Locals(TokenKey).(*jwt.Token)

	// Retrieve claims from token
	claims := accessToken.Claims.(*dtos.CustomTokenClaims)

	return claims
}
