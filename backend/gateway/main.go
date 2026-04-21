package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/gateway/gateways"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := gateways.NewGatewayMux()
	ctx := context.Background()

	identityConn, err := grpc.NewClient(fmt.Sprintf("%v:%v", os.Getenv("IDENTITY_SERVICE_HOST"), os.Getenv("IDENTITY_SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Identity conn fell: %v", err.Error())
	}
	newsConn, err := grpc.NewClient(fmt.Sprintf("%v:%v", os.Getenv("NEWS_SERVICE_HOST"), os.Getenv("NEWS_SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("News conn fell: %v", err.Error())
	}

	userConn, err := grpc.NewClient(fmt.Sprintf("%v:%v", os.Getenv("USER_SERVICE_HOST"), os.Getenv("USER_SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("User conn fell: %v", err.Error())
	}

	pb.RegisterIdentityServiceHandler(ctx, mux, identityConn)
	pb.RegisterNewsServiceHandler(ctx, mux, newsConn)
	pb.RegisterUserServiceHandler(ctx, mux, userConn)

	// corsMiddleware := middleware.NewCORSMiddleware()
	// handler := corsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":51000", mux))
}
