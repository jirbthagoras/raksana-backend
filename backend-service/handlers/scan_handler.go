package handlers

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5"
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
	g.Get("/trash", h.handleGetAllScans)
	g.Post("/greenprint/:id", h.handleGenerateGreenprint)
	g.Get("/greenprint/:id", h.handleGetGreenprint)
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
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	ctx := context.Background()

	cnf := helpers.NewConfig()
	bucketName := cnf.GetString("AWS_BUCKET")
	key := fmt.Sprintf("scans/%v/%s", userId, file.Filename)

	_, err = h.AWSClient.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		slog.Error("Failed to upload to S3", "err", err)
		return err
	}

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
		ImageKey:    key,
	})
	if err != nil {
		slog.Error("Failed to insert scan", "err", err)
		return err
	}

	for _, i := range modelResponse.Items {
		_, err := h.Repository.CreateItems(ctx, repositories.CreateItemsParams{
			ScanID:      scan.ID,
			UserID:      int64(userId),
			Name:        i.Name,
			Description: i.Description,
			Value:       i.Value,
		})
		if err != nil {
			slog.Error("Failed to insert items", "err", err)
			return err
		}
	}

	bucketUrl := cnf.GetString("AWS_URL")

	modelResponse.ImageKey = bucketUrl + scan.ImageKey

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": modelResponse,
	})
}

func (h *ScanHandler) handleGetAllScans(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("Failed to get subject from token", "err", err)
		return nil
	}

	ctx := context.Background()

	scanRes, err := h.Repository.GetAllUserScans(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get scan result", "err", err)
		return err
	}

	var scans []models.AIResponseScan
	for _, s := range scanRes {
		resItem, err := h.Repository.GetItemsByScanId(ctx, s.ID)
		if err != nil {
			slog.Error("Failed to get item", "err", err)
			return err
		}

		var items []models.ResponseItems
		for _, i := range resItem {
			var isHavingGreenPrint bool = true

			_, err := h.Repository.GetGreenprints(ctx, i.ID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					isHavingGreenPrint = false
				} else {
					slog.Error("Failed to get greenprint", "err", err)
					return err
				}
			}

			items = append(items, models.ResponseItems{
				Id:               int(i.ID),
				Name:             i.Name,
				Description:      i.Description,
				Value:            i.Value,
				HavingGreenprint: isHavingGreenPrint,
			})
		}

		cnf := helpers.NewConfig()
		bucketUrl := cnf.GetString("AWS_URL")

		scans = append(scans, models.AIResponseScan{
			Title:       s.Title,
			Description: s.Description,
			ImageKey:    bucketUrl + s.ImageKey,
			Items:       items,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"scans": scans,
		},
	})
}

func (h *ScanHandler) handleGenerateGreenprint(c *fiber.Ctx) error {
	itemId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
		return err
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		slog.Error("Failed to get the subject", "err", err)
		return fiber.ErrUnauthorized
	}

	ctx := context.Background()

	resItem, err := h.Repository.GetItemsById(ctx, int64(itemId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Item not found")
		}
		slog.Error("Failed to get item by id")
		return err
	}

	if userId != int(resItem.UserID) {
		return fiber.NewError(fiber.StatusBadRequest, "Item not found")
	}

	cnf := helpers.NewConfig()

	// Generate greeprint
	model, err := configs.InitModel(h.AIClient.Genai, cnf, configs.GreenPrint)
	if err != nil {
		return err
	}

	session := model.StartChat()
	session.History = []*genai.Content{}

	aiMsg := fmt.Sprintf("Saya hendak membuat sebuah tutorial atau langkah-langkah terperinci untuk membuat: %s, dengan deskripsi sebagai berikut: %s", resItem.Name, resItem.Description)
	resp, err := session.SendMessage(ctx, genai.Text(aiMsg))
	if err != nil {
		slog.Error("Failed to send message", "err", err)
	}

	responseMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		responseMsg += fmt.Sprintf("%v\n", part)
	}

	var greenprintRes models.AIResponseGreenprint
	err = json.Unmarshal([]byte(responseMsg), &greenprintRes)
	if err != nil {
		slog.Error("Failed to parse Gemini response content", "err", err)
		return err
	}

	gp, err := h.Repository.CreateGreenprint(ctx, repositories.CreateGreenprintParams{
		ItemID:              int64(itemId),
		ImageKey:            "anjas kelas",
		Description:         greenprintRes.Description,
		Title:               greenprintRes.Title,
		SustainabilityScore: greenprintRes.SustainabilityScore,
	})
	if err != nil {
		slog.Error("Failed to create greenprint", "err", err)
		return err
	}

	for _, s := range greenprintRes.Steps {
		_, err := h.Repository.CreateSteps(ctx, repositories.CreateStepsParams{
			GreenprintID: gp.ID,
			Description:  s.Description,
		})
		if err != nil {
			slog.Error("Failed to create greenprint", "err", err)
			return err
		}
	}

	for _, m := range greenprintRes.Materials {
		_, err = h.Repository.CreateMaterials(ctx, repositories.CreateMaterialsParams{
			GreenprintID: gp.ID,
			Name:         m.Name,
			Description:  m.Description,
			Price:        m.Price,
			Quantity:     m.Quantity,
		})
	}

	for _, t := range greenprintRes.Tools {
		_, err = h.Repository.CreateTools(ctx, repositories.CreateToolsParams{
			GreenprintID: gp.ID,
			Name:         t.Name,
			Description:  t.Description,
			Price:        t.Price,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": greenprintRes,
	})
}

func (h *ScanHandler) handleGetGreenprint(c *fiber.Ctx) error {
	itemId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	ctx := context.Background()

	greenprintRes, err := h.Repository.GetGreenprints(ctx, int64(itemId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Greenprint not found")
		}
		slog.Error("Failed to get greenprint", "err", err)
		return err
	}

	materialsRes, err := h.Repository.GetMaterials(ctx, greenprintRes.ID)
	if err != nil {
		slog.Error("Failed to get materials", "err", err)
		return err
	}

	toolsRes, err := h.Repository.GetTools(ctx, greenprintRes.ID)
	if err != nil {
		slog.Error("Failed to get tools", "err", err)
		return err
	}

	stepsRes, err := h.Repository.GetSteps(ctx, greenprintRes.ID)
	if err != nil {
		slog.Error("Failed to get tools", "err", err)
		return err
	}

	var res = models.AIResponseGreenprint{
		Title:               greenprintRes.Title,
		Description:         greenprintRes.Description,
		EstimatedTime:       greenprintRes.EstimatedTime,
		SustainabilityScore: greenprintRes.SustainabilityScore,
		CreatedAt:           greenprintRes.CreatedAt.Time.Format("2006-01-02 15:04"),
		Tools:               []models.ResponseTool{},
		Materials:           []models.ResponseMaterial{},
		Steps:               []models.ResponseStep{},
		Text:                "",
	}

	for _, m := range materialsRes {
		res.Materials = append(res.Materials, models.ResponseMaterial{
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			Quantity:    m.Quantity,
		})
	}

	for _, t := range toolsRes {
		res.Tools = append(res.Tools, models.ResponseTool{
			Name:        t.Name,
			Description: t.Description,
			Price:       t.Price,
		})
	}

	for _, s := range stepsRes {
		res.Steps = append(res.Steps, models.ResponseStep{
			Description: s.Description,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})

}
