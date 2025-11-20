package server

import (
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func Start(cfg *config.Config, logger zerolog.Logger) {
	app := fiber.New()

	db := database.SetupDatabase(cfg)

	// Attach logger + db to ctx for every request
	app.Use(middlewares.ContextMiddleware(logger, db))

	setRoutes(cfg, app)

	logger.Info().Msg(fmt.Sprintf("API started on :%s", cfg.ApiPort))
	if err := app.Listen(fmt.Sprintf(":%s", cfg.ApiPort)); err != nil {
		logger.Fatal().Err(err)
	}
}
