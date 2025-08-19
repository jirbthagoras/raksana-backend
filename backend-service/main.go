package main

import (
	"jirbthagoras/raksana-backend/app"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})

	conn := app.GetConnection()

	repository := repositories.New(conn)
	validator := validator.New()

	api := server.Group("/api")

	router := app.NewAppRouter(validator, repository)
	router.RegisterRoute(api)

	if err := server.Listen(":3000"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
