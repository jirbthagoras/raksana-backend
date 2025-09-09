package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
)

type ScanHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*TreasureHandler
	*QuestHandler
	*EventHandler
	*configs.AWSClient
	*configs.AIClient
}

func NewScanHandler(
	v *validator.Validate,
	r *repositories.Queries,
	th *TreasureHandler,
	qh *QuestHandler,
	eh *EventHandler,
	aws *configs.AWSClient,
	ai *configs.AIClient,
) *ScanHandler {
	return &ScanHandler{
		Validator:       v,
		Repository:      r,
		TreasureHandler: th,
		QuestHandler:    qh,
		EventHandler:    eh,
		AWSClient:       aws,
		AIClient:        ai,
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
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("Failed to get the subject", "err", err)
		return fiber.ErrUnauthorized
	}

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

	cnf := helpers.NewConfig()
	model, err := configs.InitModel(h.AIClient.Genai, cnf, configs.TrashScanner)
	if err != nil {
		slog.Error("Failed to init model", "err", err)
		return err
	}

	session := model.StartChat()
	session.History = []*genai.Content{}

	reqMsg, err := json.Marshal(output.Labels)
	if err != nil {
		slog.Error("Failed to parse the json to []bytes", "err", err)
		return err
	}

	resp, err := session.SendMessage(ctx, genai.Text(reqMsg))
	if err != nil {
		slog.Error("Failed to send the message", "err", err)
		return err
	}

	responseMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		responseMsg += fmt.Sprintf("%v\n", part)
	}

	var modelResponse models.AIResponseScan
	err = json.Unmarshal([]byte(responseMsg), &modelResponse)
	if err != nil {
		slog.Error("Failed to parse Gemini response content", "err", err)
		return err
	}

	scan, err := h.Repository.CreateScans(ctx, repositories.CreateScansParams{
		UserID:      int64(userId),
		Title:       modelResponse.Title,
		Description: modelResponse.Description,
	})
	if err != nil {
		slog.Error("Failed to insert scan", "err", err)
		return err
	}

	for _, i := range modelResponse.Items {
		_, err := h.Repository.CreateItems(ctx, repositories.CreateItemsParams{
			ScanID:      scan.ID,
			Name:        i.Name,
			Description: i.Description,
			Value:       i.Value,
		})
		if err != nil {
			slog.Error("Failed to insert items", "err", err)
			return err
		}

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": modelResponse,
	})
}
