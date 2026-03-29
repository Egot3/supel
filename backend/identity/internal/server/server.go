package server

import (
	"context"

	pb "github.com/Egot3/supel/backend/contracts"
	jwtutils "github.com/Egot3/supel/backend/identity/internal/JWTutils"
	"github.com/Egot3/supel/backend/identity/internal/repositories"
	"github.com/Egot3/supel/backend/identity/types"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
}

func NewIdentityServer() *IdentityServer {
	return &IdentityServer{
	} //duh
}

func (s *IdentityServer) GetToken(ctx context.Context, req *pb.TokenPayload) (*pb.Token, error) {
	resp, err := jwtutils.GenerateToken(req.GetUuid(),types.UserRole(req.GetRole().String()))
	
	return &pb.Token{Body: resp}, err
}

func (s *IdentityServer) ValidateToken(ctx context.Context, req *pb.Token) (*pb.TokenPayload, error) {
	body, err := jwtutils.ValidateToken(req.GetBody())
	if err != nil {
		return nil, err
	}
	user, err := repositories.User(ctx, body.ID)

	

	return &pb.TokenPayload{
		Uuid: user.UUID,
		Role: pb.Role(pb.Role_value[string(body.Role)]),
	}, err
}