package interceptor

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func HeaderInjector() runtime.ServeMuxOption {
	return runtime.WithMetadata(
		func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}

			if userUUID := r.Header.Get("user-uuid"); userUUID != "" {
				md.Set("user-uuid", userUUID)
			}
			if userRole := r.Header.Get("user-role"); userRole != "" {
				md.Set("user-role", userRole)
			}

			return md
		},
	)
}
