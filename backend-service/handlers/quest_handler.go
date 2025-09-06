package handlers

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type QuestHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.PointService
	*services.JournalService
	*services.LeaderboardService
	*services.StreakService
}

func NewQuestHandler(
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
	ls *services.LeaderboardService,
	ss *services.StreakService,
) *QuestHandler {
	return &QuestHandler{
		Repository:         r,
		PointService:       ps,
		JournalService:     js,
		LeaderboardService: ls,
		StreakService:      ss,
	}
}

func (h *QuestHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/quests")
	g.Use(helpers.TokenMiddleware)
}
