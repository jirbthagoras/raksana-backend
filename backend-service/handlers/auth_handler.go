package handlers

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
}

func NewAuthHandler(
	v *validator.Validate,
	r *repositories.Queries,
) *AuthHandler {
	return &AuthHandler{
		Validator:  v,
		Repository: r,
	}
}

func (h *AuthHandler) RegisterRoute(router fiber.Router) {
	g := router.Group("/auth")
	g.Post("/register", h.handlRegister)
	g.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "tested",
		})
	})
}

func (h *AuthHandler) handlRegister(c *fiber.Ctx) error {
	req := &models.PostUserRegister{}
	err := h.Validator.Struct(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return helpers.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "err", err.Error())
	}

	_, err = h.Repository.CreateUser(context.Background(), repositories.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	})
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fiber.NewError(fiber.StatusBadRequest, "Username already exists")
		}

		slog.Error(err.Error())
		return err
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Successfully registered user",
	})
}
