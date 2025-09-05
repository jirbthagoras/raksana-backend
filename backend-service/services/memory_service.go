package services

import (
	"context"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
)

type MemoryService struct {
	Repository *repositories.Queries
}

func NewMemoryService(
	r *repositories.Queries,
) *MemoryService {
	return &MemoryService{
		Repository: r,
	}
}

func (s *MemoryService) CreateMemory(description string, fileKey string, userId int) (int, error) {
	id, err := s.Repository.CreateMemory(context.Background(), repositories.CreateMemoryParams{
		UserID:      int64(userId),
		Description: description,
		FileKey:     fileKey,
	})
	if err != nil {
		slog.Error("Failed to create memory", "err", err)
		return 0, err
	}
	return int(id), nil
}
