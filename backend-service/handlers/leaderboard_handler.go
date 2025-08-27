package handlers

import (
	"jirbthagoras/raksana-backend/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type LeaderboardHandler struct {
	Validator  *validator.Validate
	Redis      *redis.Client
	Repository *repositories.Queries
}

func NewLeaderboardHandler(
	v *validator.Validate,
	rd *redis.Client,
	r *repositories.Queries,
) *LeaderboardHandler {
	return &LeaderboardHandler{
		Validator:  v,
		Redis:      rd,
		Repository: r,
	}
}

func (h *LeaderboardHandler) RegisterRoutes(router fiber.Router) {
	_ = router.Group("/leaderboard")
}
