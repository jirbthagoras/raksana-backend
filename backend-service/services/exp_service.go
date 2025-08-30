package services

import (
	"context"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
)

type ExpService struct {
	Repository *repositories.Queries
	*JournalService
}

func NewExpService(
	rp *repositories.Queries,
	s *JournalService,
) *ExpService {
	return &ExpService{
		Repository:     rp,
		JournalService: s,
	}
}

func (s *ExpService) IncreaseExp(userId int, expGain int) error {
	profile, err := s.Repository.IncreaseExp(context.Background(), repositories.IncreaseExpParams{
		ExpGain: int32(expGain),
		UserID:  int32(userId),
	})
	if err != nil {
		slog.Error("Failed to update profile", "err", err)
		return err
	}

	// i think it will repeatedly increase level based on gained exp
	for profile.CurrentExp >= profile.ExpNeeded {
		profile.CurrentExp -= profile.ExpNeeded

		profile.ExpNeeded = int64(helpers.CalculateExpNeeded(int(profile.Level)))

		level, err := s.Repository.UpdateLevelAndExpNeeded(
			context.Background(),
			repositories.UpdateLevelAndExpNeededParams{
				ExpNeeded: profile.ExpNeeded,
				UserID:    int64(userId),
			},
		)
		if err != nil {
			slog.Error("Failed to update profiles exp_needed and level", "err", err)
			return err
		}

		// append log as system log
		logMsg := fmt.Sprintf("Baru saja naik level! Sekarang level %v", level)
		err = s.JournalService.AppendLog(&models.PostLogAppend{
			IsSystem:  true,
			IsPrivate: false,
			Text:      logMsg,
		}, userId)
		if err != nil {
			return err
		}

		profile.Level = level
	}

	return nil
}
