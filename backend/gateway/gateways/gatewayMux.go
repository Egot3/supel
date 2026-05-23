package gateways

import (
	"github.com/Egot3/supel/backend/gateway/interceptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewGatewayMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		interceptor.CookieSetter(),
		interceptor.HeaderInjector(),
		interceptor.RequestUUIDInjectorIn(),
		interceptor.RequestUUIDInjectorOut(),

		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			switch s {
			case "user-uuid":
				return s, true
			default:
				return runtime.DefaultHeaderMatcher(s)
			}
		}),

		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
		}),
	)
}
