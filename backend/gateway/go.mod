module github.com/Egot3/supel/backend/gateway

go 1.26.1

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	github.com/rs/cors v1.11.1
)

require go.opentelemetry.io/otel v1.40.0 // indirect

require (
	github.com/Egot3/supel/backend/contracts v0.1.19
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)
