package handlers

import (
	"jirbthagoras/raksana-backend/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type PacketHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
}

func NewPacketHandler(
	v *validator.Validate,
	r *redis.Client,
	rp *repositories.Queries,
) *PacketHandler {
	return &PacketHandler{
		Validator:  v,
		Repository: rp,
	}
}

func (h *PacketHandler) RegisterRoutes(router fiber.Router) {
	_ = router.Group("/packet")
}

func (h *PacketHandler) handleGeneratePacket(c *fiber.Ctx) error {
	return nil

}
