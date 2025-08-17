package main

import (
	"jirbthagoras/raksana-backend/helpers"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New(fiber.Config{
		ErrorHandler: helpers.ErrorHandler,
	})

	if err := server.Listen(":3000"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
