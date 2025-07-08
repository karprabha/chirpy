package handler

import (
	"fmt"
	"net/http"

	"github.com/karprabha/chirpy/internal/config"
)

func AdminMetrics(cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<html>\n<head>\n<title>Chirpy Metrics</title>\n</head>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.FileServerHits.Load())
	})
}

func AdminReset(cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cfg.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		cfg.FileServerHits.Store(0)

		err := cfg.Queries.DeleteAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server reset"))
	})
}
