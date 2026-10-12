package route

import (
	"fmt"
	"net/http"
	"server/internal/constant"
	"server/internal/model"
	"server/internal/port"
	"server/internal/safemap"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

type TierMaker struct {
	port.TierMakerPort
	svc           port.TierMakerService
	safemap       *safemap.TierSafeMap
	globalCounter *atomic.Int64
}

func NewTierMakerRoute(svc port.TierMakerService, safemap *safemap.TierSafeMap, globalCounter *atomic.Int64) port.TierMakerPort {
	return &TierMaker{
		svc:           svc,
		safemap:       safemap,
		globalCounter: globalCounter,
	}
}

func (t *TierMaker) MakeTier(e *fiber.Ctx) error {
	t.globalCounter.Add(1)
	return e.Status(http.StatusOK).JSON(model.APIResponse{
		Message: "DONE",
	})
}

func (t *TierMaker) GetTiers(e *fiber.Ctx) error {
	resp := model.TierResponse{}
	resp.TotalCount = t.globalCounter.Load()

	// Keep track of the best tier for each character.
	type characterInfo struct {
		tier  constant.Tier
		rate  string
		value float64
	}

	best := make(map[string]characterInfo)

	tierPriority := make(map[constant.Tier]int)
	for i, tier := range constant.AllTier {
		tierPriority[tier] = len(constant.AllTier) - i
	}

	for _, key := range constant.AllTier {
		data, ok := t.safemap.SafeLoad(string(key))
		if !ok {
			continue
		}

		rates := t.convert(data.Character)

		for name, rate := range rates {
			value, err := strconv.ParseFloat(
				strings.TrimSuffix(rate, "%"),
				64,
			)
			if err != nil {
				continue
			}

			current, exists := best[name]
			if !exists ||
				value > current.value ||
				(value == current.value &&
					tierPriority[data.Tier] > tierPriority[current.tier]) {
				best[name] = characterInfo{
					tier:  data.Tier,
					rate:  rate,
					value: value,
				}
			}
		}
	}

	// Build the response, keeping one character per tier.
	tierCharacters := make(map[constant.Tier]map[string]string)

	for name, info := range best {
		if tierCharacters[info.tier] == nil {
			tierCharacters[info.tier] = make(map[string]string)
		}

		tierCharacters[info.tier][name] = info.rate
	}

	for _, tier := range constant.AllTier {
		characters := tierCharacters[tier]
		if len(characters) == 0 {
			continue
		}

		resp.TierResponse = append(resp.TierResponse, model.CharacterRate{
			Tier:      tier,
			Character: characters,
		})
	}

	if len(resp.TierResponse) == 0 {
		resp.Empty()
	}

	return e.Status(http.StatusOK).JSON(&resp)
}

func (t *TierMaker) convert(data map[string]int) map[string]string {
	total := t.globalCounter.Load()
	res := make(map[string]string, len(data))

	for k, v := range data {
		if total == 0 {
			res[k] = "0"
			continue
		}

		p := float64(v) / float64(total) * 100
		res[k] = fmt.Sprintf("%.2f%%", p)
	}

	return res
}
