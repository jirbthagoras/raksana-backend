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

	g.Post("/monthly", h.handleCreateMonthlyRecap)
	g.Get("/monthly/me", h.handleGetMonthlyRecap)
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
		var completedAt string = task.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		var createdAt string = task.CreatedAt.Time.Format("2006-01-02 15:04:05")
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

	if userTasks.AssignedTask == 0 {
		stringCompletionRate = "0%"
	}
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

	var recapResponse models.AIResponseRecap
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

func (h *RecapHandler) handleCreateMonthlyRecap(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	now := time.Now()
	lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()

	if now.Day() != lastDay {
		return fiber.NewError(fiber.StatusBadRequest, "Hari ini bukan akhir bulan")
	}

	resLogs, err := h.Repository.GetLastMonthUserLogs(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get last month logs", "err", err)
		return err
	}

	var logs []models.ResponseGetLogs
	for _, l := range resLogs {
		logs = append(logs, models.ResponseGetLogs{
			Text:      l.Text,
			IsSystem:  l.IsSystem,
			CreatedAt: l.CreatedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	resHist, err := h.Repository.GetLastMonthUserHistories(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get last month histories", "err", err)
		return err
	}

	var hists []models.ResponseHistory
	for _, h := range resHist {
		hists = append(hists, models.ResponseHistory{
			Name:      h.Name,
			Type:      h.Type,
			Category:  h.Category,
			Amount:    int(h.Amount),
			CreatedAt: h.CreatedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	statistics, err := h.Repository.GetUserStatistic(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get statistics", "err", err)
		return err
	}

	req := models.RequestGetMonthlyRecap{
		Statistics: models.RequestStatistics{
			Challenges:    int(statistics.Challenges),
			Quests:        int(statistics.Quests),
			Events:        int(statistics.Events),
			Treasures:     int(statistics.Treasures),
			LongestStreak: int(statistics.LongestStreak),
		},
		Logs:      logs,
		Histories: hists,
	}

	reqMsg, err := json.Marshal(req)
	if err != nil {
		slog.Error("Something wrong while marshaling", "err", err)
		return err
	}

	cnf := helpers.NewConfig()
	model, err := configs.InitModel(h.AIClient.Genai, cnf, configs.RecapMonthly)

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(reqMsg))
	if err != nil {
		slog.Error("Something error while sending message to GenAI", "err", err)
		return err
	}

	latestRecap, err := h.Repository.GetLatestMonhtlyRecap(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get latest recap", "err", err)
	}

	if latestRecap.IsThisMonth {
		return fiber.NewError(fiber.StatusBadRequest, "Kamu sudah merekap bulan ini")
	}

	responseMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		responseMsg += fmt.Sprintf("%v\n", part)
	}

	var modelResponse models.AIResponseRecap
	err = json.Unmarshal([]byte(responseMsg), &modelResponse)
	if err != nil {
		slog.Error("Failed to parse Gemini response content", "err", err)
		return err
	}

	userTasks, err := h.Repository.CountUserTask(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to count user tasks", "err", err)
		return err
	}
	var completionRate float64 = 0.0

	completionRate = float64(userTasks.CompletedTask) * 100.0 / float64(userTasks.AssignedTask)
	stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

	if userTasks.AssignedTask == 0 {
		stringCompletionRate = "0%"
	}

	recapId, err := h.Repository.CreateMonthlyRecap(ctx, repositories.CreateMonthlyRecapParams{
		UserID:         int64(userId),
		Summary:        modelResponse.Summary,
		Tips:           modelResponse.Tips,
		GrowthRating:   modelResponse.GrowthRating,
		AssignedTask:   int32(userTasks.AssignedTask),
		CompletedTask:  int32(userTasks.CompletedTask),
		CompletionRate: stringCompletionRate,
	})
	if err != nil {
		slog.Error("Failed to create monthly recap", "err", err)
		return err
	}

	err = h.Repository.CreateRecapDetails(ctx, repositories.CreateRecapDetailsParams{
		MonthlyRecapID: recapId,
		Challenges:     statistics.Challenges,
		Quests:         statistics.Quests,
		Events:         statistics.Events,
		Treasures:      statistics.Treasures,
		LongestStreak:  statistics.LongestStreak,
	})
	if err != nil {
		slog.Error("Failed to create recap detail", "err", err)
		return err
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("failed to load timezone")
		return fmt.Errorf("failed to load timezone: %w", err)
	}
	todayDate := time.Now().In(loc).Format("2006-01")
	if modelResponse.GrowthRating == "5" || modelResponse.GrowthRating == "4" {
		logMsg := fmt.Sprintf("Saya baru saja mendapatkan growth rating %s di monthly recap %s milik saya!", modelResponse.GrowthRating, todayDate)
		err := h.JournalService.AppendLog(&models.PostLogAppend{
			Text:      logMsg,
			IsSystem:  true,
			IsPrivate: false,
		}, userId)
		if err != nil {
			return err
		}
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	var response = models.ResponseMonthlyRecap{
		Summary:        modelResponse.Summary,
		Tips:           modelResponse.Tips,
		AssignedTask:   int32(userTasks.AssignedTask),
		CompletedTask:  int32(userTasks.CompletedTask),
		CompletionRate: stringCompletionRate,
		GrowthRating:   modelResponse.GrowthRating,
		Type:           "monthly",
		CreatedAt:      todayDate,
		Challenges:     int(statistics.Challenges),
		Quests:         int(statistics.Quests),
		Events:         int(statistics.Events),
		Treasures:      int(statistics.Treasures),
		LongestStreak:  int(statistics.LongestStreak),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": response,
	})
}

func (h *RecapHandler) handleGetMonthlyRecap(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	res, err := h.Repository.GetAllMonthlyRecapsWithDetails(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get recaps", "err", err)
		return err
	}

	var monthlyRecaps []models.ResponseMonthlyRecap
	for _, r := range res {
		monthlyRecaps = append(monthlyRecaps, models.ResponseMonthlyRecap{
			Summary:        r.Summary,
			Tips:           r.Tips,
			CreatedAt:      r.RecapCreatedAt.Time.Format("2006-01-02 15:04"),
			AssignedTask:   r.AssignedTask,
			CompletedTask:  r.CompletedTask,
			CompletionRate: r.CompletionRate,
			GrowthRating:   r.GrowthRating,
			Challenges:     int(r.Challenges.Int32),
			Events:         int(r.Events.Int32),
			Treasures:      int(r.Treasures.Int32),
			Quests:         int(r.Quests.Int32),
			LongestStreak:  int(r.LongestStreak.Int32),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"monthly_recaps": monthlyRecaps,
		},
	})
}
