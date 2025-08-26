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

func NewStreakHandler(v *validator.Validate, r *redis.Client, rp *repositories.Queries) *StreakHandler {
	return &StreakHandler{
		Validator:  v,
		Redis:      r,
		Repository: rp,
	}
}

func (h *StreakHandler) RegisterRoute(router fiber.Router) {
	g := router.Group("/streak")
	g.Post("/", h.handleUpdateStreak)
}

func (h *StreakHandler) handleUpdateStreak(c *fiber.Ctx) error {
	ctx := context.Background()
	ttl := helpers.SecondsUntilMidnight()
	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("err", err)
		return err
	}

	streakKey := fmt.Sprintf("user:%s:streak", id)
	flagKey := fmt.Sprintf("user:%s:checkin_flag", id)

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
			"streak": newStreak,
		},
	})
}
