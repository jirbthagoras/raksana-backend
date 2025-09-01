package services

import (
	"context"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
)

type JournalService struct {
	Repository *repositories.Queries
}

func NewJournalService(
	rp *repositories.Queries,
) *JournalService {
	return &JournalService{
		Repository: rp,
	}
}

func (s *JournalService) AppendLog(req *models.PostLogAppend, userId int) error {
	_, err := s.Repository.CreateLog(context.Background(), repositories.CreateLogParams{
		UserID:    int64(userId),
		Text:      req.Text,
		IsSystem:  req.IsSystem,
		IsPrivate: req.IsPrivate,
	})
	if err != nil {
		slog.Error("Failed to append log", "err", err.Error())
		return err
	}

	return nil
}

func (s *JournalService) GetLogs(id int, isPrivate bool) ([]models.ResponseGetLogs, error) {
	var logs []models.ResponseGetLogs

	result, err := s.Repository.GetLogs(context.Background(), repositories.GetLogsParams{
		UserID:    int64(id),
		IsPrivate: isPrivate,
	})
	if err != nil {
		slog.Error("Failed to get logs", "err", err)
		return logs, err
	}

	logs = []models.ResponseGetLogs{}
	for _, log := range result {
		logs = append(logs, models.ResponseGetLogs{
			Text:      log.Text,
			IsSystem:  log.IsSystem,
			IsPrivate: log.IsPrivate,
			CreatedAt: log.CreatedAt,
		})
	}

	return logs, nil
}
