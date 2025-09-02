package handlers

import (
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/configs"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var allowedVideoTypes = []string{
	"video/mp4",
	"video/webm",
	"video/ogg",
}

type FileHandler struct {
	*configs.AWSClient
	Validator *validator.Validate
}

func NewFileHandler(
	v *validator.Validate,
	aws *configs.AWSClient,
) *FileHandler {
	return &FileHandler{
		Validator: v,
		AWSClient: aws,
	}
}

func (h *FileHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/file")
	g.Use(helpers.TokenMiddleware)
	g.Post("/presign", h.handleCreatePresigned)

}

func (h *FileHandler) handleCreatePresigned(c *fiber.Ctx) error {
	req := &models.PostFilePresigned{}
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

	option := c.Query("type")
	var key string

	id := uuid.New().String()
	ext := filepath.Ext(req.Filename)

	switch option {
	case "profile":
		key = fmt.Sprintf("profiles/%v/%s%s", userId, id, ext)
	case "memory":
		key = fmt.Sprintf("memories/%v/%s%s", userId, id, ext)
	case "scan":
		key = "scan/"
	default:
		return fiber.NewError(fiber.StatusBadRequest, "Unrecognized query param")
	}

	for _, t := range allowedVideoTypes {
		if strings.EqualFold(req.ContentType, t) {
			return fiber.NewError(fiber.StatusBadRequest, "Allowed content type: image/png, image/jpg, and videos/mp4")
		}
	}

	fileUrl, presignedReq, err := h.AWSClient.CreatePresignUrlPutObject(key, req.ContentType)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"presigned_url": presignedReq.URL,
			"file_url":      fileUrl,
		},
	})
}
