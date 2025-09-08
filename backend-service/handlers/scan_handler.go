package handlers

import (
	"context"
	"errors"
	"io"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ScanHandler struct {
	Validator *validator.Validate
	*TreasureHandler
	*QuestHandler
	*EventHandler
	*configs.AWSClient
}

func NewScanHandler(
	v *validator.Validate,
	th *TreasureHandler,
	qh *QuestHandler,
	eh *EventHandler,
	aws *configs.AWSClient,
) *ScanHandler {
	return &ScanHandler{
		Validator:       v,
		TreasureHandler: th,
		QuestHandler:    qh,
		EventHandler:    eh,
		AWSClient:       aws,
	}
}

func (h *ScanHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/scan")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleScan)
	g.Post("/trash", h.handleScanTrash)
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

func (h *ScanHandler) handleScanTrash(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		slog.Error("Failed to take image", "err", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		slog.Error("Failed to open the image", "err", err)
		return err
	}

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	ctx := context.Background()

	output, err := h.AWSClient.RekognitionClient.DetectLabels(ctx, &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: fileBytes,
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(75.0),
	})
	if err != nil {
		slog.Error("Something wrong with the image scanning", "err", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})

}
