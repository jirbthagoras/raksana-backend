package handlers

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ActivityHandler struct {
	Repository *repositories.Queries
	Validator  *validator.Validate
}

func NewActivityHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *ActivityHandler {
	return &ActivityHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *ActivityHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/activity")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetActivityMap)
	g.Get("/:id", h.handleGetActivityMapByUserId)
}

func (h *ActivityHandler) handleGetActivityMap(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	contributionsRes, err := h.Repository.GetUserContributions(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user contributions", "err", err)
		return err
	}

	var contributions []models.ResponseContributions
	for _, c := range contributionsRes {
		contributions = append(contributions, models.ResponseContributions{
			Id:          c.ContributionID,
			Name:        c.Name,
			Description: c.Description,
			Latitude:    c.Latitude,
			Longitude:   c.Longitude,
			PointGain:   float64(c.PointGain),
		})
	}

	attendanceRes, err := h.Repository.GetUserAttendances(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user contributions", "err", err)
		return err
	}

	var attendances []models.ResponseAttendances
	for _, a := range attendanceRes {
		attendances = append(attendances, models.ResponseAttendances{
			Id:          a.AttendanceID,
			Name:        a.DetailName,
			Description: a.DetailDescription,
			Latitude:    a.Latitude,
			Longitude:   a.Longitude,
			PointGain:   a.PointGain,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"contributions": contributions,
			"attendances":   attendances,
		},
	})
}

func (h *ActivityHandler) handleGetActivityMapByUserId(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	ctx := context.Background()

	contributionsRes, err := h.Repository.GetUserContributions(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user contributions", "err", err)
		return err
	}

	var contributions []models.ResponseContributions
	for _, c := range contributionsRes {
		contributions = append(contributions, models.ResponseContributions{
			Id:          c.ContributionID,
			Name:        c.Name,
			Description: c.Description,
			Latitude:    c.Latitude,
			Longitude:   c.Longitude,
			PointGain:   float64(c.PointGain),
		})
	}

	attendanceRes, err := h.Repository.GetUserAttendances(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user contributions", "err", err)
		return err
	}

	var attendances []models.ResponseAttendances
	for _, a := range attendanceRes {
		attendances = append(attendances, models.ResponseAttendances{
			Id:          a.AttendanceID,
			Name:        a.DetailName,
			Description: a.DetailDescription,
			Latitude:    a.Latitude,
			Longitude:   a.Longitude,
			PointGain:   a.PointGain,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"contributions": contributions,
			"attendances":   attendances,
		},
	})
}
