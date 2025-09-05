package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5"
)

type RecapHandler struct {
	Repository *repositories.Queries
	*configs.AIClient
	*services.JournalService
	*services.StreakService
}

func NewRecapHandler(
	r *repositories.Queries,
	ai *configs.AIClient,
	js *services.JournalService,
	ss *services.StreakService,
) *RecapHandler {
	return &RecapHandler{
		Repository:     r,
		AIClient:       ai,
		JournalService: js,
		StreakService:  ss,
	}
}

func (h *RecapHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/recap")
	g.Use(helpers.TokenMiddleware)
	g.Post("/weekly", h.handleCreateWeeklyRecap)
	g.Get("/weekly/me", h.handleGetWeeklyRecap)
	g.Get("/weekly/:id", h.handleCreateWeeklyRecap)
}

func (h *RecapHandler) handleCreateWeeklyRecap(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("Failed to get current timezone", "err", err)
		return err
	}

	today := time.Now().In(loc).Weekday()
	if today != time.Sunday {
		return fiber.NewError(fiber.StatusBadRequest, "Hari ini bukan waktu yang tepat untuk weekly recap")
	}

	cnf := helpers.NewConfig()
	aiModel, err := configs.InitModel(h.AIClient.Genai, cnf, configs.RecapWeekly)
	if err != nil {
		slog.Error("Failed to init model", "err", err)
		return err
	}

	ctx := context.Background()
	res, err := h.Repository.GetLastWeekTasks(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get last week tasks", "err", err)
		return err
	}

	var tasks []models.InputTask
	for _, task := range res {
		var completedAt string = task.UpdatedAt.Time.In(loc).Format("2006-01-02 15:04:05")
		var createdAt string = task.CreatedAt.Time.In(loc).Format("2006-01-02 15:04:05")
		tasks = append(tasks, models.InputTask{
			Name:        task.Name,
			Description: task.Description,
			Difficulty:  task.Difficulty,
			Completed:   task.Completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		})
	}

	var todayDate string = time.Now().In(loc).Format("2006-01-02")

	session := aiModel.StartChat()
	session.History = []*genai.Content{}

	var isFirstTime bool = false

	latestRecap, err := h.Repository.GetLatestRecap(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isFirstTime = true
		} else {
			slog.Error("Failed to get latest recap", "err", err)
			return err
		}
	}

	if latestRecap.CreatedAt.Time.Format("2006-01-02") == todayDate {
		return fiber.NewError(fiber.StatusBadRequest, "Anda sudah mengambil weekly recap minggu ini")
	}

	userTasks, err := h.Repository.CountUserTask(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to count user tasks", "err", err)
		return err
	}
	var completionRate float64 = 0.0

	completionRate = float64(userTasks.CompletedTask) * 100.0 / float64(userTasks.AssignedTask)
	stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

	var inputRecap = models.InputRecap{
		Date:               todayDate,
		AssignedTask:       int(userTasks.AssignedTask),
		CompletedTask:      int(userTasks.AssignedTask),
		TaskCompletionRate: stringCompletionRate,
		Tasks:              tasks,
	}

	var reqRecap = models.RequestGetRecap{
		InputRecap: inputRecap,
	}

	if !isFirstTime {
		reqRecap.PreviousRecap = latestRecap
	}

	msg, err := json.Marshal(reqRecap)
	if err != nil {
		slog.Error("Failed to casts Request Recap", "err", err)
		return err
	}

	resp, err := session.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		slog.Error("Failed to send message to generative ai", "err", err)
		return err
	}

	responseMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		responseMsg += fmt.Sprintf("%v\n", part)
	}

	var recapResponse models.AIResponsWeeklyRecap
	err = json.Unmarshal([]byte(responseMsg), &recapResponse)
	if err != nil {
		slog.Error("Failed to parse Gemini response content", "err", err)
		return err
	}

	if recapResponse.GrowthRating == "5" || recapResponse.GrowthRating == "4" {
		logMsg := fmt.Sprintf("Saya baru saja mendapatkan growth rating %s di weekly recap %s milik saya!", recapResponse.GrowthRating, todayDate)
		err := h.JournalService.AppendLog(&models.PostLogAppend{
			Text:      logMsg,
			IsSystem:  true,
			IsPrivate: false,
		}, userId)
		if err != nil {
			return err
		}
	}

	err = h.Repository.CreateWeeklyRecap(ctx, repositories.CreateWeeklyRecapParams{
		UserID:         int64(userId),
		Tips:           recapResponse.Tips,
		Summary:        recapResponse.Summary,
		AssignedTask:   int32(userTasks.AssignedTask),
		CompletedTask:  int32(userTasks.CompletedTask),
		CompletionRate: stringCompletionRate,
		GrowthRating:   recapResponse.GrowthRating,
	})
	if err != nil {
		slog.Error("Failed to create weekly recaps", "err", err)
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"recap": fiber.Map{
				"date":                 todayDate,
				"summary":              recapResponse.Summary,
				"tips":                 recapResponse.Tips,
				"assigned_tasks":       userTasks.AssignedTask,
				"completed_tasks":      userTasks.CompletedTask,
				"task_completion_rate": stringCompletionRate,
				"growth_rating":        recapResponse.GrowthRating,
			},
		},
	})
}

func (h *RecapHandler) handleGetWeeklyRecap(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetWeeklyRecaps(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get weekly recaps", "err", err)
		return err
	}

	var recaps []models.ResponseRecap

	for _, recap := range res {
		recaps = append(recaps, models.ResponseRecap{
			Summary:            recap.Summary,
			Tips:               recap.Tips,
			CompletedTask:      recap.CompletedTask,
			AssignedTask:       recap.AssignedTask,
			TaskCompletionRate: recap.CompletionRate,
			CreatedAt:          recap.CreatedAt.Time.Format("2006-01-02 15:04"),
			GrowthRating:       recap.GrowthRating,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"recaps": recaps,
		},
	})
}
