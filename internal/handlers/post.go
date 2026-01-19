package handlers

import (
	"fmt"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func GetPost(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	log.Info().Msg("GET post")

	pl := struct {
		Uid uint `json:"uid"`
	}{}

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}

	nlimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error parsing int")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad limit",
		})

	}

	if len(c.Body()) == 0 {
		// wir wollen posts von "allen" bekommen
		posts, err := services.GetGlobalPosts(c.UserContext(), int(nlimit))
		if err != nil {
			log.Error().Err(err).Msg("error getting posts")
			return returnInternalError(c)
		}

		return c.Status(fiber.StatusOK).JSON(posts)

	}

	// dieser parser ist eigentlich terror shit, man kann ein leeres obj ("{}") eingeben und kriegt struct {uid: 0, bid: 0} zurück xD
	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("JSON Parser Error!")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})
	}

	posts, err := services.GetPost(c.UserContext(), pl.Uid, int(nlimit))
	if err != nil {
		log.Error().Err(err).Msg("error getting posts")
		return returnInternalError(c)
	}

	return c.Status(fiber.StatusOK).JSON(posts)
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
		ImageURL string `json:"imageurl"`
	}{}

	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("json parsing error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})

	}

	if err = services.CreatePost(c.UserContext(), uint(id), pl.Bid, pl.Content, pl.ImageURL); err != nil {
		log.Error().Err(err).Msg("error creating post")
		return returnInternalError(c)
	}

	return nil
}

func returnInternalError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(
		dtos.GenericRestErrorResponse{
			Description: "Internal Server Error",
		})
}
