package handlers

import (
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PointHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*configs.AIClient
	*services.PointService
	*services.JournalService
}

func NewPointHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
) *PointHandler {
	return &PointHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
	}
}

func (h *PointHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/point")
	g.Use(helpers.TokenMiddleware)
}
