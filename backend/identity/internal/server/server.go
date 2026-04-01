package server

import (
	"context"
	"log"

	pb "github.com/Egot3/supel/backend/contracts"
	jwtutils "github.com/Egot3/supel/backend/identity/internal/JWTutils"
	"github.com/Egot3/supel/backend/identity/internal/database/repositories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
}

func NewIdentityServer() *IdentityServer {
	return &IdentityServer{} //duh
}

func (s *IdentityServer) GenerateToken(ctx context.Context, req *pb.TokenRequest) (*pb.Token, error) {
	log.Println("getting token", req)

	user, err := repositories.User(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found in db")
	}

	resp, err := jwtutils.GenerateToken(user.UUID, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.Token{Token: resp}, nil
}

func (s *IdentityServer) ValidateToken(ctx context.Context, req *pb.Token) (*pb.TokenPayload, error) {
	body, err := jwtutils.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	user, err := repositories.User(ctx, body.ID)

	return &pb.TokenPayload{
		Uuid: user.UUID,
		Role: string(user.Role),
	}, err
}
