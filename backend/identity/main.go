package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	database.InitDB()
	ctx := context.Background()
	if err := database.RunMigrations(ctx, database.DB); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}
	var tableNames []string
	database.DB.NewSelect().
		TableExpr("information_schema.tables").
		Column("table_name").
		Where("table_schema = ?", "public").
		Scan(ctx, &tableNames)

	log.Println(tableNames)
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

	go func() {
		log.Printf("identity Service gRPC server on %s", port)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	gwMux := runtime.NewServeMux()

	selfConn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	authClient := pb.NewIdentityServiceClient(selfConn)

	httpMux := http.NewServeMux()
	httpMux.Handle("/v1/", gwMux)

	httpMux.Handle("/internal/identity/validate", server.NewForwardIdentityHandler(authClient))
	log.Println("HTTP server listening on :9030")
	if err := http.ListenAndServe(":9030", httpMux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
