package port

import (
	"server/internal/safemap"

	"github.com/gofiber/fiber/v2"
)

type TierMakerPort interface {
	MakeTier(e *fiber.Ctx) error
	GetTiers(e *fiber.Ctx) error
}
type TierMakerService interface {
	MakeTier(tier map[string][]string, safemap *safemap.TierSafeMap) error
}
