package handlers

import (
	"context"
	"errors"
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
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.LeaderboardService
}

func NewAuthHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ls *services.LeaderboardService,
) *AuthHandler {
	return &AuthHandler{
		Validator:          v,
		Repository:         r,
		LeaderboardService: ls,
	}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/auth")
	g.Post("/register", h.handleRegister)
	g.Post("/login", h.handleLogin)
	g.Get("/me", h.handleMe)
	g.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "tested",
		})
	})
}

func (h *AuthHandler) handleRegister(c *fiber.Ctx) error {
	req := &models.PostUserRegister{}
	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "err", err.Error())
	}

	ctx := context.Background()

	user, err := h.Repository.CreateUser(ctx, repositories.CreateUserParams{
		Username: req.Username,
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
	})
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// return fiber.NewError(fiber.StatusBadRequest, err.Error())
			switch pgErr.ConstraintName {
			case "users_email_unique":
				return fiber.NewError(fiber.StatusBadRequest, "Email already used")
			case "users_username_unique":
				return fiber.NewError(fiber.StatusBadRequest, "Username already exists")
			}
		}

		slog.Error(err.Error())
		return err
	}

	profile, err := h.Repository.CreateProfile(ctx, repositories.CreateProfileParams{
		UserID:    user.ID,
		ExpNeeded: 50,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	userId := strconv.Itoa(int(user.ID))

	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")
	link := bucketUrl + profile.ProfileKey

	err = h.LeaderboardService.SetUserInfo(userId, user.Username, link)
	if err != nil {
		return err
	}

	err = h.LeaderboardService.UpdatePoint(userId, 0)
	if err != nil {
		return err
	}

	err = h.Repository.CreateStatistics(ctx, user.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"username": req.Username,
			"name":     req.Name,
			"email":    req.Email,
		},
	})
}

func (h *AuthHandler) handleLogin(c *fiber.Ctx) error {
	req := &models.PostUserLogin{}

	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	ctx := context.Background()

	user, err := h.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Email does not exists")
		}
		slog.Error("Failed to get user with such email", "err", err.Error())
		return err
	}

	if user.IsAdmin {
		return fiber.NewError(fiber.StatusUnauthorized, "Anda adalah admin")
	}

	if ok := helpers.CheckPassword(req.Password, user.Password); !ok {
		return fiber.NewError(fiber.StatusBadRequest, "Password does not match")
	}

	expiry := time.Now().Add(720 * time.Hour)
	token, err := helpers.GenerateToken(
		int(user.ID),
		user.Username,
		user.Email,
		expiry,
	)

	return c.Status(200).JSON(fiber.Map{
		"data": fiber.Map{
			"token": token,
		},
	})
}

func (h *AuthHandler) handleMe(c *fiber.Ctx) error {
	token, err := helpers.GetTokenFromRequest(c)
	if err != nil {
		slog.Error("Failed to get token from request", "err", err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Token is not attached")
	}

	_, claims, err := helpers.ValidateToken(token)
	if err != nil {
		slog.Error("Failed to validate token", "err", err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Token is not valid")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": claims,
	})
}
