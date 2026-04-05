package server

import (
	"context"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/news/internal/database/repositories"
	"github.com/Egot3/supel/backend/news/internal/models"
	storage "github.com/Egot3/supel/backend/news/internal/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NewsSever struct {
	pb.UnimplementedNewsServiceServer
	storageService storage.StorageService
}

func NewNewsService(storageService storage.StorageService) *NewsSever {
	return &NewsSever{
		storageService: storageService,
	}
}

func (s *NewsSever) CreateNew(ctx context.Context, req *pb.CreateNewRequest) (*pb.CreateNewResponse, error) {
	// there is an RBAC call, you just don't see it in MVP
	fileKeys := req.GetFileKeys()

	createdNew, err := repositories.CreateNew(ctx, models.New{
		Caption: req.Caption,
		Body:    req.GetBody(),
	}, fileKeys)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var imageGets []string
	for _, key := range fileKeys {
		imgGet, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageGets = append(imageGets, imgGet)
	}
	return &pb.CreateNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    createdNew.UserUUID,
			Caption:   createdNew.Caption,
			Body:      createdNew.Body,
			ImageUrls: imageGets,
			CreatedAt: timestamppb.New(createdNew.CreatedAt),
		},
	}, nil
}
