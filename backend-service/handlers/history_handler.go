package handlers

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type HistoryHandler struct {
	Repository *repositories.Queries
}

func NewHistoryHandler(
	r *repositories.Queries,
) *HistoryHandler {
	return &HistoryHandler{
		Repository: r,
	}
}

func (h *HistoryHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/history")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetHistories)
}

func (h *HistoryHandler) handleGetHistories(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetUserHistories(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get histories")
		return nil
	}

	var histories []models.ResponseHistory
	for _, h := range res {
		histories = append(histories, models.ResponseHistory{
			Name:      h.Name,
			Type:      h.Type,
			Category:  h.Category,
			Amount:    int(h.Amount),
			CreatedAt: h.CreatedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"histories": histories,
		},
	})
}
