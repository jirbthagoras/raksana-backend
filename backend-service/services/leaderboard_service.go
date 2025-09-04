package services

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type LeaderboardService struct {
	Redis   *redis.Client
	KeyName string
}

func NewLeaderboardService(
	redis *redis.Client,
) *LeaderboardService {
	return &LeaderboardService{
		Redis:   redis,
		KeyName: "raksana:leaderboard",
	}
}

func (s *LeaderboardService) UpdatePoint(userId string, points float64) error {
	ctx := context.Background()
	_, err := s.Redis.ZAdd(ctx, s.KeyName, redis.Z{
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
	_, err := s.Redis.ZIncrBy(ctx, s.KeyName, points, userId).Result()
	if err != nil {
		slog.Error("Faield to incraese point", "err", err)
		return err
	}
	return nil
}

func (s *LeaderboardService) GetUserScore(userId string) (float64, error) {
	ctx := context.Background()
	score, err := s.Redis.ZScore(ctx, s.KeyName, userId).Result()
	if err != nil {
		slog.Error("Failed to get user score", "err", err)
		return 0, err
	}
	return score, nil
}

func (s *LeaderboardService) GetUserRank(userId string) (int64, error) {
	ctx := context.Background()
	rank, err := s.Redis.ZRevRank(ctx, s.KeyName, userId).Result()
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
	ID       string
	Name     string
	ImageUrl string
	Score    int
	Rank     int
	IsUser   bool
}

func (s *LeaderboardService) GetUserInfo(userId string) (UserInfoRedis, error) {
	ctx := context.Background()
	var user UserInfoRedis
	name, err := s.Redis.HGet(ctx, "user:"+userId, "name").Result()
	if err != nil {
		return user, err
	}
	imageUrl, err := s.Redis.HGet(ctx, "user:"+userId, "name").Result()
	if err != nil {
		return user, err
	}
	return UserInfoRedis{
		ID:       userId,
		Name:     name,
		ImageUrl: imageUrl,
	}, nil
}

func (s *LeaderboardService) GetTopLeaderboard(currentUserId string) ([]UserInfoRedis, error) {
	ctx := context.Background()
	results, err := s.Redis.ZRevRangeWithScores(ctx, s.KeyName, 0, 100000).Result()
	if err != nil {
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
		row.Score = int(z.Score)
		row.IsUser = false
		if currentUserId == userId {
			row.IsUser = true
		}

		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}
