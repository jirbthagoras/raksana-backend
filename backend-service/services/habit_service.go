package services

import (
	"context"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"sync"
)

type HabitService struct {
	Repository *repositories.Queries
	Lock       sync.Mutex
	*StreakService
}

func NewHabitService(
	rp *repositories.Queries,
	ss *StreakService,
) *HabitService {
	return &HabitService{
		Repository:    rp,
		StreakService: ss,
	}
}

func (s *HabitService) CheckHabitState(
	ctx context.Context,
	packet repositories.Packet,
	userId int,
) error {
	// Lock the process so that no race con potention
	s.Lock.Lock()
	defer s.Lock.Unlock()

	assignedTask, err := s.Repository.CountAssignedTask(
		ctx,
		repositories.CountAssignedTaskParams{
			UserID:   int64(userId),
			PacketID: packet.ID,
		},
	)

	if assignedTask == 0 {
		return nil
	}

	completionRate := float64(packet.CompletedTask) * 100.0 / float64(assignedTask)

	habits, err := s.Repository.GetHabitsByPacketId(ctx, packet.ID)
	if err != nil {
		slog.Error("Failed to get all habits")
		return nil
	}

	currentStreak, err := s.StreakService.GetCurrentStreak(userId)
	if err != nil {
		return err
	}

	for _, habit := range habits {
		switch habit.Difficulty {
		case "normal":
			if completionRate >= 50 && currentStreak >= 3 {
				s.Repository.UnlockHabit(ctx, habit.ID)
			}
		case "hard":
			if completionRate >= 70 && currentStreak >= 7 {
				s.Repository.UnlockHabit(ctx, habit.ID)
			}
		}
	}

	return nil
}

func (s *HabitService) GetAllHabits(packetId int64) ([]repositories.Habit, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	habits, err := s.Repository.GetHabitsByPacketId(context.Background(), packetId)
	if err != nil {
		slog.Error("Failed to get habits", "err", err)
		return nil, err
	}

	return habits, nil
}

func (s *HabitService) GetUnlockedHabits(packetId int64) ([]repositories.Habit, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	unlockedHabits, err := s.Repository.GetUnlockedHabitsByPacketId(
		context.Background(),
		packetId,
	)
	if err != nil {
		slog.Error("Failed to get unlocked habits habits", "err", err)
		return nil, err
	}

	return unlockedHabits, nil
}
