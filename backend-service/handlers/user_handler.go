package handlers

import (
	"context"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.StreakService
	*services.HabitService
	*services.JournalService
	*services.ExpService
	Mu sync.Mutex
}

func NewUserHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *UserHandler {
	return &UserHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	g1 := router.Group("/user")
	g1.Use(helpers.TokenMiddleware)

	g2 := router.Group("/profile")
	g2.Use(helpers.TokenMiddleware)
	g2.Get("/", h.handleGetProfile)
}

func (h *UserHandler) handleGetProfile(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetUserProfileStatistic(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get user profile and statistics")
		return err
	}

	tasks, err := h.Repository.CountUserTask(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to count", "err", err)
		return err
	}
	var completionRate float64 = 0.0

	if tasks.AssignedTask != 0 || tasks.CompletedTask != 0 {
		completionRate = float64(tasks.CompletedTask) * 100.0 / float64(tasks.AssignedTask)
	}

	stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

	profile := models.ResponseGetUserProfileStatistic{
		Name:               res.Name,
		Username:           res.Username,
		Email:              res.Email,
		CurrentExp:         res.CurrentExp,
		ExpNeeded:          res.ExpNeeded,
		Level:              res.Level,
		Points:             res.Points,
		ProfileUrl:         res.ProfileUrl,
		Challenges:         res.Challenges,
		Events:             res.Events,
		Quests:             res.Quests,
		Treasures:          res.Treasures,
		TaskCompletionRate: stringCompletionRate,
		CompletedTask:      int32(tasks.CompletedTask),
		AssignedTask:       int32(tasks.AssignedTask),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": profile,
	})
}
