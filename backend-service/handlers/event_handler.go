package handlers

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.PointService
	*services.JournalService
	*services.StreakService
	Mu sync.Mutex
}

func NewEventHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
	ss *services.StreakService,
) *EventHandler {
	return &EventHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
		StreakService:  ss,
	}
}

func (h *EventHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/event ")
	g.Use(helpers.TokenMiddleware)
}
