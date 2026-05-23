package types

import grpb "github.com/Egot3/supel/backend/contracts/group"

type GroupType string

const (
	UNSPECIFIED = "UNSPECIFIED"
	GROUP       = "GROUP"
	CLUB        = "CLUB"
)

func GroupTypeFromProto(gt grpb.GroupType) GroupType {
	switch gt {
	case grpb.GroupType_GROUP:
		return GROUP
	case grpb.GroupType_CLUB:
		return CLUB
	default:
		return UNSPECIFIED
	}
}
