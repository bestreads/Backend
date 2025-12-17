package handlers

import (
	"fmt"
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type postResponse struct {
	Pfp      string
	Username string
	Book     database.Book
	Content  string
	Image    string
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

	// wichtige informationen:
	// 1. den namen der preload dinger *RICHTIG* zu schreiben (auch groß- und kleinschreibung)
	// 2. gorm wird (sporadisch) andere sachen preloaden, keine ahnung wieso?
	// 3. eigentlich wollte ich eine api machen, mit der man nur die metadaten von einem user lädt.
	//    würde ein bisschen hübscher in json aussehen, das habe ich aber mal nicht gemacht
	posts, err := gorm.G[database.Post](database.GlobalDB).Preload("User", nil).Preload("Book", nil).Where("user_id = ? AND book_id = ?", pl.Uid, pl.Bid).Find(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Database query error")
		return returnInternalError(c)
	}

	returnData, err := convert(posts)
	if err != nil {
		log.Error().Err(err).Msg("Image retrieval error")
		return returnInternalError(c)
	}

	return c.Status(fiber.StatusOK).JSON(returnData)
}

// gerade wird das bild noch automatisch aus dem storage wieder zurückgeholt,
// in zukunft vllt durch eine "data" api oder so
//
// man müsste das hier eig auch wegmachen und durch verschachtelte structs
// (sachen mit FKs in der db) machen
func convert(p []database.Post) ([]postResponse, error) {
	res := make([]postResponse, len(p))
	for i, post := range p {
		imageData, err := database.FileRetrieve(post.ImageHash, database.PostImage)
		if err != nil {
			return make([]postResponse, 0), err
		}

		res[i] = postResponse{
			Pfp:      post.User.Pfp,
			Username: post.User.Username,
			Book:     post.Book,
			Content:  post.Content,
			Image:    imageData,
		}
	}

	return res, nil
}

func CreatePost(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())

	id, err := strconv.ParseUint(c.Params("ID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing id error")
		return returnInternalError(c)
	}

	log.Info().Msg(fmt.Sprintf("POST post for user %d", id))

	pl := struct {
		Bid      uint   `json:"bid"`
		Content  string `json:"content"`
		B64Image string `json:"b64image"`
	}{}

	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("json parsing error")
		return returnInternalError(c)
	}

	// leeres bild wird einfach "0", das ist okay glaube ich
	hash, err := database.FileStore(pl.B64Image, database.PostImage)
	if err != nil {
		log.Error().Err(err).Msg("file storage error")
		return returnInternalError(c)
	}

	post := database.Post{
		UserID:    uint(id),
		BookID:    pl.Bid,
		Content:   pl.Content,
		ImageHash: strconv.Itoa(hash),
	}

	if err := gorm.G[database.Post](database.GlobalDB).Create(c.Context(), &post); err != nil {
		log.Error().Err(err).Msg("Database post insertion error")
		return returnInternalError(c)
	}

	return nil
}

func returnInternalError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "Internal Server Error",
	})
}
