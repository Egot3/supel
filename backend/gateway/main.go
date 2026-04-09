package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/gateway/gateways"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := gateways.NewGatewayMux()
	ctx := context.Background()

	identityConn, err := grpc.NewClient("identity-service.identity-domain.svc.cluster.local:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Identity conn fell: %v", err.Error())
	}
	newsConn, err := grpc.NewClient("news-service.news-domain.svc.cluster.local:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("News conn fell: %v", err.Error())
	}

	pb.RegisterIdentityServiceHandler(ctx, mux, identityConn)
	pb.RegisterNewsServiceHandler(ctx, mux, newsConn)

	log.Fatal(http.ListenAndServe(":51000", mux))
}
