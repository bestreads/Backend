package handlers

import (
	"errors"
	"strconv"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

func StartFollow(c *fiber.Ctx) error {
	return followInner(c, false)
}

func StopFollow(c *fiber.Ctx) error {
	return followInner(c, true)
}

func GetFollowers(c *fiber.Ctx) error {
	return getFollowInner(c, false)
}

func GetFollowing(c *fiber.Ctx) error {
	return getFollowInner(c, true)
}

// --- implementierung ---

func getFollowInner(c *fiber.Ctx, following bool) error {
	log := middlewares.Logger(c.UserContext())

	userId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		log.Err(err).Str("id", c.Params("id")).Msg("invalid user ID format")
		return c.Status(fiber.StatusBadRequest).
			JSON(dtos.GenericRestErrorResponse{
				Description: "User ID must be a valid positive number",
			})
	}

	ids, err := services.GetFollow(c.UserContext(), uint(userId), following)
	if err != nil {
		log.Error().Err(err).Msg("failed to query following/follower information")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error follow information",
		})

	}

	return c.Status(fiber.StatusOK).JSON(ids)
}

func followInner(c *fiber.Ctx, unfollow bool) error {
	log := middlewares.Logger(c.UserContext())
	this_id, err := middlewares.User(c).GetId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error getting user id",
		})
	}

	other_id, err := getIdQuery(c)
	if err != nil {
		log.Err(err).Msg("error parsing target id")
		return c.Status(fiber.StatusBadRequest).JSON(dtos.GenericRestErrorResponse{
			Description: "invalid id in requests parameter",
		})
	}

	if err := services.SetFollow(c.UserContext(), this_id, other_id, unfollow); err != nil {
		log.Err(err).Msg("failed to update followers")
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.GenericRestErrorResponse{
			Description: "error updating followers",
		})

	}

	return c.Status(fiber.StatusOK).JSON(dtos.GenericRestResponse{
		Message: "ok",
	})

}

func getIdQuery(c *fiber.Ctx) (uint, error) {
	id_str := c.Query("id")
	if id_str == "" {
		return 0, errors.New("id is empty")
	}

	other_id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(other_id), nil

}
