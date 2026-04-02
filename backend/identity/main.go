package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	database.InitDB()
	ctx := context.Background()
	if err := database.RunMigrations(ctx, database.DB); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}

	// if len(os.Getenv("RABBIT_HOST")) != 0 {
	// 	conn := amqp091.Connection{}
	// 	subscriber, err := sub.NewSubscriber(&conn)
	// 	if err != nil {
	// 		log.Fatalf("couldn't create channel of subscriber: %v", err)
	// 	}
	// 	sp := sub.SubscriberPackage{
	// 		Queue:     os.Getenv("RABBIT_QUEUE"),
	// 		Consumer:  os.Getenv("RABBIT_CONSUMER"),
	// 		AutoAck:   false,
	// 		Exclusive: true,
	// 		NoLocal:   false,
	// 		NoWait:    false,
	// 		Args:      nil,
	// 	}
	// 	subsFunc, err := subscriber.StartSubscriberFunc(sp)
	// 	if err != nil {
	// 		log.Fatalf("couldn't subscribe: %v", err)
	// 	}
	// 	ch := make(chan amqp091.Delivery, 1)
	// 	go subsFunc(ch)
	// 	go handlers.HandleSyncMessage(ch)
	// }

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("identity", grpc_health_v1.HealthCheckResponse_SERVING)

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
