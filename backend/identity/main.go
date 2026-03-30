package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/server"
	"google.golang.org/grpc"
)

func main(){
	ctx := context.Background()
	if err := database.RunMigrations(ctx, database.DB); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}

	grpcServer := grpc.NewServer()

	productServer := server.NewIdentityServer()
	pb.RegisterIdentityServiceServer(grpcServer, productServer)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("identity Service gRPC server on %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}