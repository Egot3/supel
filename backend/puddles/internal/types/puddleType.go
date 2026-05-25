package types

import ppb "github.com/Egot3/supel/backend/contracts/puddle"

type PuddleType string

const (
	UNSPECIFIED = "UNSPECIFIED"
	GROUP       = "GROUP"
	ONEONONE    = "ONEONONE"
	CHANNEL     = "CHANNEL"
)

func PuddleTypePrToTy(protoPuddle ppb.PuddleType) PuddleType {
	switch protoPuddle {
	case ppb.PuddleType_ONEONONE:
		return ONEONONE
	case ppb.PuddleType_CHANNEL:
		return CHANNEL
	case ppb.PuddleType_GROUP:
		return GROUP
	default:
		return UNSPECIFIED
	}
}
