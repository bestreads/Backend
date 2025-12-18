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
	log := middlewares.Logger(c.UserContext())
	id, err := strconv.ParseUint(c.Params("ID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing id error")
		return err
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "1"
	}

	log.Info().Msg(fmt.Sprintf("GET library for user %d with limit %s", id, limit))

	nlimit, err := strconv.ParseUint(limit, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("error converting limit to int")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "Invalid limit",
		})
	}

	result, err := services.QueryLibrary(c.UserContext(), uint(id), nlimit)
	if err != nil {
		log.Error().Err(err).Msg("error getting user library")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error getting user library",
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func AddToLibrary(c *fiber.Ctx) error {
	log := middlewares.Logger(c.UserContext())
	id, err := strconv.ParseUint(c.Params("ID"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("parsing id error")
		return err
	}

	pl := struct {
		Bid   uint               `json:"bid"`
		State database.ReadState `json:"state"`
	}{}

	if err := c.BodyParser(&pl); err != nil {
		log.Error().Err(err).Msg("json parser error")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "JSON invalid",
		})
	}

	if err := services.AddToLibrary(c.UserContext(), uint(id), pl.Bid, pl.State); err != nil {
		log.Error().Err(err).Msg("failed to add book to user library")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "Failed to update user library",
		})
	}

	return nil
}

func UpdateReadingStatus(c *fiber.Ctx) error {
	return nil
}

func DeleteFromLibrary(c *fiber.Ctx) error {
	return nil
}
