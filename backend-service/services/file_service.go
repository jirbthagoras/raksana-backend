package services

import (
	"fmt"
	"jirbthagoras/raksana-backend/configs"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var allowedVideoTypes = []string{
	"video/mp4",
	"video/webm",
	"video/ogg",
}

type FileService struct {
	*configs.AWSClient
}

func NewFileService(
	aws *configs.AWSClient,
) *FileService {
	return &FileService{
		AWSClient: aws,
	}
}

func (h *FileService) CreatePresignedURL(option string, userId string, filename string, contentType string) (string, string, error) {
	var key string
	id := uuid.New().String()
	ext := filepath.Ext(filename)
	switch option {
	case "profile":
		key = fmt.Sprintf("profiles/%v/%s%s", userId, id, ext)
	case "memory":
		key = fmt.Sprintf("memories/%v/%s%s", userId, id, ext)
	case "scan":
		key = "scan/"
	default:
		return "", "", fiber.NewError(fiber.StatusBadRequest, "Unrecognized query param")
	}

	for _, t := range allowedVideoTypes {
		if strings.EqualFold(contentType, t) {
			return "", "", fiber.NewError(fiber.StatusBadRequest, "Allowed content type: image/png, image/jpg, and videos/mp4")
		}
	}

	presignedReq, err := h.AWSClient.CreatePresignUrlPutObject(key, contentType)
	if err != nil {
		return "", "", err
	}

	return presignedReq.URL, key, nil
}
