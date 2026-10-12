package service

import (
	"server/internal/port"
	"server/internal/safemap"
)

type TierMaker struct {
	port.TierMakerService
}

func NewTierMakerService() port.TierMakerService {
	return &TierMaker{}
}

func (t *TierMaker) MakeTier(tier map[string][]string, safemap *safemap.TierSafeMap) error {
	return nil
}
