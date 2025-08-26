package handlers

import (
	"context"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type StreakHandler struct {
	Validator  *validator.Validate
	Redis      *redis.Client
	Repository *repositories.Queries
}

func NewStreakHandler(
	v *validator.Validate,
	r *redis.Client,
	rp *repositories.Queries,
) *StreakHandler {
	return &StreakHandler{
		Validator:  v,
		Redis:      r,
		Repository: rp,
	}
}

func (h *StreakHandler) RegisterRoute(router fiber.Router) {
	g := router.Group("/streak")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleUpdateStreak)
	g.Get("/", h.handleGetStreak)

}

func (h *StreakHandler) handleUpdateStreak(c *fiber.Ctx) error {
	ctx := context.Background()
	ttl := helpers.SecondsUntilMidnight()
	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("Faield to get subject from token", "err", err)
		return err
	}

	streakKey := fmt.Sprintf("user:%v:streak", id)
	flagKey := fmt.Sprintf("user:%v:checkin_flag", id)

	exists, err := h.Redis.Exists(ctx, flagKey).Result()
	if err != nil {
		slog.Error("Failed to check keyval existence")
		return err
	}

	if exists > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "You've already checked in today")
	}

	newStreak, err := h.Redis.Incr(ctx, streakKey).Result()
	if err != nil {
		slog.Error("Failed to incr keyval", "err", err)
		return err
	}

	_, err = h.Redis.Set(ctx, flagKey, 1, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		slog.Error("Failed to set new keyval", "err", err)
		return err
	}

	_, err = h.Redis.Expire(ctx, streakKey, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		slog.Error("Failed to set expire to a keyval", "err", err)
		return err
	}

	stat, err := h.Repository.GetStatisticByUserID(ctx, int64(id))
	if err != nil {
		slog.Error("Failed to get statistics data")
		return err
	}

	if newStreak > int64(stat.LongestStreak) {
		h.Repository.UpdateLongestStreak(ctx, repositories.UpdateLongestStreakParams{
			UserID:        int64(id),
			LongestStreak: int32(newStreak),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}

func (h *StreakHandler) handleGetStreak(c *fiber.Ctx) error {
	id, err := helpers.GetSubjectFromToken(c)
	streakKey := fmt.Sprintf("user:%v:streak", id)
	if err != nil {
		slog.Error("Failed to get subject from token", "err", err)
		return err
	}

	exists, err := h.Redis.Get(context.Background(), streakKey).Result()
	if err != nil {
		slog.Error("Failed to get keyval")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"streak": exists,
		},
	})
}
