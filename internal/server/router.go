package server

import (
	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func setRoutes(cfg *config.Config, app *fiber.App) {
	basePath := app.Group(cfg.ApiBasePath)
	v1 := basePath.Group("/v1")

	v1.Get("/health", handlers.Health)
	v1.Get("/post", handlers.Post)
}
