package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/Egot3/supel/backend/contracts"
	carefulness "github.com/Egot3/supel/backend/identity/internal"
	jwtutils "github.com/Egot3/supel/backend/identity/internal/JWTutils"
	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/database/repositories"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	userClient pb.UserServiceClient
}

func NewIdentityServer(userClient pb.UserServiceClient) *IdentityServer {
	return &IdentityServer{userClient: userClient} //not so duh
}

func (s *IdentityServer) ValidateToken(ctx context.Context, req *pb.Token) (*pb.TokenPayload, error) {
	body, err := jwtutils.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Bad token %v", err)
		return nil, err
	}
	user, err := repositories.UserById(ctx, body.UserID)
	if err != nil || user == nil {
		log.Println("user not found with this token", err)
		return nil, status.Error(codes.Unauthenticated, "User not found")
	}

	return &pb.TokenPayload{
		Uuid: user.UUID,
		Role: string(user.Role),
	}, nil
}

func (s *IdentityServer) RemintToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {
	token, err := jwtutils.RemintToken(req.GetToken())

	return &pb.Token{
		Token: token,
	}, err //мы тут не колеса изобретаем
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*emptypb.Empty, error) {

	log.Printf("email: %v, password: %v", req.Email, req.Password)
	uuid, role, err := repositories.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, carefulness.InvalidCreditantials) {
			log.Printf("Bad creds %v", err)

			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			log.Printf("something is Bad %v", err)

			return nil, status.Error(codes.Internal, "Internal server erroro")
		}
	}

	token, err := jwtutils.GenerateToken(uuid, role)
	if err != nil {
		log.Printf("Error generating token %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", fmt.Sprintf("auth_token=%v; HttpOnly; Secure; SameSite=Lax; Path=/", token),
	)); err != nil {
		return nil, status.Error(codes.Internal, "grpc Cookie setting error")
	}

	return nil, nil
}

func (s *IdentityServer) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	//log.Printf("email: %v, password: %v", req.Email, req.Password)
	uuid, role, err := repositories.Register(ctx, req.Email, req.Password)

	if err != nil {
		if errors.Is(err, carefulness.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else {
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	_, err = s.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Nickname: req.Nickname,
		Uuid:     uuid,
	})
	if err != nil {
		log.Printf("Failed to create New user: %v", err.Error())
		return nil, status.Error(codes.Internal, "failed to create pub user")
	}

	token, err := jwtutils.GenerateToken(uuid, role)
	if err != nil {
		log.Printf("err: %v", err)
		return nil, status.Error(codes.Internal, "Token gen error")
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", fmt.Sprintf("auth_token=%v; HttpOnly; Secure; SameSite=Lax; Path=/", token),
	)); err != nil {
		return nil, status.Error(codes.Internal, "grpc Cookie setting error")
	}

	return nil, nil
}

func (s *IdentityServer) DisableUser(ctx context.Context, req *pb.DisableUserRequest) (*emptypb.Empty, error) {
	err := database.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := s.userClient.DeleteUser(ctx, &pb.DeleteUserRequest{
			Uuid: req.GetUuid(),
		})
		if err != nil {
			log.Printf("Failed to delete user at user %v", err.Error())
			return err
		}

		err = repositories.DisableUserTx(ctx, tx, req.GetUuid())
		if err != nil {
			log.Printf("Failed to disable user in identity %v", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, carefulness.UserNotFound) {
			return nil, status.Error(codes.NotFound, "User to delete was not found")
		}
		return nil, status.Error(codes.Internal, "Unable to delete user")
	}

	return nil, nil
}
