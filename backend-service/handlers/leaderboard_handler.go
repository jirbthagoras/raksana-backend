package handlers

import (
	"jirbthagoras/raksana-backend/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type LeaderboardHandler struct {
	Redis      *redis.Client
	Repository *repositories.Queries
}

func NewLeaderboardHandler(
	rd *redis.Client,
	r *repositories.Queries,
) *LeaderboardHandler {
	return &LeaderboardHandler{
		Redis:      rd,
		Repository: r,
	}
}

func (h *LeaderboardHandler) RegisterRoutes(router fiber.Router) {
	_ = router.Group("/leaderboard")
}
