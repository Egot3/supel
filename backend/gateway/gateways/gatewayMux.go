package gateways

import (
	"github.com/Egot3/supel/backend/gateway/interceptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func NewGatewayMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		interceptor.CookieSetter(),
		interceptor.HeaderIjector(),

		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			switch s {
			case "user-uuid", "user-role":
				return s, true
			default:
				return runtime.DefaultHeaderMatcher(s)
			}
		}),
	)
}
