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

type CheckHabitStateReturn struct {
	PacketName string `json:"packet_name,omitempty"`
	Difficulty string `json:"difficulty,omitempty"`
	IsUnlock   bool   `json:"is_unlocked"`
}

func (s *HabitService) CheckHabitState(
	ctx context.Context,
	packet repositories.Packet,
	userId int,
) (CheckHabitStateReturn, error) {

	var res = CheckHabitStateReturn{
		IsUnlock: false,
	}

	packetTask, err := s.Repository.CountPacketTasks(
		ctx,
		repositories.CountPacketTasksParams{
			UserID:   int64(userId),
			PacketID: packet.ID,
		},
	)

	if err != nil {
		slog.Error("Failed sum assigned task", "err", err)
		return res, err
	}

	if packetTask.AssignedTask == 0 {
		return res, nil
	}

	completionRate := int(float64(packet.CompletedTask) * 100.0 / float64(packetTask.AssignedTask))

	habits, err := s.Repository.GetLockedHabits(ctx, packet.ID)
	if err != nil {
		slog.Error("Failed to get all habits")
		return res, nil
	}

	currentStreak, err := s.StreakService.GetCurrentStreak(ctx, int64(userId))
	if err != nil {
		return res, err
	}

	for _, habit := range habits {
		switch habit.Difficulty {
		case "normal":
			if completionRate >= 50 && currentStreak >= 3 {
				s.Repository.UnlockHabit(ctx, habit.ID)
				res.PacketName = packet.Name
				res.IsUnlock = true
				res.Difficulty = "sedang"
			}
		case "hard":
			if completionRate >= 70 && currentStreak >= 5 {
				s.Repository.UnlockHabit(ctx, habit.ID)
				res.PacketName = packet.Name
				res.IsUnlock = true
				res.Difficulty = "sulit"
			}
		}
	}

	return res, nil
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
