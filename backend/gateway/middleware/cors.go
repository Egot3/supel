package middleware

import (
	"net/http"
	"os"

	"github.com/rs/cors"
)

func NewCORSMiddleware() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug:            os.Getenv("ENV") == "development",
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPut,
			http.MethodPost,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		MaxAge: 300,
	})
	return c.Handler
}
