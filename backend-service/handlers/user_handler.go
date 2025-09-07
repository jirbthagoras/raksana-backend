package handlers

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Repository *repositories.Queries
	Validator  *validator.Validate
	*configs.AWSClient
	*services.UserService
	*services.LeaderboardService
	*services.FileService
	Mu sync.Mutex
}

func NewUserHandler(
	v *validator.Validate,
	r *repositories.Queries,
	us *services.UserService,
	ls *services.LeaderboardService,
	fs *services.FileService,
	a *configs.AWSClient,
) *UserHandler {
	return &UserHandler{
		Validator:          v,
		Repository:         r,
		UserService:        us,
		LeaderboardService: ls,
		FileService:        fs,
		AWSClient:          a,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	g1 := router.Group("/user")
	g1.Use(helpers.TokenMiddleware)
	g1.Get("/", h.handleGetAllUsers)

	g2 := router.Group("/profile")
	g2.Use(helpers.TokenMiddleware)
	g2.Get("/me", h.handleGetProfile)
	g2.Get("/:id", h.handleGetProfileById)
	g2.Put("/picture", h.handleUpdateProfilePicture)
}

func (h *UserHandler) handleGetProfile(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	profile, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": profile,
	})
}

func (h *UserHandler) handleGetProfileById(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to get packet id", "err", err)
	}

	profile, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": profile,
	})
}

func (h *UserHandler) handleUpdateProfilePicture(c *fiber.Ctx) error {
	req := &models.PutUserEditProfile{}

	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	ctx := context.Background()

	profile, err := h.Repository.GetUserProfile(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user prodile", "err", err)
		return err
	}

	cnf := helpers.NewConfig()
	bucketName := cnf.GetString("AWS_BUCKET")
	bucketUrl := cnf.GetString("AWS_URL")

	presignedUrl, key, err := h.FileService.CreatePresignedURL("profile", strconv.Itoa(userId), req.Filename, req.ContentType)
	if err != nil {
		return err
	}

	_ = h.AWSClient.DeleteObject(bucketName, profile.ProfileKey)

	err = h.Repository.UpdateUserProfile(ctx, repositories.UpdateUserProfileParams{
		UserID:     int64(userId),
		ProfileKey: key,
	})

	imageUrl := bucketUrl + key

	err = h.LeaderboardService.UpdateProfile(strconv.Itoa(userId), imageUrl)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"presigned_url": presignedUrl,
		},
	})
}

func (h *UserHandler) handleGetAllUsers(c *fiber.Ctx) error {
	ctx := context.Background()
	res, err := h.Repository.GetAllUser(ctx)
	if err != nil {
		slog.Error("Failed to get all users", "err", err)
		return err
	}

	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")

	var users []models.ResponseUser
	for _, user := range res {
		streak, err := h.StreakService.GetCurrentStreak(int(user.UserID))
		if err != nil {
			return err
		}
		users = append(users, models.ResponseUser{
			ID:         int(user.UserID),
			Level:      int(user.Level),
			Username:   user.Username,
			Email:      user.Email,
			ProfileUrl: bucketUrl + user.ProfileKey,
			Streak:     streak,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"users": users,
		},
	})
}
