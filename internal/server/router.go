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
	v1.Get("/books/search", handlers.BookSearch)
	v1.Get("/books/:bid", handlers.GetBook)

	// ?limit=n
	v1.Get("/post", handlers.GetPost)

	v1.Post("/login", handlers.Login)
	v1.Post("/user", handlers.CreateUser)
	v1.Get("/user/profile/:id", handlers.GetUserProfile)

	// ?type=0|1|2
	v1.Put("/media", handlers.SaveFile)
	// ?type=0|1|2
	v1.Get("/media/:KEY", handlers.GetFile)

	v1user := v1.Group("/user/:ID") // vllt hier so eine auth middleware
	v1user.Post("/post", handlers.CreatePost)

	// ?limit=n
	v1user.Get("/lib", handlers.GetLibrary)
	v1user.Post("/lib", handlers.AddToLibrary)
	v1user.Put("/lib/:BID", handlers.UpdateReadingStatus)
	v1user.Delete("/lib/:BID", handlers.DeleteFromLibrary)
}
