package model

import (
	"server/internal/constant"
)

type TierRequest struct {
	Data      map[string][]string `json:"data"` // map[tier][array of char]
	Ramaining int                 `json:"remaining"`
}

type TierResponse struct {
	TotalCount   int64           `json:"count"`
	TierResponse []CharacterRate `json:"tier_response"`
}

type CharacterRate struct {
	Tier      constant.Tier     `json:"tier"`
	Character map[string]string `json:"character"`
}

func (t *TierResponse) Empty() {
	e := map[string]string{}
	for _, v := range constant.AllTier {
		CharacterRate := CharacterRate{
			Tier:      v,
			Character: e,
		}
		t.TierResponse = append(t.TierResponse, CharacterRate)
	}
}

type APIResponse struct {
	Message string `json:"message"`
}
