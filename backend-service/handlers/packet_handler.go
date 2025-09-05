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
	"jirbthagoras/raksana-backend/services"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
)

type PacketHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*configs.AIClient
	*services.JournalService
	*services.PacketService
	*services.StreakService
}

func NewPacketHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ai *configs.AIClient,
	js *services.JournalService,
	ps *services.PacketService,
	ss *services.StreakService,
) *PacketHandler {
	return &PacketHandler{
		Validator:      v,
		Repository:     r,
		AIClient:       ai,
		JournalService: js,
		PacketService:  ps,
		StreakService:  ss,
	}
}

func (h *PacketHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/packet")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleGeneratePacket)
	g.Get("/me", h.handleGetAllPackets)
	g.Get("/:id", h.handleGetPacketById)
	g.Get("/detail/:id", h.handleGetPacketDetail)
}

func (h *PacketHandler) handleGetAllPackets(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	packets, err := h.PacketService.GetALlPackets(int64(userId))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"packets": packets,
		},
	})
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
		return err
	}

	result, err := h.Repository.CountUserActivePackets(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to count active packets", "err", err)
		return err
	}

	if result != 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Anda sudah memiliki beberapa packet aktif!")
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
		var locked bool = true
		var weight int

		switch habit.Difficulty {
		case "easy":
			locked = false
			weight = 70
		case "normal":
			weight = 25
		case "hard":
			weight = 5
		}

		_, err = h.Repository.CreateHabit(ctx, repositories.CreateHabitParams{
			PacketID:    packetId,
			Name:        habit.Name,
			Description: habit.Description,
			Difficulty:  habit.Difficulty,
			Locked:      locked,
			Weight:      int32(weight),
		})
		if err != nil {
			slog.Error("Failed to insert row into habits", "err", err)
			return err
		}
	}

	logMsg := fmt.Sprintf("Baru saja membuat packet baru dengan nama: %s ayo dicek!", ecoachResponse.Name)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		IsSystem:  true,
		IsPrivate: false,
		Text:      logMsg,
	}, userId)
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"packet": ecoachResponse,
		},
	})
}

func (h *PacketHandler) handleGetPacketDetail(c *fiber.Ctx) error {
	packetId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	packetDetail, err := h.PacketService.GetPacketDetail(packetId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": packetDetail,
	})
}

func (h *PacketHandler) handleGetPacketById(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	packets, err := h.PacketService.GetALlPackets(int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": packets,
	})
}
