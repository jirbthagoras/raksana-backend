package main

import (
	"jirbthagoras/raksana-backend/app"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New(fiber.Config{
		ErrorHandler: exceptions.ErrorHandler,
	})

	// open connection
	conn := app.GetConnection()
	defer conn.Close()

	// redis
	redisConn := app.NewRedisClient()

	// init repo
	repository := repositories.New(conn)
	validator := validator.New()

	api := server.Group("/api")

	router := app.NewAppRouter(validator, repository, redisConn)
	router.RegisterRoute(api)

	// go func() {
	if err := server.Listen(":3000"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit
	// slog.Info("Shutting down gracefully.")
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// if err := server.ShutdownWithContext(ctx); err != nil {
	// 	slog.Error("Error shutting down fiber server", "err", err.Error())
	// }
	//
	// slog.Info("Server stopped successfully.")
}
