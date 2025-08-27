package services

import (
	"context"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type StreakService struct {
	Redis      *redis.Client
	Repository *repositories.Queries
}

func NewStreakService(r *redis.Client, rp *repositories.Queries) *StreakService {
	return &StreakService{
		Redis:      r,
		Repository: rp,
	}
}

func (s *StreakService) UpdateStreak(id int) error {
	ctx := context.Background()
	ttl := helpers.SecondsUntilMidnight()

	streakKey := fmt.Sprintf("user:%v:streak", id)
	flagKey := fmt.Sprintf("user:%v:checkin_flag", id)

	exists, err := s.Redis.Exists(ctx, flagKey).Result()
	if err != nil {
		slog.Error("Failed to check keyval existence")
		return err
	}

	if exists > 0 {
		slog.Error("User already checked in today")
		return fiber.NewError(fiber.StatusBadRequest, "You've already checked in today")
	}

	newStreak, err := s.Redis.Incr(ctx, streakKey).Result()
	if err != nil {
		slog.Error("Failed to incr keyval", "err", err)
		return err
	}

	_, err = s.Redis.Set(ctx, flagKey, 1, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		slog.Error("Failed to set new keyval", "err", err)
		return err
	}

	_, err = s.Redis.Expire(ctx, streakKey, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		slog.Error("Failed to set expire to a keyval", "err", err)
		return err
	}

	stat, err := s.Repository.GetStatisticByUserID(ctx, int64(id))
	if err != nil {
		slog.Error("Failed to get statistics data")
		return err
	}

	if newStreak > int64(stat.LongestStreak) {
		s.Repository.UpdateLongestStreak(ctx, repositories.UpdateLongestStreakParams{
			UserID:        int64(id),
			LongestStreak: int32(newStreak),
		})
	}

	return nil
}

func (s *StreakService) GetCurrentStreak(id int) (int, error) {
	streakKey := fmt.Sprintf("user:%v:streak", id)

	result, err := s.Redis.Get(context.Background(), streakKey).Result()
	if err != nil {
		slog.Error("Failed to get keyval")
		return 0, err
	}

	streak, err := strconv.Atoi(result)
	if err != nil {
		slog.Error("Failed to convert value")
		return 0, err
	}

	return streak, nil
}
