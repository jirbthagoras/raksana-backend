package app

import (
	"jirbthagoras/raksana-backend/handlers"
	"jirbthagoras/raksana-backend/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AppRouter struct {
	*handlers.AuthHandler
}

func NewAppRouter(v *validator.Validate, r *repositories.Queries) *AppRouter {
	return &AppRouter{
		AuthHandler: handlers.NewAuthHandler(v, r),
	}
}

func (r *AppRouter) RegisterRoute(router fiber.Router) {
	r.AuthHandler.RegisterRoute(router)
}
