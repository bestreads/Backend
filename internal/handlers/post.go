package handlers

import (
	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type postResponse struct {
	Uid     uint
	Bid     uint
	Content string
	Image   string
}

func GetPost(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("GET post")

	pl := struct {
		Uid uint `json:"uid"`
		Bid uint `json:"bid"`
	}{}

	// dieser parser ist eigentlich terror shit, man kann ein leeres obj ("{}") eingeben und kriegt struct {uid: 0, bid: 0} zurück xD
	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("JSON Parser Error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "malformed request body (invalid json?)",
		})
	}

	posts, err := gorm.G[database.Post](database.GlobalDB).Where("user_id = ? AND book_id = ?", pl.Uid, pl.Bid).Find(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(convert(posts))
}

func convert(p []database.Post) []postResponse {
	res := make([]postResponse, len(p))
	for i, post := range p {
		res[i] = postResponse{
			Uid:     post.UserID,
			Bid:     post.BookID,
			Content: post.Content,
			Image:   post.Image,
		}
	}

	return res
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
