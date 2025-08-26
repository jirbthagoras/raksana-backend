package handlers

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type JournalHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
}

func NewJournalHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *JournalHandler {
	return &JournalHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *JournalHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/log")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleAppendJournal)
	g.Get("/", h.handleGetLogs)
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

	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("err", err)
		return err
	}

	if err != nil {
		slog.Error("Failed to convert string into integer", "err", err.Error())
		return err
	}

	_, err = h.Repository.CreateLog(context.Background(), repositories.CreateLogParams{
		UserID:    int64(id),
		Text:      req.Text,
		IsMarked:  false,
		IsSystem:  req.IsSystem,
		IsPrivate: req.IsPrivate,
	})
	if err != nil {
		slog.Error("Failed to append log", "err", err.Error())
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}

func (h *JournalHandler) handleGetLogs(c *fiber.Ctx) error {
	// Get as string with default
	isPrivateParam := c.Query("is_private", "false") // default false
	isSystemParam := c.Query("is_system", "false")   // default false
	isMarkedParam := c.Query("is_marked", "false")

	// Convert string -> bool
	isPrivate, err := strconv.ParseBool(isPrivateParam)
	if err != nil {
		isPrivate = false
	}

	isSystem, err := strconv.ParseBool(isSystemParam)
	if err != nil {
		isSystem = false
	}

	isMarked, err := strconv.ParseBool(isMarkedParam)
	if err != nil {
		isMarked = false
	}

	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	result, err := h.Repository.GetLogs(context.Background(), repositories.GetLogsParams{
		UserID:    int64(id),
		IsMarked:  isMarked,
		IsSystem:  isSystem,
		IsPrivate: isPrivate,
	})

	logs := []models.ResponseGetLogs{}
	for _, log := range result {
		logs = append(logs, models.ResponseGetLogs{
			Text:      log.Text,
			IsMarked:  log.IsMarked,
			IsSystem:  log.IsSystem,
			IsPrivate: log.IsPrivate,
			CreatedAt: log.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"logs": logs,
		},
	})
}
