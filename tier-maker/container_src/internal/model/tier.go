package model

import (
	"server/internal/constant"
)

type TierList struct {
	Tier      constant.Tier  `json:"tier"`
	Character map[string]int `json:"character"`
}
