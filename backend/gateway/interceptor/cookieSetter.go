package interceptor

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func CookieSetter() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(
		func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
			meta, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				return nil
			}

			cookies := meta.HeaderMD.Get("set-cookie")
			for _, cookie := range cookies {
				w.Header().Add("set-cookie", cookie)
			}

			return nil
		},
	)
}
