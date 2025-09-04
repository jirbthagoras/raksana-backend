package services

import "jirbthagoras/raksana-backend/repositories"

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
