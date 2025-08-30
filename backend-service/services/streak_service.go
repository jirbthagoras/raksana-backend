package services

import (
	"context"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"strconv"
	"time"

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

func (s *StreakService) UpdateStreak(ctx context.Context, id int64) error {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	streakKey := fmt.Sprintf("user:%d:streak", id)
	lastCheckinKey := fmt.Sprintf("user:%d:last_checkin", id)
	flagKey := fmt.Sprintf("user:%d:checkin_flag", id)

	exists, err := s.Redis.Exists(ctx, flagKey).Result()
	if err != nil {
		return fmt.Errorf("redis check failed: %w", err)
	}
	if exists > 0 {
		return nil
	}

	lastCheckin, err := s.Redis.Get(ctx, lastCheckinKey).Result()
	if err == redis.Nil {
		lastCheckin = ""
	} else if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	var newStreak int64

	switch {
	case lastCheckin == today:
		return nil

	case lastCheckin == yesterday:
		newStreak, err = s.Redis.Incr(ctx, streakKey).Result()
		if err != nil {
			return fmt.Errorf("redis incr failed: %w", err)
		}

	default:
		err = s.Redis.Set(ctx, streakKey, 1, 0).Err()
		if err != nil {
			return fmt.Errorf("redis reset streak failed: %w", err)
		}
		newStreak = 1
	}

	if err := s.Redis.Set(ctx, lastCheckinKey, today, 0).Err(); err != nil {
		return fmt.Errorf("redis set last_checkin failed: %w", err)
	}

	ttl := helpers.SecondsUntilMidnight()
	if err := s.Redis.Set(ctx, flagKey, 1, time.Duration(ttl)*time.Second).Err(); err != nil {
		return fmt.Errorf("redis set flag failed: %w", err)
	}

	stat, err := s.Repository.GetStatisticByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("db get statistic failed: %w", err)
	}

	if newStreak > int64(stat.LongestStreak) {
		if err := s.Repository.UpdateLongestStreak(ctx, repositories.UpdateLongestStreakParams{
			UserID:        id,
			LongestStreak: int32(newStreak),
		}); err != nil {
			return fmt.Errorf("db update longest streak failed: %w", err)
		}
	}

	return nil
}

func (s *StreakService) GetCurrentStreak(id int) (int, error) {
	streakKey := fmt.Sprintf("user:%v:streak", id)

	result, err := s.Redis.Get(context.Background(), streakKey).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		slog.Error("Failed to get keyval", "err", err)
		return 0, err
	}

	streak, err := strconv.Atoi(result)
	if err != nil {
		slog.Error("Failed to convert value", "err", err)
		return 0, err
	}

	return streak, nil
}
