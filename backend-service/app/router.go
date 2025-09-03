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
	*handlers.UserHandler
	*handlers.FileHandler
	*handlers.MemoryHandler
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
	packetService := services.NewPacketService(r)
	userService := services.NewUserService(r)

	cnf := helpers.NewConfig()
	aiClient := configs.InitAiClient(cnf)
	awsClient := configs.InitAWSClient(cnf)

	return &AppRouter{
		AuthHandler:        handlers.NewAuthHandler(v, r),
		JournalHandler:     handlers.NewJournalHandler(v, journalService),
		LeaderboardHandler: handlers.NewLeaderboardHandler(rd, r),
		StreakHandler:      handlers.NewStreakHandler(rd, streakService),
		PacketHandler:      handlers.NewPacketHandler(v, r, aiClient, journalService, packetService),
		TaskHandler:        handlers.NewTaskHandler(r, streakService, habitService, journalService, expService),
		UserHandler:        handlers.NewUserHandler(v, r, userService, awsClient),
		FileHandler:        handlers.NewFileHandler(v, awsClient),
		MemoryHandler:      handlers.NewMemoryHandler(v, r, awsClient),
	}
}

func (r *AppRouter) RegisterRoute(router fiber.Router) {
	r.AuthHandler.RegisterRoutes(router)
	r.JournalHandler.RegisterRoutes(router)
	r.LeaderboardHandler.RegisterRoutes(router)
	r.StreakHandler.RegisterRoutes(router)
	r.PacketHandler.RegisterRoutes(router)
	r.TaskHandler.RegisterRoutes(router)
	r.UserHandler.RegisterRoutes(router)
	r.FileHandler.RegisterRoutes(router)
	r.MemoryHandler.RegisterRoutes(router)
}
