package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type TreasureHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.PointService
	*services.JournalService
	*services.StreakService
}

func NewTreasureHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
	ss *services.StreakService,
) *TreasureHandler {
	return &TreasureHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
		StreakService:  ss,
	}
}

func (h *TreasureHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/treasure")
	g.Use(helpers.TokenMiddleware)
	g.Get("/me", h.handlGetCurrentUserTreasures)
	g.Get("/:id", h.handleGetUserTreasure)
}

func (h *TreasureHandler) handleClaimTreasure(c *fiber.Ctx) error {
	req := &models.ActivityRequest{}

	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err)
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

	_, payload, err := helpers.ValidateActivityToken(req.Token)
	if err != nil {
		slog.Error("Failed to validate token")
		return fiber.NewError(fiber.StatusBadRequest, "Token Invalid")
	}

	ctx := context.Background()

	treasure, err := h.Repository.GetTreasureByCodeId(ctx, payload.Subject)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Treasure tidak ditemukan")
		}
		slog.Error("Failed to get treasure from db", "err", err)
		return err
	}

	err = h.Repository.CreateClaimed(ctx, repositories.CreateClaimedParams{
		UserID:     int64(userId),
		TreasureID: treasure.ID,
	})
	if err != nil {
		slog.Error("Failed to insert into db", "err", err)
		return err
	}

	err = h.Repository.DeactivateTreasure(ctx, treasure.ID)
	if err != nil {
		slog.Error("Failed to update the row", "err", err)
		return err
	}

	profile, err := h.Repository.GetUserProfile(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user profile", "err", err)
		return err
	}

	historyMsg := fmt.Sprintf("Mendapatkan poin treasure: %s", treasure.Name)
	_, err = h.PointService.UpdateUserPoint(int64(userId), treasure.PointGain, historyMsg, "treasure", int(profile.Level))
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	logMsg := fmt.Sprintf("Baru saja mendapatkan treasure: '%s', memperoleh poin: %v", treasure.Name, treasure.PointGain)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)
	if err != nil {
		return nil
	}

	_, err = h.Repository.IncreaseTreasuresFieldByOne(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to update quest field", "err", err)
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"type":    "treasures",
			"message": "success",
			"treasure": fiber.Map{
				"name":       treasure.Name,
				"point_gain": treasure.PointGain,
			},
		},
	})
}

func (h *TreasureHandler) handlGetCurrentUserTreasures(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetAllClaimedTreasure(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get claimed treasure", "err", err)
		return err
	}

	var treasures []models.ResponseTreasure
	for _, treasure := range res {
		treasures = append(treasures, models.ResponseTreasure{
			Name:      treasure.Name,
			PointGain: int(treasure.PointGain),
			ClaimedAt: treasure.ClaimedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": treasures,
	})
}

func (h *TreasureHandler) handleGetUserTreasure(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	res, err := h.Repository.GetAllClaimedTreasure(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed to get claimed treasure", "err", err)
		return err
	}

	var treasures []models.ResponseTreasure
	for _, treasure := range res {
		treasures = append(treasures, models.ResponseTreasure{
			Name:      treasure.Name,
			PointGain: int(treasure.PointGain),
			ClaimedAt: treasure.ClaimedAt.Time.Format("2006-01-02 15:04"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"treasures": treasures,
		},
	})
}
