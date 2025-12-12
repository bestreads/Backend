package handlers

import (
	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPost(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("GET demopost")

	pl := struct {
		Uid uint `json:"uid"`
		Bid uint `json:"bid"`
	}{}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/400
	if err := c.BodyParser(&pl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "malformed request body (invalid json?)",
		})
	}

	// maybe gucken, ob bid/uid nicht da sind und dann alle posts zurückgeben?

	post, err := gorm.G[database.Post](database.GlobalDB).Where("uid = ? AND bid = ?", pl.Uid, pl.Bid).First(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

func CreatePost(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("POST demopost")

	pl := struct {
		Uid     uint   `json:"uid"`
		Bid     uint   `json:"bid"`
		Content string `json:"content"`
		// Image   string `json:"image"`
	}{}

	if err := c.BodyParser(&pl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "malformed request body (invalid json?)",
		})
	}

	post := database.Post{
		UserID:  pl.Uid,
		BookID:  pl.Bid,
		Content: pl.Content,
		// Image: pl.Image,
	}

	if err := gorm.G[database.Post](database.GlobalDB).Create(c.Context(), &post); err != nil {
		return err
	}

	return nil
}
