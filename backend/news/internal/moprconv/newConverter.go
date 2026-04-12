package moprconv

import (
	pb "github.com/Egot3/supel/backend/contracts"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Egot3/supel/backend/news/internal/models"
)

func NewConverter(newModel *models.New, bodyUrl *string, imageUrls []string) (newProto *pb.New) {
	newProto = &pb.New{
		NewId:     newModel.NewUUID,
		UserId:    newModel.UserUUID,
		Caption:   newModel.Caption,
		BodyUrl:   bodyUrl,
		ImageUrls: imageUrls,
		CreatedAt: timestamppb.New(newModel.CreatedAt),
	}

	return newProto
}
