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

	v1user := v1.Group("/user")
	v1user.Post("/", handlers.CreateUser)

	auth := v1.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh",
		middlewares.Protected(cfg, log, types.RefreshToken),
		handlers.TokenRefresh,
	)
	auth.Post("/logout", handlers.Logout)

	v1.Get("/health", handlers.Health)

	// Apply the authentication middleware to a new sub-group
	v1Protected := v1.Group("/", middlewares.Protected(cfg, log, types.AccessToken))

	// --- Protected routes ---

	v1userProtected := v1Protected.Group("/user")
	v1userProtected.Get("/", handlers.GetOwnUser)
  v1user.Put("/", handlers.ChangeUserData)
	v1userProtected.Get("/profile/:id", handlers.GetUserProfile)

	// ?limit=n
	v1libProtected := v1userProtected.Group("/:ID/lib")
	v1libWithoutUserIdProtected := v1userProtected.Group("/lib")
	v1libProtected.Get("/", handlers.GetLibrary)
	v1libProtected.Post("/", handlers.AddToLibrary)
	v1libWithoutUserIdProtected.Put("/review", handlers.UpdateReview)
	v1libProtected.Put("/:BID", handlers.UpdateReadingStatus)
	v1libProtected.Delete("/:BID", handlers.DeleteFromLibrary)

	v1booksProtected := v1Protected.Group("/book")
	v1booksProtected.Get("/search", handlers.BookSearch)
	v1booksProtected.Get("/:bid", handlers.GetBook)

	// ?limit=n
	v1Protected.Get("/post", handlers.GetPost)

	v1mediaProtected := v1Protected.Group("/media")
	// ?type=0|1|2
	v1mediaProtected.Put("/", handlers.SaveFile)
	// ?type=0|1|2
	v1mediaProtected.Get("/:KEY", handlers.GetFile)
}
