package server

import (
	"context"
	"errors"
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

func (s *IdentityServer) ValidateToken(ctx context.Context, req *pb.Token) (*pb.TokenPayload, error) {
	body, err := jwtutils.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Bad token %v", err)
		return nil, err
	}
	user, err := repositories.UserById(ctx, body.ID)

	return &pb.TokenPayload{
		Uuid: user.UUID,
		Role: string(user.Role),
	}, err
}

func (s *IdentityServer) RemintToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {
	token, err := jwtutils.RemintToken(req.GetToken())

	return &pb.Token{
		Token: token,
	}, err //мы тут не колеса изобретаем
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.Token, error) {

	log.Printf("email: %v, password: %v", req.Email, req.Password)
	uuid, role, err := repositories.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, errors.New("Invalid credetantials")) {
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

	return &pb.Token{
		Token: token,
	}, nil
}

func (s *IdentityServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Token, error) {
	uuid, role, err := repositories.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errors.New("User with this email alreay exists")) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else {
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	// conn, err := grpc.NewClient(fmt.Sprintf("%v:%v",
	// 	os.Getenv("USER_HOST"),
	// 	os.Getenv("USER_PORT")),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()))

	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "Couldn't register with name, try again later")
	// }
	// defer conn.Close()

	// client := pb.NewUserClient(conn)

	token, err := jwtutils.GenerateToken(uuid, role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Token{
		Token: token,
	}, nil
}
