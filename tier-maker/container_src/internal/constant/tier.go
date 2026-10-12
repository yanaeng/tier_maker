package constant

type Tier string

const (
	TierS Tier = "S"
	TierA Tier = "A"
	TierB Tier = "B"
	TierC Tier = "C"
	TierD Tier = "D"
)

var AllTier = []Tier{TierS, TierA, TierB, TierC, TierD}
