package handlers

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type StreakHandler struct {
	StreakService *services.StreakService
}

func NewStreakHandler(
	r *redis.Client,
	s *services.StreakService,
) *StreakHandler {
	return &StreakHandler{
		StreakService: s,
	}
}

func (h *StreakHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/streak")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetStreak)
}

func (h *StreakHandler) handleGetStreak(c *fiber.Ctx) error {
	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	streak, err := h.StreakService.GetCurrentStreak(ctx, int64(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"streak": streak,
		},
	})
}
