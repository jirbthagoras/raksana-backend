package handlers

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LeaderboardHandler struct {
	*services.LeaderboardService
}

func NewLeaderboardHandler(
	ls *services.LeaderboardService,
) *LeaderboardHandler {
	return &LeaderboardHandler{
		LeaderboardService: ls,
	}
}

func (h *LeaderboardHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/leaderboard")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleLeaderboard)
}

func (h *LeaderboardHandler) handleLeaderboard(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	leaderboard, err := h.LeaderboardService.GetTopLeaderboard(strconv.Itoa(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": leaderboard,
	})
}
