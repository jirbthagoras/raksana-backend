package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type PointHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	Mu         sync.Mutex
	*configs.AIClient
	*services.PointService
	*services.JournalService
}

func NewPointHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
) *PointHandler {
	return &PointHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
	}
}

func (h *PointHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/point")
	g.Use(helpers.TokenMiddleware)
	g.Get("/", h.handleGetCurrentBalance)
	g.Post("/", h.handleConvertPoint)
}

func (h *PointHandler) handleGetCurrentBalance(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx := context.Background()

	profile, err := h.Repository.GetUserProfile(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user profile", "err", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"balance": profile.Points,
		},
	})
}

func (h *PointHandler) handleConvertPoint(c *fiber.Ctx) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	req := &models.RequestPostConvert{}

	err = c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	ctx := context.Background()

	profile, err := h.Repository.GetUserProfile(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user profile", "err", err)
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	cnf := helpers.NewConfig()
	convertionRate := cnf.GetInt("CONVERTION_RATE")

	pointTotal := req.Amount * convertionRate
	if pointTotal > int(profile.Points) {
		return fiber.NewError(fiber.StatusBadRequest, "Saldo anda tidak cukup")
	}

	_, err = h.Repository.DecreaseUserPoints(ctx, repositories.DecreaseUserPointsParams{
		UserID: int64(userId),
		Points: int64(pointTotal),
	})

	region, err := h.Repository.GetRegionById(ctx, int64(req.RegionId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "region tidak ditemukan")
		}
		slog.Error("Failed to get region", "err", err)
		return err
	}

	err = h.Repository.IncreaseRegionTreeAmount(ctx, repositories.IncreaseRegionTreeAmountParams{
		TreeAmount: int32(req.Amount),
		ID:         region.ID,
	})
	if err != nil {
		slog.Error("Failed to update region", "err", err)
		return err
	}

	histMsg := fmt.Sprintf("Konversi poin ke pohon untuk region %s dalam jumlah %v pohon", region.Name, req.Amount)
	err = h.Repository.AppendHistry(ctx, repositories.AppendHistryParams{
		UserID:   int64(userId),
		Name:     histMsg,
		Category: "Convert",
		Type:     "output",
		Amount:   int32(req.Amount * convertionRate),
	})

	logMsg := fmt.Sprintf("Saya baru suaja menukar %v GP menjadi pohon dengan jumlah %v di region: %s", pointTotal, req.Amount, region.Name)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})

}
