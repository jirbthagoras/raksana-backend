package handlers

import (
	"errors"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ScanHandler struct {
	Validator *validator.Validate
	*TreasureHandler
	*QuestHandler
	*EventHandler
}

func NewScanHandler(
	v *validator.Validate,
	th *TreasureHandler,
	qh *QuestHandler,
	eh *EventHandler,
) *ScanHandler {
	return &ScanHandler{
		Validator:       v,
		TreasureHandler: th,
		QuestHandler:    qh,
		EventHandler:    eh,
	}
}

func (h *ScanHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/scan")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleScan)
}

func (h *ScanHandler) handleScan(c *fiber.Ctx) error {
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

	_, payload, err := helpers.ValidateActivityToken(req.Token)
	if err != nil {
		slog.Error("Failed to validate token")
		return fiber.NewError(fiber.StatusBadRequest, "Token Invalid")
	}

	switch payload.Type {
	case "treasure":
		return h.TreasureHandler.handleClaimTreasure(c)
	case "quest":
		return h.QuestHandler.handleContribute(c)
	case "event":
		return h.EventHandler.handleAttend(c)
	default:
		return fiber.ErrBadRequest
	}
}
