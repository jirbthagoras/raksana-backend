package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type TaskHandler struct {
	Repository *repositories.Queries
	*services.StreakService
	*services.HabitService
	*services.JournalService
	*services.ExpService
	Mu sync.Mutex
}

func NewTaskHandler(
	r *repositories.Queries,
	ss *services.StreakService,
	sh *services.HabitService,
	js *services.JournalService,
	es *services.ExpService,
) *TaskHandler {
	return &TaskHandler{
		Repository:     r,
		StreakService:  ss,
		HabitService:   sh,
		JournalService: js,
		ExpService:     es,
	}
}

func (h *TaskHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/task")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetTodayTask)
	g.Put("/:id", h.handleCompleteTask)
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

	var tasks []models.ResponseGetTask

	if len(todayTasks) != 0 {
		for _, task := range todayTasks {
			tasks = append(tasks, models.ResponseGetTask{
				Id:          int(task.ID),
				Name:        task.Name,
				Description: task.Description,
				Difficulty:  task.Difficulty,
				Completed:   task.Completed,
				CreatedAt:   task.CreatedAt.Time.Format("2006-01-02 15:04"),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": fiber.Map{
				"tasks": tasks,
			},
		})
	}

	activePacket, err := h.Repository.GetUserActivePackets(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"data": fiber.Map{
					"tasks":   []string{},
					"message": "Kamu tidak memiliki packet aktif.",
				},
			})
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
			Id:          int(task.ID),
			Name:        task.Name,
			Description: task.Description,
			Difficulty:  task.Difficulty,
			Completed:   task.Completed,
			CreatedAt:   task.CreatedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"tasks": tasks,
		},
	})
}

func (h *TaskHandler) handleCompleteTask(c *fiber.Ctx) error {
	// Lock the process because ts can happens so multiple times at the same time hehe
	h.Mu.Lock()
	defer h.Mu.Unlock()

	taskId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to parse id from route parameters")
		return err
	}

	ctx := context.Background()

	res, err := h.Repository.GetTaskById(ctx, int64(taskId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Task not found")
		}
		slog.Error("Failed to get specific task", "err", err)
		return err
	}

	if res.Completed {
		return fiber.NewError(fiber.StatusBadRequest, "Task already completed")
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	activePacket, err := h.Repository.GetUserActivePackets(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "You have no active packet, please create a packet first")
		}
		slog.Error("Failed to get active packets", "err", err)
		return err
	}

	task, err := h.Repository.CompleteTask(ctx, repositories.CompleteTaskParams{
		UserID: int64(userId),
		ID:     int64(taskId),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Task is not valid")
		}
		slog.Error("Failed to update task status")
		return err
	}

	err = h.Repository.IncreasePacketCompletedTask(ctx, activePacket.ID)
	if err != nil {
		slog.Error("Failed to update packet")
		return err
	}

	activePacket.CompletedTask++
	if activePacket.CompletedTask >= activePacket.ExpectedTask {
		err = h.Repository.CompletePacket(ctx, activePacket.ID)
		if err != nil {
			slog.Error("Failed to complete packet", "err", err)
			return err
		}

		packetTask, err := h.Repository.CountPacketTasks(ctx,
			repositories.CountPacketTasksParams{
				UserID:   int64(userId),
				PacketID: activePacket.ID,
			},
		)
		if err != nil {
			slog.Error("Failed to count assigned task", "err", err)
			return err
		}

		var completionRate float64 = 0.0

		if packetTask.AssignedTask != 0 {
			completionRate = float64(activePacket.CompletedTask) * 100.0 / float64(packetTask.AssignedTask)
		}

		logMsg := fmt.Sprintf("Aku baru saja menyelesaikan packet %s! Dengan completion rate: %v", activePacket.Name, completionRate) + "%"
		err = h.JournalService.AppendLog(&models.PostLogAppend{
			Text:      logMsg,
			IsSystem:  true,
			IsPrivate: false,
		}, userId)
	}

	err = h.HabitService.CheckHabitState(ctx, activePacket, userId)
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	expGain, err := helpers.CheckExpGain(task.Difficulty)
	if err != nil {
		slog.Error("Failed to get exp gain", "err", err)
		return err
	}

	todayTask, err := h.Repository.GetTodayTasks(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get today tasks", "err", err)
		return err
	}

	if len(todayTask) <= 0 {
		err := h.JournalService.AppendLog(&models.PostLogAppend{
			Text:      "Aku baru saja menyelesaikan semua task hari ini!",
			IsSystem:  true,
			IsPrivate: false,
		}, userId)
		if err != nil {
			return err
		}
	}

	levelUp, level, err := h.ExpService.IncreaseExp(userId, expGain)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"message":       "sucessfully completed task",
			"leveled_up":    levelUp,
			"current_level": level,
		},
	})
}
