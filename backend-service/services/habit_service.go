package services

import (
	"context"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
)

type HabitService struct {
	Repository *repositories.Queries
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
	packetTask, err := s.Repository.CountPacketTasks(
		ctx,
		repositories.CountPacketTasksParams{
			UserID:   int64(userId),
			PacketID: packet.ID,
		},
	)
	if err != nil {
		slog.Error("Failed sum assigned task", "err", err)
		return err
	}

	if packetTask.AssignedTask == 0 {
		return nil
	}

	completionRate := float64(packet.CompletedTask) * 100.0 / float64(packetTask.AssignedTask)

	habits, err := s.Repository.GetPacketHabits(ctx, packet.ID)
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
			if completionRate >= 70 && currentStreak >= 5 {
				s.Repository.UnlockHabit(ctx, habit.ID)
			}
		}
	}

	return nil
}

func (s *HabitService) GetAllHabits(packetId int64) ([]repositories.Habit, error) {
	habits, err := s.Repository.GetPacketHabits(context.Background(), packetId)
	if err != nil {
		slog.Error("Failed to get habits", "err", err)
		return nil, err
	}

	return habits, nil
}

func (s *HabitService) GetUnlockedHabits(packetId int64) ([]repositories.Habit, error) {
	unlockedHabits, err := s.Repository.GetPacketUnlockedHabits(
		context.Background(),
		packetId,
	)
	if err != nil {
		slog.Error("Failed to get unlocked habits habits", "err", err)
		return nil, err
	}

	return unlockedHabits, nil
}
