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
	*handlers.TaskHandler
}

func NewAppRouter(
	v *validator.Validate,
	r *repositories.Queries,
	rd *redis.Client,
) *AppRouter {
	journalService := services.NewJournalService(r)
	streakService := services.NewStreakService(rd, r)
	habitService := services.NewHabitService(r, streakService)
	expService := services.NewExpService(r, journalService)

	cnf := helpers.NewConfig()
	aiClient := configs.InitAiClient(cnf)

	return &AppRouter{
		AuthHandler:        handlers.NewAuthHandler(v, r),
		JournalHandler:     handlers.NewJournalHandler(v, journalService),
		LeaderboardHandler: handlers.NewLeaderboardHandler(v, rd, r),
		StreakHandler:      handlers.NewStreakHandler(v, rd, streakService),
		PacketHandler:      handlers.NewPacketHandler(v, r, aiClient, journalService),
		TaskHandler:        handlers.NewTaskHandler(v, r, streakService, habitService, journalService, expService),
	}
}

func (r *AppRouter) RegisterRoute(router fiber.Router) {
	r.AuthHandler.RegisterRoutes(router)
	r.JournalHandler.RegisterRoutes(router)
	r.LeaderboardHandler.RegisterRoutes(router)
	r.StreakHandler.RegisterRoutes(router)
	r.PacketHandler.RegisterRoutes(router)
	r.TaskHandler.RegisterRoutes(router)
}
