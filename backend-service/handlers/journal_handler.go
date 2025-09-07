package handlers

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type JournalHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.JournalService
	*services.StreakService
}

func NewJournalHandler(
	v *validator.Validate,
	r *repositories.Queries,
	s *services.JournalService,
	ss *services.StreakService,
) *JournalHandler {
	return &JournalHandler{
		Repository:     r,
		Validator:      v,
		JournalService: s,
		StreakService:  ss,
	}
}

func (h *JournalHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/log")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleAppendJournal)
	g.Get("/", h.handleGetLogs)
	g.Get("/:id", h.handleGetLogsByUserId)
}

func (h *JournalHandler) handleAppendJournal(c *fiber.Ctx) error {
	req := &models.PostLogAppend{}
	req.IsPrivate = false
	req.IsSystem = false
	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
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

	err = h.JournalService.AppendLog(req, userId)
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(context.Background(), int64(userId))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}

func (h *JournalHandler) handleGetLogs(c *fiber.Ctx) error {
	isPrivateParam := c.Query("is_private", "false")
	isPrivate, err := strconv.ParseBool(isPrivateParam)
	if err != nil {
		isPrivate = false
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetLogs(context.Background(), repositories.GetLogsParams{
		UserID:    int64(userId),
		IsPrivate: isPrivate,
	})

	var logs = []models.ResponseGetLogs{}
	for _, log := range res {
		logs = append(logs, models.ResponseGetLogs{
			Text:      log.Text,
			IsSystem:  log.IsSystem,
			IsPrivate: log.IsPrivate,
			CreatedAt: log.CreatedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"logs": logs,
		},
	})
}

func (h *JournalHandler) handleGetLogsByUserId(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	logs, err := h.JournalService.GetLogs(userId, false)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"logs": logs,
		},
	})

}
