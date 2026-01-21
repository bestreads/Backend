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

	limit := c.Query("limit")
	if limit == "" {
		limit = "-1"
	}

	userId := c.Query("userId")
	if userId == "" {
		userId = "0"
	}

	nlimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error parsing int")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad limit",
		})

	}

	nUserId, err := strconv.ParseInt(userId, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error parsing int")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "bad userId",
		})

	}

	if nUserId == 0 {
		// wir wollen posts von "allen" bekommen
		posts, err := services.GetGlobalPosts(c.UserContext(), int(nlimit))
		if err != nil {
			log.Error().Err(err).Msg("error getting posts")
			return returnInternalError(c)
		}

		return c.Status(fiber.StatusOK).JSON(posts)
	}

	posts, err := services.GetPost(c.UserContext(), uint(nUserId), int(nlimit))
	if err != nil {
		log.Error().Err(err).Msg("error getting posts")
		return returnInternalError(c)
	}

	return c.Status(fiber.StatusOK).JSON(posts)
}

func CreatePost(c *fiber.Ctx) error {
	user := middlewares.User(c)
	log := middlewares.Logger(c.UserContext())

	// Get user id from token
	userId, err := user.GetId()
	if err != nil {
		msg := "Failed to get user id"
		log.Error().Err(err).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	log.Info().Msg(fmt.Sprintf("POST post for user %d", userId))

	payload := struct {
		Bid     uint   `json:"bid"`
		Content string `json:"content"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		log.Error().Err(err).Msg("json parsing error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})
	}

	if payload.Bid == 0 {
		log.Error().Msg("Invalid bookID: bid is missing or 0")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "book id invalid or missing",
		})
	}

	if payload.Content == "" {
		log.Error().Err(err).Msg("No content: " + payload.Content)
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "Content must be present",
		})
	}

	if err = services.CreatePost(c.UserContext(), uint(userId), payload.Bid, payload.Content); err != nil {
		log.Error().Err(err).Msg("error creating post")
		return returnInternalError(c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func returnInternalError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(
		dtos.GenericRestErrorResponse{
			Description: "Internal Server Error",
		})
}
