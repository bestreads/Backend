package middlewares

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ctxKey string

const (
	LoggerKey ctxKey = "logger"
	DBKey     ctxKey = "db"
)

func ContextMiddleware(logger zerolog.Logger, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Add logger to ctx
		ctx = context.WithValue(ctx, LoggerKey, logger.With().
			Str("request_id", strconv.FormatUint(c.Context().ID(), 10)).
			Logger(),
		)

		// Add db to ctx
		ctx = context.WithValue(ctx, DBKey, db)

		c.SetUserContext(ctx)
		return c.Next()
	}
}

func Logger(ctx context.Context) zerolog.Logger {
	return ctx.Value(LoggerKey).(zerolog.Logger)
}

func DB(ctx context.Context) *gorm.DB {
	return ctx.Value(DBKey).(*gorm.DB)
}
