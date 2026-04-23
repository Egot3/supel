package server

import (
	"context"
	"log"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/user/internal/database/repositories"
	"github.com/Egot3/supel/backend/user/internal/models"
	"github.com/Egot3/supel/backend/user/internal/moprconv"
	storage "github.com/Egot3/supel/backend/user/internal/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserSever struct {
	pb.UnimplementedUserServiceServer
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

func NewUserService(storageService storage.StorageService) *UserSever {
	return &UserSever{
		storageService: storageService,
	}
}

func (s *UserSever) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	err := repositories.CreateUser(ctx, req.GetNickname(), req.GetUuid())

	if err != nil {
		log.Printf("Failed to create user %v: %v", req.GetUuid(), err.Error())
		return nil, status.Error(codes.Internal, "Failed to create pub user")
	}
	return nil, err
}

func (s *UserSever) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userUuid, role, ok := UserFromContext(ctx)
	if !ok {
		log.Printf("identity failure in deletion")
		return nil, status.Error(codes.DataLoss, "Identity failure")
	}
	targetUuid := req.GetUuid()

	if role != "ADMIN" && userUuid != req.GetUuid() {
		log.Printf("user %v doesn't own %v and not an admin(%v)", userUuid, targetUuid, role)
		return nil, status.Error(codes.PermissionDenied, "NOT ENOUGH *POWER*")
	}

	err := repositories.DeleteUser(ctx, targetUuid)
	if err != nil {
		log.Printf("couldn't delete user %v: %v", targetUuid, err.Error())
		return nil, status.Error(codes.Internal, "couldn't delete user")
	}

	return nil, nil
}

func (s *UserSever) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := repositories.GetUser(ctx, req.GetUuid())
	if err != nil {
		log.Printf("Failed to get user %v: %v", req.Uuid, err.Error())
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: moprconv.UserMoToPr(user),
	}, nil
}

func (s *UserSever) PatchUser(ctx context.Context, req *pb.PatchUserRequest) (*emptypb.Empty, error) {
	userUuid, role, ok := UserFromContext(ctx)
	if !ok {
		log.Printf("identity failure in deletion")
		return nil, status.Error(codes.DataLoss, "Identity failure")
	}
	targetUuid := req.GetUuid()

	if role != "ADMIN" && userUuid != req.GetUuid() {
		log.Printf("user %v doesn't own %v and not an admin(%v)", userUuid, targetUuid, role)
		return nil, status.Error(codes.PermissionDenied, "NOT. ENOUGH. POWER")
	}

	patched := &models.User{
		UUID: targetUuid,
	}
	if req.Nickname != nil {
		patched.Nickname = *req.Nickname
	}
	if req.Description != nil {
		patched.Description = req.Description
	}
	if req.Status != nil {
		patched.Status = req.Status
	}
	if req.StatusReactionKey != nil {
		patched.StatusReactionKey = req.StatusReactionKey
	}
	if req.StatusExpirationDate != nil {
		t := req.StatusExpirationDate.AsTime()
		patched.StatusExpiration = &t
	}

	err := repositories.PatchUser(ctx, patched)
	if err != nil {
		log.Printf("Couldn't patch user %v: %v", targetUuid, err.Error())
		return nil, status.Error(codes.Internal, "couldn't patch users")
	}

	return nil, nil
}
