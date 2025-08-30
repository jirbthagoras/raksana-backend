package handlers

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type TaskHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.StreakService
	*services.HabitService
}

func NewTaskHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ss *services.StreakService,
	sh *services.HabitService,
) *TaskHandler {
	return &TaskHandler{
		Validator:     v,
		Repository:    r,
		StreakService: ss,
		HabitService:  sh,
	}
}

func (h *TaskHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/task")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetTodayTask)
}

func (h *TaskHandler) handleGetTodayTask(c *fiber.Ctx) error {
	ctx := context.Background()
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	todayTasks, err := h.Repository.GetTodayTasks(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get today tasks", "err", err)
		return err
	}

	if len(todayTasks) != 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": todayTasks,
		})
	}

	activePacket, err := h.Repository.GetActivePacketsByUserId(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "You have no active packet, please create a packet first")
		}
		slog.Error("Failed to get active packets", "err", err)
		return err
	}

	unlockedHabits, err := h.HabitService.GetUnlockedHabits(activePacket.ID)
	if err != nil {
		return err
	}

	var taskPerDay int = int(activePacket.TaskPerDay)
	var remainingTask int = int(activePacket.ExpectedTask) - int(activePacket.CompletedTask)

	if remainingTask < taskPerDay {
		taskPerDay = remainingTask
	}

	randomizedHabits := helpers.PickMultiple(unlockedHabits, taskPerDay)

	var tasks []models.ResponseGetTask

	for _, habit := range randomizedHabits {
		task, err := h.Repository.CreateTask(ctx, repositories.CreateTaskParams{
			HabitID:     habit.ID,
			UserID:      int64(userId),
			PacketID:    activePacket.ID,
			Name:        habit.Name,
			Description: habit.Description,
			Difficulty:  habit.Difficulty,
		})
		if err != nil {
			slog.Error("Failed to insert tasks", "err", err)
			return err
		}

		tasks = append(tasks, models.ResponseGetTask{
			Name:        task.Name,
			Description: task.Description,
			Difficulty:  task.Difficulty,
			Completed:   task.Completed,
			CreatedAt:   task.CreatedAt.Time,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"tasks": tasks,
		},
	})
}
