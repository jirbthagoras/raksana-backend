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
	Validator *validator.Validate
	*services.JournalService
}

func NewJournalHandler(
	v *validator.Validate,
	s *services.JournalService,
) *JournalHandler {
	return &JournalHandler{
		Validator:      v,
		JournalService: s,
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

	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	err = h.JournalService.AppendLog(req, id)
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

	id, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	logs, err := h.Repository.GetLogs(context.Background(), repositories.GetLogsParams{
		UserID:    int64(id),
		IsPrivate: isPrivate,
	})

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
