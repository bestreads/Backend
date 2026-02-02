package server

import (
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/handlers"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
)

func setRoutes(cfg *config.Config, log zerolog.Logger, app *fiber.App) {
	// Create the main API base path group
	basePath := app.Group(cfg.ApiBasePath)

	// Initialize api version groups
	v1 := basePath.Group("/v1")

	// --- Public routes ---

	swaggerGroupName := "/swagger"
	swaggerFileName := "docs.yaml"
	swaggerGroup := v1.Group(swaggerGroupName)

	// Swagger document provider
	swaggerGroup.Get(fmt.Sprintf("/%s", swaggerFileName), func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.yaml")
	})

	// Swagger UI provider
	swaggerGroup.Get("/docs/*", swagger.New(
		swagger.Config{ // Custom url/path for file provider:
			URL: fmt.Sprintf("%s://%s:%s%s/v1%s/%s",
				cfg.ApiProtocol, cfg.ApiDomain, cfg.ApiProdPort, cfg.ApiBasePath, swaggerGroupName, swaggerFileName),
		}))

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
	v1userProtected.Put("/", handlers.ChangeUserData)
	v1userProtected.Get("/:id", handlers.GetUserProfile)

	// library
	// ?limit=n
	v1LibProtected := v1Protected.Group("/lib")
	v1LibProtected.Get("/", handlers.GetLibrary)
	v1LibProtected.Post("/", handlers.AddToLibrary)
	v1LibProtected.Put("/review", handlers.UpdateReview)
	v1LibProtected.Put("/:BID", handlers.UpdateReadingStatus)
	v1LibProtected.Delete("/:BID", handlers.DeleteFromLibrary)

	v1booksProtected := v1Protected.Group("/book")
	v1booksProtected.Get("/search", handlers.BookSearch)
	v1booksProtected.Get("/:bid", handlers.GetBook)

	// ?limit=n
	v1Protected.Get("/post", handlers.GetPost)
	v1Protected.Post("/post", handlers.CreatePost)
	v1Protected.Delete("/post", handlers.DeletePost)

	v1mediaProtected := v1Protected.Group("/media")
	// ?type=0|1|2
	v1mediaProtected.Put("/", handlers.SaveFile)
	// ?type=0|1|2
	v1mediaProtected.Get("/:KEY", handlers.GetFile)
}
