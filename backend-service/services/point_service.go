package services

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
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

func (s *PointService) UpdateUserPoint(userId int64, pointGain int64, name string, category string, userLevel int) (repositories.Profile, error) {
	ctx := context.Background()

	multiplier := helpers.GetMultiplier(userLevel)
	realPoint := int(float64(pointGain) * multiplier)
	profile, err := s.Repository.IncreaseUserPoints(ctx, repositories.IncreaseUserPointsParams{
		UserID: userId,
		Points: int64(realPoint),
	})
	if err != nil {
		slog.Error("Failed to increase user points", "err", err)
		return profile, err
	}

	err = s.LeaderboardService.IncrPoint(strconv.Itoa(int(userId)), float64(realPoint))
	if err != nil {
		return profile, err
	}

	err = s.Repository.AppendHistry(ctx, repositories.AppendHistryParams{
		UserID:   userId,
		Amount:   int32(realPoint),
		Type:     "input",
		Category: category,
		Name:     name,
	})

	return profile, nil
}
