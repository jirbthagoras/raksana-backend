package app

import (
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/handlers"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AppRouter struct {
	*handlers.AuthHandler
	*handlers.JournalHandler
	*handlers.LeaderboardHandler
	*handlers.StreakHandler
	*handlers.PacketHandler
}

func NewAppRouter(
	v *validator.Validate,
	r *repositories.Queries,
	rd *redis.Client,
) *AppRouter {
	journalService := services.NewJournalService(r)
	streakService := services.NewStreakService(rd, r)

	cnf := helpers.NewConfig()
	aiClient := configs.InitAiClient(cnf)

	return &AppRouter{
		AuthHandler:        handlers.NewAuthHandler(v, r),
		JournalHandler:     handlers.NewJournalHandler(v, journalService),
		LeaderboardHandler: handlers.NewLeaderboardHandler(v, rd, r),
		StreakHandler:      handlers.NewStreakHandler(v, rd, streakService),
		PacketHandler:      handlers.NewPacketHandler(v, r, aiClient),
	}
}

func (r *AppRouter) RegisterRoute(router fiber.Router) {
	r.AuthHandler.RegisterRoutes(router)
	r.JournalHandler.RegisterRoutes(router)
	r.LeaderboardHandler.RegisterRoutes(router)
	r.StreakHandler.RegisterRoutes(router)
	r.PacketHandler.RegisterRoutes(router)
}
