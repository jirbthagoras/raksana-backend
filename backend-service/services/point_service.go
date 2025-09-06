package services

import (
	"context"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"strconv"
)

type PointService struct {
	Repository *repositories.Queries
	*LeaderboardService
}

func NewPointService(
	rp *repositories.Queries,
	ls *LeaderboardService,
) *PointService {
	return &PointService{
		Repository:         rp,
		LeaderboardService: ls,
	}
}

func (s *PointService) UpdateUserPoint(userId int64, pointGain int64) (repositories.Profile, error) {
	ctx := context.Background()
	profile, err := s.Repository.IncreaseUserPoints(ctx, repositories.IncreaseUserPointsParams{
		UserID: userId,
		Points: pointGain,
	})
	if err != nil {
		slog.Error("Failed to increase user points", "err", err)
		return profile, err
	}

	err = s.LeaderboardService.IncrPoint(strconv.Itoa(int(userId)), float64(pointGain))
	if err != nil {
		return profile, err
	}

	return profile, nil
}
