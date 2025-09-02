package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

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
	g.Post("/", h.handleCreateMemory)
	g.Get("/me", h.handleGetMemories)
	g.Get("/:id", h.handleGetMemoriesById)
}

func (h *MemoryHandler) handleCreateMemory(c *fiber.Ctx) error {
	req := &models.PostMemoryCreate{}
	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	err = h.Repository.CreateMemory(context.Background(), repositories.CreateMemoryParams{
		UserID:      int64(userId),
		Description: req.Description,
		FileUrl:     req.FileUrl,
	})
	if err != nil {
		slog.Error("Failed to create memory", "err", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}

func (h *MemoryHandler) handleGetMemories(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%+v", userId))

	res, err := h.Repository.GetMemoryWithParticipation(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get memories")
		return err
	}
	slog.Info(fmt.Sprintf("%+v", res))

	var memories []models.ResponseMemory

	for _, memory := range res {
		c := models.ToResponseMemory(memory)
		memories = append(memories, c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"memories": memories,
		},
	})
}

func (h *MemoryHandler) handleGetMemoriesById(c *fiber.Ctx) error {

	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	res, err := h.Repository.GetMemoryWithParticipation(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get memories")
		return err
	}

	var memories []models.ResponseMemory

	for _, memory := range res {
		c := models.ToResponseMemory(memory)
		memories = append(memories, c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"memories": memories,
		},
	})

}
