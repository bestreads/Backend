package server

import (
	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/handlers"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func setRoutes(cfg *config.Config, log zerolog.Logger, app *fiber.App) {
	// Create the main API base path group
	basePath := app.Group(cfg.ApiBasePath)

	// Initialize api version groups
	v1 := basePath.Group("/v1")

	// --- Public routes ---

	auth := v1.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh",
		middlewares.Protected(cfg, log, types.RefreshToken),
		handlers.TokenRefresh,
	)
	auth.Post("/logout", handlers.Logout)

	v1.Get("/health", handlers.Health)
	v1.Get("/books/search", handlers.BookSearch)
	v1.Get("/books/:bid", handlers.GetBook)

	// ?type=0|1|2
	v1.Put("/media", handlers.SaveFile)
	// ?type=0|1|2
	v1.Get("/media/:KEY", handlers.GetFile)

	// Apply the authentication middleware to a new sub-group
	v1Protected := v1.Group("/", middlewares.Protected(cfg, log, types.AccessToken))

	// ?limit=number
	v1Protected.Get("/post", handlers.GetPost)

	// --- Protected routes ---
	v1user := v1Protected.Group("/user")
	v1user.Get("/", handlers.GetOwnUser)
	v1user.Post("/", handlers.CreateUser)
	v1user.Put("/", handlers.ChangeUserData)
	v1user.Get("/profile/:id", handlers.GetUserProfile)
	v1user.Post("/post", handlers.CreatePost)

	v1userWithId := v1user.Group("/:ID")
	// ?limit=n
	v1userWithId.Get("/lib", handlers.GetLibrary)
	v1userWithId.Post("/lib", handlers.AddToLibrary)
	v1userWithId.Put("/lib/:BID", handlers.UpdateReadingStatus)
	v1userWithId.Delete("/lib/:BID", handlers.DeleteFromLibrary)
}
