package main

import (
	"context"
	"jirbthagoras/raksana-backend/app"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// init repo
	repository := repositories.New(conn)
	validator := validator.New()

	// init fcm
	_ = app.InitFCMClient()

	api := server.Group("/api")

	router := app.NewAppRouter(validator, repository)
	router.RegisterRoute(api)

	go func() {
		if err := server.Listen(":3000"); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down gracefully.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.ShutdownWithContext(ctx); err != nil {
		slog.Error("Error shutting down fiber server", "err", err.Error())
	}

	conn.Close()
	slog.Info("Server stopped successfully.")
}
