package moprconv

import (
	grpb "github.com/Egot3/supel/backend/contracts/group"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/egot3/supel/backend/group/internal/models"
)

func GroupMoToPr(mo *models.Group) *grpb.Group {
	grType, ok := grpb.GroupType_value[string(mo.GroupType)]
	if !ok {
		grType = int32(grpb.GroupType_UNSPECIFIED)
	}
	return &grpb.Group{
		UUID:        mo.UUID.String(),
		GroupType:   grpb.GroupType(grType),
		Name:        mo.Name,
		Description: mo.Description,
		CreatedAt:   timestamppb.New(mo.CreatedAt),
		UpdatedAt:   timestamppb.New(mo.UpdatedAt),
	}
}
