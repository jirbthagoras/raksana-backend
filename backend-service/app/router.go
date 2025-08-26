package app

import (
	"jirbthagoras/raksana-backend/handlers"
	"jirbthagoras/raksana-backend/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AppRouter struct {
	*handlers.AuthHandler
	*handlers.JournalHandler
	*handlers.LeaderboardHandler
	*handlers.StreakHandler
}

func NewAppRouter(
	v *validator.Validate,
	r *repositories.Queries,
	rd *redis.Client,
) *AppRouter {
	return &AppRouter{
		AuthHandler:        handlers.NewAuthHandler(v, r),
		JournalHandler:     handlers.NewJournalHandler(v, r),
		LeaderboardHandler: handlers.NewLeaderboardHandler(v, rd, r),
		StreakHandler:      handlers.NewStreakHandler(v, rd, r),
	}
}

func (r *AppRouter) RegisterRoute(router fiber.Router) {
	r.AuthHandler.RegisterRoute(router)
	r.JournalHandler.RegisterRoutes(router)
	r.LeaderboardHandler.RegisterRoutes(router)
	r.StreakHandler.RegisterRoute(router)
}
