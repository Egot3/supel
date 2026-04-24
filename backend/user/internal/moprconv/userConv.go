package moprconv

import (
	pb "github.com/Egot3/supel/backend/contracts"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Egot3/supel/backend/user/internal/models"
)

func UserMoToPr(model *models.User, avatarUrl *string) *pb.User {
	return &pb.User{
		Uuid:              model.UUID,
		Nickname:          model.Nickname,
		Description:       model.Description,
		AvatarUrl:         avatarUrl,
		Status:            model.Status,
		StatusReactionKey: model.StatusReactionKey,
		CreatedAt:         timestamppb.New(model.CreatedAt),
	}
}

func UserPrToMo(proto *pb.User) *models.User {
	return &models.User{
		UUID:              proto.Uuid,
		Nickname:          proto.Nickname,
		Description:       proto.Description,
		Status:            proto.Status,
		StatusReactionKey: proto.StatusReactionKey,
		CreatedAt:         proto.CreatedAt.AsTime(),
	}
}
