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
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type ChallengeHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.MemoryService
	*services.PointService
	*services.JournalService
	*services.FileService
	*services.StreakService
}

func NewChallengeHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ms *services.MemoryService,
	ps *services.PointService,
	js *services.JournalService,
	fs *services.FileService,
	ss *services.StreakService,
) *ChallengeHandler {
	return &ChallengeHandler{
		Validator:      v,
		Repository:     r,
		MemoryService:  ms,
		PointService:   ps,
		JournalService: js,
		FileService:    fs,
		StreakService:  ss,
	}
}

func (h *ChallengeHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/challenge")
	g.Use(helpers.TokenMiddleware)
	g.Post("/", h.handleParticipate)
	g.Get("/today", h.handleGetTodayChallenge)
	g.Get("/", h.handleGetAllChallenges)
	g.Get("/:id", h.handleGetChallengeParticipants)
}

func (h *ChallengeHandler) handleParticipate(c *fiber.Ctx) error {
	req := &models.PostCreateParticipation{}

	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse body", "err", err)
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

	ctx := context.Background()

	challenge, err := h.Repository.GetChallengeWithDetail(ctx)
	if err != nil {
		slog.Error("Failed to get challenge with details", "err", err)
		return err
	}

	participation, err := h.Repository.CheckParticipation(ctx, repositories.CheckParticipationParams{
		UserID:      int64(userId),
		ChallengeID: challenge.ChallengeID,
	})
	if err != nil {
		slog.Error("Failed to get current timezone", "err", err)
		return err
	}

	if participation != 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Anda sudah berpartisipasi dalam tantangan ini")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("Failed to get current timezone", "err", err)
		return err
	}

	var todayDate string = time.Now().In(loc).Format("2006-01-02")

	if todayDate != challenge.CreatedAt.Time.Format("2006-01-02") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid challenge")
	}

	presignedUrl, fileKey, err := h.FileService.CreatePresignedURL(
		"memory",
		strconv.Itoa(userId),
		req.FileName,
		req.ContentType,
	)
	if err != nil {
		return err
	}

	memoryId, err := h.MemoryService.CreateMemory(req.Description, fileKey, userId)
	if err != nil {
		return err
	}

	_, err = h.Repository.CreateParticipation(context.Background(), repositories.CreateParticipationParams{
		MemoryID:    int64(memoryId),
		UserID:      int64(userId),
		ChallengeID: int64(challenge.ChallengeID),
	})
	if err != nil {
		slog.Error("Failed to insert row to participation", "err", err)
		return err
	}

	_, err = h.PointService.UpdateUserPoint(int64(userId), challenge.PointGain)
	if err != nil {
		return err
	}

	logMsg := fmt.Sprintf("Baru saja berpartisipasi dalam challenge harian day-%v, dan mendapatkan poin sebesar %v ", challenge.Day, challenge.PointGain)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)
	if err != nil {
		return err
	}

	_, err = h.Repository.IncreaseChallengesFieldByOne(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to increase challenges field by one", "err", err)
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"presigned_url": presignedUrl,
		},
	})
}

func (h *ChallengeHandler) handleGetTodayChallenge(c *fiber.Ctx) error {
	res, err := h.Repository.GetTodayChallenge(context.Background())
	if err != nil {
		slog.Error("Failed to get today challenge", "err", err)
		return err
	}

	var challenge = models.ResponseChallenge{
		ID:          int(res.ChallengeID),
		Name:        res.Name,
		Description: res.Description,
		Difficulty:  res.Difficulty,
		Day:         int(res.Day),
		PointGain:   int(res.PointGain),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": challenge,
	})
}

func (h *ChallengeHandler) handleGetAllChallenges(c *fiber.Ctx) error {
	res, err := h.Repository.GetAllChallenges(context.Background())
	if err != nil {
		slog.Error("Failed to get all challenges", "err", err)
		return err
	}

	var challenges []models.ResponseChallenge
	for _, challenge := range res {
		challenges = append(challenges, models.ResponseChallenge{
			ID:          int(challenge.ChallengeID),
			Name:        challenge.Name,
			Description: challenge.Description,
			Difficulty:  challenge.Difficulty,
			Day:         int(challenge.Day),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"challenges": challenges,
		},
	})
}

func (h *ChallengeHandler) handleGetChallengeParticipants(c *fiber.Ctx) error {
	challengeId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	_, err = h.Repository.GetChallengeWithDetailById(context.Background(), int64(challengeId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Challenge not found")
		}
		slog.Error("Failed to get challenge detail")
		return err
	}

	res, err := h.Repository.GetMemoriesByChallengeID(context.Background(), int64(challengeId))
	if err != nil {
		slog.Error("Failed to get memories relaetdd to challenge", "err", err)
		return err
	}

	var challengeMemories []models.ResponseMemory
	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")
	for _, memory := range res {
		challengeMemories = append(challengeMemories, models.ResponseMemory{
			UserID:      memory.UserID,
			UserName:    memory.Username,
			FileURL:     bucketUrl + memory.FileKey,
			Description: memory.Description,
			CreatedAt:   memory.MemoryCreatedAt.Time.Format("2006-01-02 15:00"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"participants": challengeMemories,
		},
	})
}
