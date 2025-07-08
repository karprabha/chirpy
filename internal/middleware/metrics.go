package middleware

import (
	"net/http"

	"github.com/karprabha/chirpy/internal/config"
)

func WithMetrics(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
