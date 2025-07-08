package handler

import (
	"fmt"
	"net/http"

	"github.com/karprabha/chirpy/internal/config"
)

func AdminMetrics(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html>\n<head>\n<title>Chirpy Metrics</title>\n</head>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.FileServerHits.Load())
	})
}

func AdminReset(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Store(0)
		w.Write([]byte("Hit counter reset"))
	})
}
