package server

import (
	"fmt"
	"os"

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

	// Swagger document provider - check Docker path first, then relative (local dev)
	swaggerFilePath := "./docs/swagger.yaml"
	if _, err := os.Stat("/app/docs/swagger.yaml"); err == nil {
		swaggerFilePath = "/app/docs/swagger.yaml"
	}
	swaggerGroup.Get(fmt.Sprintf("/%s", swaggerFileName), func(c *fiber.Ctx) error {
		return c.SendFile(swaggerFilePath)
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

	// user
	v1userProtected := v1Protected.Group("/user")
	v1userProtected.Get("/", handlers.GetOwnUser)
	v1userProtected.Put("/", handlers.ChangeUserData)
	v1userProtected.Get("/:id", handlers.GetUserProfile)
	v1userProtected.Get("/:id/followers", handlers.GetFollowers)
	v1userProtected.Get("/:id/following", handlers.GetFollowing)

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

	// posts
	// ?limit=n
	v1Protected.Get("/post", handlers.GetPost)
	v1Protected.Post("/post", handlers.CreatePost)
	v1Protected.Delete("/post", handlers.DeletePost)
	v1Protected.Put("/post", handlers.CreatePost)

	// media
	v1mediaProtected := v1Protected.Group("/media")
	v1mediaProtected.Put("/", handlers.SaveFile)
	v1mediaProtected.Get("/:KEY", handlers.GetFile)

	// follow
	// ?id=n
	v1Protected.Post("/follow", handlers.StartFollow)
	// ?id=n
	v1Protected.Delete("/follow", handlers.StopFollow)
}
