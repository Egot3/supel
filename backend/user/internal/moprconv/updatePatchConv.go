package moprconv

import (
	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/user/internal/models"
)

func PatchPrToUpdateMo(proto *pb.PatchUserRequest) *models.UpdateUser {
	t := proto.StatusExpirationDate.AsTime()
	return &models.UpdateUser{
		UUID:              proto.Uuid,
		Nickname:          proto.Nickname,
		Description:       proto.Description,
		Status:            proto.Status,
		StatusExpiration:  &t,
		StatusReactionKey: proto.StatusReactionKey,
	}
}
