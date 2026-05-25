package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/user/internal/database/repositories"

	"github.com/Egot3/supel/backend/user/internal/moprconv"
	storage "github.com/Egot3/supel/backend/user/internal/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	storageService storage.StorageService
}

func UserFromContext(ctx context.Context) (userID string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	userID = md.Get("user-uuid")[0]
	return userID, !(len(userID) == 0)
}

func NewUserService(storageService storage.StorageService) *UserServer {
	return &UserServer{
		storageService: storageService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	UUID := req.GetUuid()
	key := fmt.Sprintf("orgs/ETSEvilCorp/user/avatar/%v", UUID)

	err := repositories.CreateUser(ctx, req.GetNickname(), UUID, key)

	if err != nil {
		log.Printf("Failed to create user %v: %v", UUID, err.Error())
		return nil, status.Error(codes.Internal, "Failed to create pub user")
	}
	return nil, err
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	targetUuid := req.GetUuid()

	/* 	if role != "ADMIN" && userUuid != req.GetUuid() {
		log.Printf("user %v doesn't own %v and not an admin(%v)", userUuid, targetUuid, role)
		return nil, status.Error(codes.PermissionDenied, "NOT ENOUGH *POWER*")
	} */

	err := repositories.DeleteUser(ctx, targetUuid)
	if err != nil {
		log.Printf("couldn't delete user %v: %v", targetUuid, err.Error())
		return nil, status.Error(codes.Internal, "couldn't delete user")
	}

	return nil, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := repositories.GetUser(ctx, req.GetUuid())
	if err != nil {
		log.Printf("Failed to get user %v: %v", req.Uuid, err.Error())
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	avatarUrl := ""
	if aUrl, err := s.storageService.GETurl(ctx, user.AvatarKey); err == nil {
		avatarUrl = aUrl
	}

	return &pb.GetUserResponse{
		User: moprconv.UserMoToPr(user, avatarUrl),
	}, nil
}

func (s *UserServer) PatchUser(ctx context.Context, req *pb.PatchUserRequest) (*emptypb.Empty, error) {
	userUuid, ok := UserFromContext(ctx)
	if !ok {
		log.Printf("identity failure in deletion")
		return nil, status.Error(codes.DataLoss, "Identity failure")
	}
	targetUuid := req.GetUuid()

	//RBAC
	if userUuid != req.GetUuid() {
		log.Printf("user %v doesn't own %v", userUuid, targetUuid)
		return nil, status.Error(codes.PermissionDenied, "NOT. ENOUGH. POWER")
	}

	err := repositories.PatchUser(ctx, moprconv.PatchPrToUpdateMo(req))
	if err != nil {
		log.Printf("Couldn't patch user %v: %v", targetUuid, err.Error())
		return nil, status.Error(codes.Internal, "couldn't patch users")
	}

	return nil, nil
}

func (s *UserServer) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarResponse, error) {
	userUuid, ok := UserFromContext(ctx)
	if !ok {
		log.Printf("identity failure in avatar changing")
		return nil, status.Error(codes.DataLoss, "Identity failure")
	}
	targetUuid := req.GetUuid()

	//RBAC
	if userUuid != req.GetUuid() {
		log.Printf("user %v doesn't own %v and not an admin()", userUuid, targetUuid)
		return nil, status.Error(codes.PermissionDenied, "NOT. ENOUGH. POWER")
	}

	key, err := repositories.GetAvatarKey(ctx, req.Uuid)
	if err != nil || key == "" {
		log.Printf("couldn't retrieve a key for %v: %v", req.Uuid, err.Error())
		return nil, status.Error(codes.Internal, "Key loss")
	}

	avatarPUTUrl, err := s.storageService.PUTurl(ctx, key, "image/webp")
	if err != nil {
		log.Printf("Couldn't get a PUT url for %v: %v", req.Uuid, err.Error())
		return nil, status.Error(codes.Internal, "Couldn't generate a PUT url")
	}
	return &pb.UploadAvatarResponse{
		AvatarUrl: avatarPUTUrl,
	}, nil
}

func (s *UserServer) GetSelf(ctx context.Context, _ *emptypb.Empty) (*pb.GetSelfResponse, error) {
	userUuid, ok := UserFromContext(ctx)
	if !ok {
		log.Printf("identity failure in getting self")
		return nil, status.Error(codes.DataLoss, "Identity failure")
	}

	userModel, err := repositories.GetUser(ctx, userUuid)
	if err != nil {
		log.Printf("Couldn't fetch self-user %v: %v", userUuid, err.Error())
		return nil, status.Error(codes.Internal, "Couldn't fetch self")
	}

	avatarUrl := ""
	if aUrl, err := s.storageService.GETurl(ctx, userModel.AvatarKey); err == nil {
		avatarUrl = aUrl
	}

	return &pb.GetSelfResponse{
		User: moprconv.UserMoToPr(userModel, avatarUrl),
	}, nil
}
