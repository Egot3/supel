package services

import (
	"fmt"
	"os"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRBACClient(i do.Injector) (rbacpb.RBACServiceClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", os.Getenv("RBAC_SERVICE_HOST"),
		os.Getenv("RBAC_SERVICE_PORT")), grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rbacClient := rbacpb.NewRBACServiceClient(conn)

	return rbacClient, nil
}
