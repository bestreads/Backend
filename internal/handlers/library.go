package handlers

import (
	"fmt"
	"strconv"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func GetLibrary(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	var userId uint
	id := c.Query("userId")
	if id != "" {
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Err(err).Msg("error parsing userId")
			return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
				Description: "Invalid userId",
			})
		}
		userId = uint(parsedId)
	} else {
		ownId, err := user.GetId()
		if err != nil {
			msg := "Failed to get user id"
			log.Error().Err(err).Msg(msg)
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
		}
		userId = ownId
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "-1"
	}

	log.Info().Msg(fmt.Sprintf("GET library for user %d with limit %s", userId, limit))

	nlimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error converting limit to int")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "Invalid limit",
		})
	}

	result, err := services.QueryLibrary(c.UserContext(), uint(userId), nlimit)
	if err != nil {
		log.Error().Err(err).Msg("error getting user library")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error getting user library",
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func AddToLibrary(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Get user id from token
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	payload := struct {
		Bid   uint               `json:"bid"`
		State database.ReadState `json:"state"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		log.Error().Err(err).Msg("json parser error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})
	}

	log.Info().Msg(fmt.Sprintf("Adding book with id: %d for user with id: %d", payload.Bid, userId))

	if err := services.AddToLibrary(c.UserContext(), userId, payload.Bid, payload.State); err != nil {
		log.Error().Err(err).Msg("failed to add book to user library")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "Failed to update user library",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func UpdateReadingStatus(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Get user id from token
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	bid, err := strconv.ParseUint(c.Params("BID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing bookid error")
		return err
	}

	log.Info().Msg(fmt.Sprintf("updating state for user %d on book %d", userId, bid))

	payload := struct {
		State database.ReadState `json:"state"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		log.Error().Err(err).Msg("json parser error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})
	}

	if err := services.UpdateReadState(c.UserContext(), userId, uint(bid), payload.State); err != nil {
		log.Error().Err(err).Msg("error updating state")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "error updating state",
		})

	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteFromLibrary(c *fiber.Ctx) error {
	user := middlewares.User(c)
	ctx := c.UserContext()
	log := middlewares.Logger(ctx)

	// Get user id from token
	userId, getUserIdErr := user.GetId()
	if getUserIdErr != nil {
		msg := "Failed to get user id"
		log.Error().Err(getUserIdErr).Msg(msg)
		return c.Status(fiber.StatusInternalServerError).
			JSON(dtos.GenericRestErrorResponse{
				Description: msg,
			})
	}

	bid, err := strconv.ParseUint(c.Params("BID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing bookid error")
		return err
	}

	log.Info().Msg(fmt.Sprintf("deleting book for user %d on book %d", userId, bid))

	if err := services.DeleteFromLibrary(c.UserContext(), userId, uint(bid)); err != nil {
		log.Error().Err(err).Msg("error deleting book from lib")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "error deleting book from library",
		})

	}

	return c.SendStatus(fiber.StatusOK)
}
