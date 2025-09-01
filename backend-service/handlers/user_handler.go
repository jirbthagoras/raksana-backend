package handlers

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Repository *repositories.Queries
	*services.UserService
	Mu sync.Mutex
}

func NewUserHandler(
	r *repositories.Queries,
	us *services.UserService,
) *UserHandler {
	return &UserHandler{
		Repository:  r,
		UserService: us,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	g1 := router.Group("/user")
	g1.Use(helpers.TokenMiddleware)

	g2 := router.Group("/profile")
	g2.Use(helpers.TokenMiddleware)
	g2.Get("/me", h.handleGetProfile)
	g2.Get("/:id", h.handleGetProfileById)
}

func (h *UserHandler) handleGetProfile(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	profile, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": profile,
	})
}

func (h *UserHandler) handleGetProfileById(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}
	profile, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": profile,
	})
}
