package utils

import (
	"server/internal/model"
	"sort"
)

func SortTierListByUpdatedAt(tierList *model.TierList) []string {
	characters := make([]string, 0, len(tierList.Character))

	for name := range tierList.Character {
		characters = append(characters, name)
	}

	sort.Slice(characters, func(i, j int) bool {
		countI := tierList.Character[characters[i]]
		countJ := tierList.Character[characters[j]]

		if countI == countJ {
			return characters[i] < characters[j]
		}

		return countI > countJ
	})

	return characters
}
