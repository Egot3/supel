package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/news/internal/database/repositories"
	"github.com/Egot3/supel/backend/news/internal/models"
	storage "github.com/Egot3/supel/backend/news/internal/s3"
	sanitizationutils "github.com/Egot3/supel/backend/news/internal/sanitizationUtils"
	"github.com/google/uuid"
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

	var imageLinks []string
	for _, key := range fileKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageLinks = append(imageLinks, imageLink)
	}

	return &pb.CreateNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    createdNew.UserUUID,
			Caption:   createdNew.Caption,
			Body:      createdNew.Body,
			ImageUrls: imageLinks,
			CreatedAt: timestamppb.New(createdNew.CreatedAt),
		},
	}, nil
}

func (s *NewsSever) GetNew(ctx context.Context, req *pb.GetNewRequest) (*pb.GetNewResponse, error) {
	createdNew, err := repositories.NewByUUID(ctx, req.GetNewId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	imageKeys, err := repositories.NewImagesByUUId(ctx, req.GetNewId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var imageLinks []string
	for _, key := range imageKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageLinks = append(imageLinks, imageLink)
	}

	return &pb.GetNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    createdNew.UserUUID,
			Caption:   createdNew.Caption,
			Body:      createdNew.Body,
			ImageUrls: imageLinks,
			CreatedAt: timestamppb.New(createdNew.CreatedAt),
		},
	}, nil
}

func (s *NewsSever) GenerateNewUploadURL(ctx context.Context, req *pb.GenerateUploadURLSsRequest) (*pb.GenerateNewUploadURLsResponse, error) {
	putUrls := make([]*pb.UploadTarget, len(req.Images))
	for _, meta := range req.Images {
		key := sanitizationutils.Slugify(
			fmt.Sprintf("orgs/ETSEvilCorp/news/attachments/%v/%v/%v",
				uuid.NewString(), time.Now().Format(time.RFC3339), meta.FileName))

		putUrl, err := s.storageService.PUTurl(ctx, key, meta.Mime)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		putUrls = append(putUrls, &pb.UploadTarget{
			UploadUrl: putUrl,
			FileKey:   key,
		})
	}

	return &pb.GenerateNewUploadURLsResponse{
		Targets: putUrls,
	}, nil
}
