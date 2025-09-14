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

type RegionHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
}

func NewRegionHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *RegionHandler {
	return &RegionHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *RegionHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/region")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetRegions)
}

func (h *RegionHandler) handleGetRegions(c *fiber.Ctx) error {
	ctx := context.Background()

	res, err := h.Repository.GetAllRegions(ctx)
	if err != nil {
		slog.Error("Failed to get regions", "err", err)
		return err
	}

	if len(res) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"data": fiber.Map{
				"regions": []string{},
			},
		})
	}

	var regions []models.ResponseRegion
	for _, region := range res {
		regions = append(regions, models.ResponseRegion{
			Id:         int(region.ID),
			Name:       region.Name,
			Location:   region.Location,
			Latitude:   region.Latitude,
			Longitude:  region.Longitude,
			TreeAmount: int(region.TreeAmount),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"regions": regions,
		},
	})
}
