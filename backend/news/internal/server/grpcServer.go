package server

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/news/internal/database/repositories"
	"github.com/Egot3/supel/backend/news/internal/models"
	"github.com/Egot3/supel/backend/news/internal/moprconv"
	storage "github.com/Egot3/supel/backend/news/internal/s3"
	sanitizationutils "github.com/Egot3/supel/backend/news/internal/sanitizationUtils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		Body:     req.BodyKey,
		BodySize: int64(req.BodySize),
	}, fileKeys)
	if err != nil {
		log.Printf("Couldn't create a new %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	imageUrls := make([]string, len(fileKeys))
	for _, key := range fileKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			log.Printf("couldn't create key for image: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageUrls = append(imageUrls, imageLink)
	}

	bodyUrl := new(string)
	bodyUrl = nil
	if createdNew.Body != nil {
		bUrl, err := s.storageService.GETurl(ctx, *createdNew.Body)
		if err != nil {
			log.Printf("Err while retrieving body: %v", err)
			return nil, status.Error(codes.Internal, "failed to create a GET url of body")
		}
		bodyUrl = &bUrl
	}

	return &pb.CreateNewResponse{
		New: &pb.New{
			NewId:     createdNew.NewUUID,
			UserId:    userUuid,
			Caption:   createdNew.Caption,
			BodyUrl:   bodyUrl,
			ImageUrls: imageUrls,
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

	imageUrls := make([]string, len(imageKeys))
	for _, key := range imageKeys {
		imageLink, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			log.Printf("couldn't create key for image: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageUrls = append(imageUrls, imageLink)
	}

	var bodyUrl string
	if createdNew.Body != nil {
		bodyUrl, err = s.storageService.GETurl(ctx, *createdNew.Body)
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
			ImageUrls: imageUrls,
			BodySize:  uint64(createdNew.BodySize),
			CreatedAt: timestamppb.New(createdNew.CreatedAt),
		},
	}, nil
}

func (s *NewsSever) GenerateNewUploadURLs(ctx context.Context, req *pb.GenerateUploadURLSsRequest) (*pb.GenerateNewUploadURLsResponse, error) {
	putUrls := make([]*pb.UploadTarget, 0, len(req.Images))
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

func (s *NewsSever) ListNews(ctx context.Context, req *pb.ListNewsRequest) (*pb.ListNewsResponse, error) {
	log.Println("Got request to list news")
	news, total, err := repositories.NewBulk(ctx, int(req.GetPage()), int(req.GetSize()))
	if err != nil {
		log.Printf("Error while listing news: %v", err)
		return nil, status.Error(codes.Internal, "Error while listing news")
	}
	log.Printf("news as models after fetching: %v", news)

	targetNews := make([]*pb.New, 0, len(news))
	for i, newEx := range news {
		bodyUrl := new(string)
		bodyUrl = nil
		if newEx.Body != nil {
			bUrl, err := s.storageService.GETurl(ctx, *newEx.Body)
			if err != nil {
				log.Printf("Err while retrieving body: %v", err)
				return nil, status.Error(codes.Internal, "failed to create a GET url of body")
			}
			bodyUrl = &bUrl
		}

		imageKeys, err := repositories.NewImagesByUUId(ctx, newEx.NewUUID)
		if err != nil {
			log.Printf("couldn't retriew image keys: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		imageUrls := make([]string, len(imageKeys))
		for _, key := range imageKeys {
			imageLink, err := s.storageService.GETurl(ctx, key)
			if err != nil {
				log.Printf("couldn't create key for image: %v", err)
				return nil, status.Error(codes.Internal, err.Error())
			}

			imageUrls = append(imageUrls, imageLink)
		}

		targetNews[i] = moprconv.NewConverter(&newEx, bodyUrl, imageUrls)
		log.Printf("TargetNew and new: %v \n %v", targetNews[i], newEx)
	}
	log.Printf("TargetNews after hydration: %v", targetNews)

	return &pb.ListNewsResponse{
		News:  targetNews,
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: uint64(total),
	}, nil
}

func (s *NewsSever) DeleteNew(ctx context.Context, req *pb.DeleteNewRequest) (*emptypb.Empty, error) {
	userUuid, role, ok := UserFromContext(ctx)
	if !ok {
		log.Println("bad creditantials: ", userUuid, role)
		return nil, status.Error(codes.Unauthenticated, "user id is not ok")
	}
	newUUID := req.NewId

	if role != "ADMIN" {
		is, err := repositories.IsCreator(ctx, userUuid, newUUID)
		if err != nil {
			log.Printf("couldn't check if user has access to deleting new: %v", err.Error())
			return nil, status.Error(codes.Internal, "failed to check ownership")
		}

		if !is {
			log.Printf("user %v doesn't own post %v and not an admin(%v)", userUuid, newUUID, role)
			return nil, status.Error(codes.PermissionDenied, "not enough POWER")
		}
	}

	err := repositories.DeleteNew(ctx, newUUID)
	if err != nil {
		log.Printf("couldn't delete new %v: %v", newUUID, err.Error())
		return nil, status.Error(codes.Internal, "Failed to delete the new")
	}

	return nil, nil
}
