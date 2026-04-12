package server

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/news/internal/database/repositories"
	"github.com/Egot3/supel/backend/news/internal/models"
	storage "github.com/Egot3/supel/backend/news/internal/s3"
	sanitizationutils "github.com/Egot3/supel/backend/news/internal/sanitizationUtils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NewsSever struct {
	pb.UnimplementedNewsServiceServer
	storageService storage.StorageService
}

func UserFromContext(ctx context.Context) (userID string, role string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", false
	}
	userID = md.Get("user-uuid")[0]
	role = md.Get("user-role")[0]
	return userID, role, !(len(userID) == 0 && len(role) == 0)
}

func NewNewsService(storageService storage.StorageService) *NewsSever {
	return &NewsSever{
		storageService: storageService,
	}
}

func (s *NewsSever) CreateNew(ctx context.Context, req *pb.CreateNewRequest) (*pb.CreateNewResponse, error) {
	// there is an RBAC call, you just don't see it in MVP
	fileKeys := req.GetImageKeys()
	log.Printf("fKeys: %v", fileKeys)

	userUuid, role, ok := UserFromContext(ctx)
	if !ok {
		log.Println("bad creditantials: ", userUuid, role)
		return nil, status.Error(codes.Unauthenticated, "user id is not ok")
	}

	createdNew, err := repositories.CreateNew(ctx, models.New{
		UserUUID: userUuid,
		Caption:  req.Caption,
		Body:     req.GetBodyKey(),
	}, fileKeys)
	if err != nil {
		log.Printf("Couldn't create a new %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var imageLinks []string
	for _, key := range fileKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			log.Printf("couldn't create key for image: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageLinks = append(imageLinks, imageLink)
	}

	var bodyUrl string
	if createdNew.Body != "" {
		bodyUrl, err = s.storageService.GETurl(ctx, createdNew.Body)
		if err != nil {
			log.Printf("Err while retrieving body: %v", err)
			return nil, status.Error(codes.Internal, "failed to create a GET url of body")
		}

	}

	return &pb.CreateNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    userUuid,
			Caption:   createdNew.Caption,
			BodyUrl:   &bodyUrl,
			ImageUrls: imageLinks,
			CreatedAt: timestamppb.New(createdNew.CreatedAt),
		},
	}, nil
}

func (s *NewsSever) GetNew(ctx context.Context, req *pb.GetNewRequest) (*pb.GetNewResponse, error) {
	createdNew, err := repositories.NewByUUID(ctx, req.GetNewId())
	if err != nil {
		log.Printf("Problem with finding a new: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	imageKeys, err := repositories.NewImagesByUUId(ctx, req.GetNewId())
	if err != nil {
		log.Printf("couldn't retriew image keys: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var imageLinks []string
	for _, key := range imageKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			log.Printf("couldn't create key for image: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageLinks = append(imageLinks, imageLink)
	}

	var bodyUrl string
	if createdNew.Body != "" {
		bodyUrl, err = s.storageService.GETurl(ctx, createdNew.Body)
		if err != nil {
			log.Printf("Err while retrieving body: %v", err)
			return nil, status.Error(codes.Internal, "failed to create a GET url of body")
		}

	}

	return &pb.GetNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    createdNew.UserUUID,
			Caption:   createdNew.Caption,
			BodyUrl:   &bodyUrl,
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

func (s *NewsSever) GenerateBodyUploadURL(ctx context.Context, req *pb.GenerateBodyUploadURLRequest) (*pb.GenerateBodyUploadURLResponse, error) {

	key := sanitizationutils.Slugify(
		fmt.Sprintf("orgs/ETSEvilCorp/news/body/%v/%v/%v",
			uuid.NewString(), time.Now().Format(time.RFC3339), req.BodyName))

	putUrl, err := s.storageService.PUTurl(ctx, key, "text/markdown")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GenerateBodyUploadURLResponse{
		Target: &pb.UploadTarget{
			UploadUrl: putUrl,
			FileKey:   key,
		},
	}, nil
}
