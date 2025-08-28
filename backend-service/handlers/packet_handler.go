package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
)

type PacketHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*configs.AIClient
}

func NewPacketHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ai *configs.AIClient,
) *PacketHandler {
	return &PacketHandler{
		Validator:  v,
		Repository: r,
		AIClient:   ai,
	}
}

func (h *PacketHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/packet")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleGeneratePacket)
}

func (h *PacketHandler) handleGeneratePacket(c *fiber.Ctx) error {
	req := &models.PostPacketCreate{}
	ctx := context.Background()
	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse body")
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token is not attached")
	}

	msg := fmt.Sprintf("Deskripsi: %s, target: %s", req.Description, req.Target)

	cnf := helpers.NewConfig()
	aiModel, err := configs.InitModel(h.AIClient.Genai, cnf, configs.Ecoach)
	if err != nil {
		return err
	}

	session := aiModel.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		slog.Error("Failed to send message to generative ai", "err", err)
		return err
	}

	responseMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		responseMsg += fmt.Sprintf("%v\n", part)
	}

	var ecoachResponse models.EcoachCreatePacketResponse
	err = json.Unmarshal([]byte(responseMsg), &ecoachResponse)
	if err != nil {
		slog.Error("Failed to parse Gemini response content", "err", err)
		return err
	}

	packetId, err := h.Repository.CreatePacket(ctx, repositories.CreatePacketParams{
		UserID:       int64(userId),
		Name:         ecoachResponse.Name,
		Target:       req.Target,
		Description:  req.Description,
		ExpectedTask: int32(ecoachResponse.ExpectedTask),
		TaskPerDay:   int32(ecoachResponse.TaskPerDay),
	})
	if err != nil {
		slog.Error("Failed to insert row into packets", "err", err)
		return err
	}

	for _, habit := range ecoachResponse.Habits {

		var locked bool = false

		if habit.Difficulty == "easy" {
			locked = true
		}

		habit.Difficulty = strings.ToLower(strings.TrimSpace(habit.Difficulty))

		_, err = h.Repository.CreateHabit(ctx, repositories.CreateHabitParams{
			PacketID:    packetId,
			Name:        habit.Name,
			Description: habit.Description,
			Difficulty:  habit.Description,
			Locked:      locked,
		})
		if err != nil {
			slog.Error("Failed to insert row into habits", "err", err)
			return err
		}

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"packet": ecoachResponse,
		},
	})
}
