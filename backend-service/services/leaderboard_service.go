package services

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type LeaderboardService struct {
	Redis *redis.Client
}

func NewLeaderboardService(
	redis *redis.Client,
) *LeaderboardService {
	return &LeaderboardService{
		Redis: redis,
	}
}

func (s *LeaderboardService) UpdatePoint(userId string, points float64) error {
	ctx := context.Background()
	_, err := s.Redis.ZAdd(ctx, "raksana:leaderboard", redis.Z{
		Score:  points,
		Member: userId,
	}).Result()
	if err != nil {
		slog.Error("Failed to add data to leaderboard", "err", err)
		return err
	}
	return nil
}

func (s *LeaderboardService) IncrPoint(userId string, points float64) error {
	ctx := context.Background()
	_, err := s.Redis.ZIncrBy(ctx, "raksana:leaderboard", points, userId).Result()
	if err != nil {
		slog.Error("Faield to incraese point", "err", err)
		return err
	}
	return nil
}

func (s *LeaderboardService) GetUserScore(userId string) (float64, error) {
	ctx := context.Background()
	score, err := s.Redis.ZScore(ctx, "raksana:leaderboard", userId).Result()
	if err != nil {
		slog.Error("Failed to get user score", "err", err)
		return 0, err
	}
	return score, nil
}

func (s *LeaderboardService) GetUserRank(userId string) (int64, error) {
	ctx := context.Background()
	rank, err := s.Redis.ZRevRank(ctx, "raksana:leaderboard", userId).Result()
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (s *LeaderboardService) SetUserInfo(userId, name string, imageUrl string) error {
	ctx := context.Background()
	_, err := s.Redis.HSet(ctx, "user:leaderboard:"+userId, "name", name, "profile", imageUrl).Result()
	if err != nil {
		slog.Error("Failed to set user info", "err", err)
		return err
	}
	return err
}

func (s *LeaderboardService) UpdateProfile(userId string, imageUrl string) error {
	ctx := context.Background()
	_, err := s.Redis.HSet(ctx, "user:leaderboard:"+userId, "profile", imageUrl).Result()
	if err != nil {
		slog.Error("Failed to set user info", "err", err)
		return err
	}
	return nil
}

type UserInfoRedis struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Points   int    `json:"points"`
	Rank     int    `json:"rank"`
	IsUser   bool   `json:"is_user"`
}

func (s *LeaderboardService) GetUserInfo(userId string) (UserInfoRedis, error) {
	ctx := context.Background()
	var user UserInfoRedis
	name, err := s.Redis.HGet(ctx, "user:leaderboard:"+userId, "name").Result()
	if err != nil {
		slog.Error("Failed to get user info: name", "err", err)
		return user, err
	}
	imageUrl, err := s.Redis.HGet(ctx, "user:leaderboard:"+userId, "profile").Result()
	if err != nil {
		slog.Error("Failed to get user info: profile", "err", err)
		return user, err
	}
	return UserInfoRedis{
		ID:       userId,
		Name:     name,
		ImageUrl: imageUrl,
	}, nil
}

func (s *LeaderboardService) GetTopLeaderboard(currentUserId string) (
	[]UserInfoRedis,
	error,
) {
	ctx := context.Background()
	results, err := s.Redis.ZRevRangeWithScores(ctx, "raksana:leaderboard", 0, 100000).Result()
	if err != nil {
		slog.Error("Failed to get", "err", err)
		return nil, err
	}

	var leaderboard []UserInfoRedis

	for i, z := range results {
		userId := z.Member.(string)
		row, err := s.GetUserInfo(userId)
		if err != nil {
			return nil, err
		}

		row.Rank = i + 1
		row.Points = int(z.Score)
		row.IsUser = false
		if currentUserId == userId {
			row.IsUser = true
		}

		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}
