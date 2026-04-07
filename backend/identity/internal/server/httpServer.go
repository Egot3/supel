package server

import (
	"fmt"
	"net/http"
	"strings"

	pb "github.com/Egot3/supel/backend/contracts"
)

func NewForwardIdentityHandler(authClient pb.IdentityServiceClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing auth header", 401)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Strange auth token", 401)
			return
		}

		resp, err := authClient.ValidateToken(r.Context(), &pb.Token{
			Token: token,
		})
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		remintedToken, err := authClient.RemintToken(r.Context(), &pb.Token{
			Token: token,
		})
		cookie := &http.Cookie{
			Name:     "Authorization",
			Value:    fmt.Sprintf("Bearer %v", remintedToken),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)

		w.Header().Set("user-uuid", resp.Uuid)
		w.Header().Set("user-role", resp.Role)
		w.WriteHeader(200)
	}
}
