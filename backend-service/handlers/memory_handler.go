package handlers

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MemoryHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
}

func NewMemoryHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *MemoryHandler {
	return &MemoryHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *MemoryHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/memory")
	g.Use(helpers.TokenMiddleware)
}

func (h *MemoruHandler) handleCreateMemory(c *fiber.Ctx) error {
	return nil
}
