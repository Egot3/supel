package interceptor

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func RequestUUIDInjectorIn() runtime.ServeMuxOption {
	return runtime.WithMetadata(
		func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}

			if uuid, err := uuid.NewV7(); err != nil {
				md.Set("request-uuid", uuid.String())
			}
			return md
		},
	)
}

func RequestUUIDInjectorOut() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(
		func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
			meta, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				return nil
			}

			reqUUIDs := meta.HeaderMD.Get("request-uuid")
			reqUUID := ""
			if len(reqUUIDs) != 0 {
				reqUUID = reqUUIDs[0]
			}

			w.Header().Add("request-uuid", reqUUID)

			return nil
		},
	)
}
