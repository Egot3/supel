package services

import (
	"fmt"
	"os"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewIdentityService(i do.Injector) (pb.IdentityServiceClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", os.Getenv("IDENTITY_SERVICE_HOST"),
		os.Getenv("IDENTITY_SERVICE_PORT")), grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	identityClient := pb.NewIdentityServiceClient(conn)

	return identityClient, nil
}
