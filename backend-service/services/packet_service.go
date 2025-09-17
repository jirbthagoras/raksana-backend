package services

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type PacketService struct {
	Repository *repositories.Queries
}

func NewPacketService(r *repositories.Queries) *PacketService {
	return &PacketService{
		Repository: r,
	}
}

func (s *PacketService) GetALlPackets(userId int64) ([]models.ResponseGetPacket, error) {
	res, err := s.Repository.GetAllPackets(context.Background(), userId)
	if err != nil {
		slog.Error("Failed to get all packets", "err", err)
		return nil, err
	}

	var packets []models.ResponseGetPacket
	for _, packet := range res {
		packetTask, err := s.Repository.CountPacketTasks(context.Background(), repositories.CountPacketTasksParams{
			UserID:   int64(userId),
			PacketID: packet.ID,
		})
		if err != nil {
			slog.Error("Failed to count packet assigned task", "err", err)
			return nil, err
		}

		completionRate := float64(packet.CompletedTask) * 100.0 / float64(packetTask.AssignedTask)
		if packetTask.AssignedTask == 0 {
			completionRate = 0
		}

		stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

		packets = append(packets, models.ResponseGetPacket{
			Id:             int32(packet.ID),
			Name:           packet.Name,
			Description:    packet.Description,
			Target:         packet.Target,
			ExpectedTask:   packet.ExpectedTask,
			CompletedTask:  packet.CompletedTask,
			CompletionRate: stringCompletionRate,
			AssignedTask:   int32(packetTask.AssignedTask),
			TaskPerDay:     packet.TaskPerDay,
			CreatedAt:      packet.CreatedAt.Time.Format("2006-01-02 15:04"),
			Completed:      packet.Completed,
		})
	}

	return packets, nil
}

func (s *PacketService) GetPacketDetail(id int) (models.ResponsePacketDetail, error) {
	ctx := context.Background()

	var habitDetail models.ResponsePacketDetail

	packetDetails, err := s.Repository.GetPacketDetail(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return habitDetail, fiber.NewError(fiber.StatusBadRequest, "No packet found ")
		}
		slog.Error("Failed to get packet detail", "err", err)
		return habitDetail, err
	}

	habits, err := s.Repository.GetPacketHabits(ctx, packetDetails.PacketID)
	if err != nil {
		slog.Error("Failed to get packet's habits", "err", err)
		return habitDetail, err
	}

	var packetHabits []models.ResponsePacketDetailHabit

	for _, habit := range habits {
		expGain, err := helpers.CheckExpGain(habit.Difficulty)
		if err != nil {
			return habitDetail, err
		}

		packetHabits = append(packetHabits, models.ResponsePacketDetailHabit{
			Name:        habit.Name,
			Description: habit.Description,
			Difficulty:  habit.Difficulty,
			Locked:      habit.Locked,
			ExpGain:     int32(expGain),
		})
	}
	packetTask, err := s.Repository.CountPacketTasks(context.Background(), repositories.CountPacketTasksParams{
		UserID:   int64(packetDetails.UserID),
		PacketID: packetDetails.PacketID,
	})
	if err != nil {
		slog.Error("Failed to count packet assigned task", "err", err)
		return habitDetail, err
	}

	completionRate := float64(packetTask.CompletedTask) / float64(packetTask.AssignedTask) * 100.0
	if packetTask.AssignedTask == 0 {
		completionRate = 0
	}

	stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

	habitDetail = models.ResponsePacketDetail{
		PacketID:           packetDetails.PacketID,
		PacketName:         packetDetails.PacketName,
		Username:           packetDetails.Username,
		Target:             packetDetails.Target,
		Description:        packetDetails.Description,
		CompletedTask:      packetDetails.CompletedTask,
		ExpectedTask:       packetDetails.ExpectedTask,
		AssignedTask:       int32(packetTask.AssignedTask),
		TaskPerDay:         packetDetails.TaskPerDay,
		TaskCompletionRate: stringCompletionRate,
		Completed:          packetDetails.Completed,
		CreatedAt:          packetDetails.CreatedAt.Time.Format("2006-01-02 15:04"),
		Habits:             packetHabits,
	}

	return habitDetail, nil
}
