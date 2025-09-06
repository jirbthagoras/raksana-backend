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
	*handlers.MemoryHandler
	*handlers.RecapHandler
	*handlers.ChallengeHandler
	*handlers.TreasureHandler
}

func NewAppRouter(
	v *validator.Validate,
	r *repositories.Queries,
	rd *redis.Client,
) *AppRouter {
	cnf := helpers.NewConfig()
	aiClient := configs.InitAiClient(cnf)
	awsClient := configs.InitAWSClient(cnf)

	journalService := services.NewJournalService(r)
	streakService := services.NewStreakService(rd, r)
	habitService := services.NewHabitService(r, streakService)
	expService := services.NewExpService(r, journalService)
	packetService := services.NewPacketService(r)
	userService := services.NewUserService(r, streakService)
	leaderboardService := services.NewLeaderboardService(rd)
	memoryService := services.NewMemoryService(r)
	pointService := services.NewPointService(r, leaderboardService)
	fileService := services.NewFileService(awsClient)

	return &AppRouter{
		AuthHandler:        handlers.NewAuthHandler(v, r, leaderboardService),
		JournalHandler:     handlers.NewJournalHandler(v, r, journalService, streakService),
		LeaderboardHandler: handlers.NewLeaderboardHandler(leaderboardService),
		StreakHandler:      handlers.NewStreakHandler(rd, streakService),
		PacketHandler:      handlers.NewPacketHandler(v, r, aiClient, journalService, packetService, streakService),
		TaskHandler:        handlers.NewTaskHandler(r, streakService, habitService, journalService, expService),
		UserHandler:        handlers.NewUserHandler(v, r, userService, leaderboardService, fileService, awsClient),
		MemoryHandler:      handlers.NewMemoryHandler(v, r, memoryService, fileService, streakService, awsClient),
		RecapHandler:       handlers.NewRecapHandler(r, aiClient, journalService, streakService),
		ChallengeHandler:   handlers.NewChallengeHandler(v, r, memoryService, pointService, journalService, fileService, streakService),
		TreasureHandler:    handlers.NewTreasureHandler(v, r, pointService, journalService),
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
	r.MemoryHandler.RegisterRoutes(router)
	r.RecapHandler.RegisterRoutes(router)
	r.ChallengeHandler.RegisterRoutes(router)
	r.TreasureHandler.RegisterRoutes(router)
}
