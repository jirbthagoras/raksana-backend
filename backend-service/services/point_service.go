package services

import (
	"context"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
)

type PointService struct {
	Repository *repositories.Queries
}

func NewPointService(
	rp *repositories.Queries,
) *PointService {
	return &PointService{
		Repository: rp,
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

	return profile, nil
}
