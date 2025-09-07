package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type QuestHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.PointService
	*services.JournalService
	*services.StreakService
	Mu sync.Mutex
}

func NewQuestHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
	ss *services.StreakService,
) *QuestHandler {
	return &QuestHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
		StreakService:  ss,
	}
}

func (h *QuestHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/quest")
	g.Use(helpers.TokenMiddleware)
	g.Get("/:id", h.handleGetContributedQuestDetails)
}

func (h *QuestHandler) handleContribute(c *fiber.Ctx) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	req := &models.ActivityRequest{}

	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err)
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	_, payload, err := helpers.ValidateActivityToken(req.Token)
	if err != nil {
		slog.Error("Failed to validate token")
		return fiber.NewError(fiber.StatusBadRequest, "Token Invalid")
	}

	ctx := context.Background()

	quest, err := h.Repository.GetQuestByCodeId(ctx, payload.Subject)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Quest tidak ditemukan")
		}
		slog.Error("Failed to get quest", "err", err)
		return err
	}

	contributors, err := h.Repository.CountQuestContributors(ctx, quest.ID)
	if err != nil {
		slog.Error("Failed to count", "err", err)
		return err
	}

	if len(contributors) >= int(quest.MaxContributors) {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Maksimal kontributor dari quest ini adalah %v orang", quest.MaxContributors))
	}

	_, err = h.Repository.CreateContributions(ctx, repositories.CreateContributionsParams{
		UserID:  int64(userId),
		QuestID: quest.ID,
	})
	if err != nil {
		slog.Error("Failed to count", "err", err)
		return err
	}

	_, err = h.PointService.UpdateUserPoint(int64(userId), quest.PointGain)
	if err != nil {
		slog.Error("Failed to count", "err", err)
		return err
	}

	logMsg := fmt.Sprintf("Baru saja berkontribusi dalam quest: %s dan mendapatkan poin sebesar: %v! Cek timeline ku!", quest.Name, quest.PointGain)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
			"type":    "quest",
		},
	})
}

func (h *QuestHandler) handleGetContributedQuestDetails(c *fiber.Ctx) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	contributionId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
		return err
	}

	ctx := context.Background()
	res, err := h.Repository.GetContributionDetails(ctx, int64(contributionId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Quest Not Found")
		}
		slog.Error("Failed to get contribution details", "err", err)
		return err
	}

	var quest = models.ResponseQuest{
		ID:              int(res.ID),
		Name:            res.Name,
		Description:     res.Description,
		Location:        res.Location,
		Latitude:        res.Latitude,
		Longitude:       res.Longitude,
		MaxContributors: int(res.MaxContributors),
		PointGain:       int(res.PointGain),
		ContributedAt:   res.ContributionDate.Time.In(loc).Format("2006-01-02 15:04"),
		CreatedAt:       res.CreatedAt.Time.In(loc).Format("2006-01-02 15:04"),
	}

	questContributors, err := h.Repository.CountQuestContributors(ctx, res.QuestID)
	if err != nil {
		slog.Error("Failed to get quest contributors", "err", err)
		return err
	}

	var contributors []models.Contributors
	for _, contributor := range questContributors {
		contributors = append(contributors, models.Contributors{
			ID:       int(contributor.ID),
			Username: contributor.Username,
		})
	}

	quest.Contributors = contributors

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": quest,
	})
}
