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
	*handlers.QuestHandler
	*handlers.EventHandler
	*handlers.ScanHandler
	*handlers.ActivityHandler
	*handlers.HistoryHandler
	*handlers.PointHandler
	*handlers.RegionHandler
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
	leaderboardService := services.NewLeaderboardService(rd)
	userService := services.NewUserService(r, streakService, leaderboardService)
	memoryService := services.NewMemoryService(r)
	pointService := services.NewPointService(r, leaderboardService)
	fileService := services.NewFileService(awsClient)

	treasureHandler := handlers.NewTreasureHandler(v, r, pointService, journalService, streakService)
	questHandler := handlers.NewQuestHandler(v, r, pointService, journalService, streakService)
	eventHandler := handlers.NewEventHandler(v, r, pointService, journalService, streakService)

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
		TreasureHandler:    treasureHandler,
		QuestHandler:       questHandler,
		EventHandler:       eventHandler,
		ScanHandler:        handlers.NewScanHandler(v, r, treasureHandler, questHandler, eventHandler, awsClient, aiClient),
		ActivityHandler:    handlers.NewActivityHandler(v, r),
		HistoryHandler:     handlers.NewHistoryHandler(r),
		PointHandler:       handlers.NewPointHandler(v, r, pointService, journalService),
		RegionHandler:      handlers.NewRegionHandler(v, r),
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
	r.QuestHandler.RegisterRoutes(router)
	r.EventHandler.RegisterRoutes(router)
	r.ScanHandler.RegisterRoutes(router)
	r.ActivityHandler.RegisterRoutes(router)
	r.HistoryHandler.RegisterRoutes(router)
	r.PointHandler.RegisterRoutes(router)
	r.RegionHandler.RegisterRoutes(router)
}
