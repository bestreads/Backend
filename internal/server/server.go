package server

import (
	"context"
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"resty.dev/v3"
)

func Start(cfg *config.Config, logger zerolog.Logger) {
	app := fiber.New()

	if cfg.DevMode {
		// CORS Middleware
		app.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:5173, http://localhost:3000",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		}))
	}

	db, dbErr := database.SetupDatabase(cfg, context.TODO())
	if dbErr != nil {
		logger.Fatal().Err(dbErr).Msg("Database connection could not be established")
	}

	logger.Info().Msg("connected to database")

	// Setup http client
	httpClient := resty.New()
	httpClient.SetDisableWarn(true) // Disable warnings, because we're sending requests from API to KC without TLS but within an isolated Docker bridge network
	defer httpClient.Close()

	// Setup validator
	validator := validator.New(validator.WithRequiredStructEnabled())

	// Attach logger + db to ctx for every request
	app.Use(middlewares.ContextMiddleware(cfg, logger, db, httpClient, validator))

	setRoutes(cfg, app)

	logger.Info().Msg(fmt.Sprintf("API started on :%s", cfg.ApiPort))
	if err := app.Listen(fmt.Sprintf(":%s", cfg.ApiPort)); err != nil {
		logger.Fatal().Err(err)
	}
}
