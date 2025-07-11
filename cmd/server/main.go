package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/karprabha/chirpy/internal/config"
	"github.com/karprabha/chirpy/internal/handler"
	"github.com/karprabha/chirpy/internal/middleware"
)

func getRootPath() string {
	rootPath, err := filepath.Abs(filepath.Join("."))
	if err != nil {
		log.Fatalf("Unable to resolve root path: %v", err)
	}

	return rootPath
}

func main() {
	appConfig := config.New()
	defer appConfig.DB.Close()

	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(getRootPath()))
	mux.Handle("/app/", middleware.WithMetrics(appConfig, http.StripPrefix("/app/", fileHandler)))

	// Admin routes
	mux.Handle("GET /admin/metrics", handler.AdminMetrics(appConfig))
	mux.Handle("POST /admin/reset", handler.AdminReset(appConfig))

	// Health route
	mux.Handle("GET /api/healthz", http.HandlerFunc(handler.Healthz))

	// Auth routes
	mux.Handle("POST /api/login", handler.Login(appConfig))
	mux.Handle("POST /api/revoke", handler.Revoke(appConfig))
	mux.Handle("POST /api/refresh", handler.Refresh(appConfig))

	// User routes
	mux.Handle("PUT /api/users", handler.UpdateUser(appConfig))
	mux.Handle("POST /api/users", handler.CreateUser(appConfig))

	// Chirp routes
	mux.Handle("POST /api/chirps", handler.CreateChirp(appConfig))
	mux.Handle("GET /api/chirps", handler.GetChirps(appConfig))
	mux.Handle("GET /api/chirps/{id}", handler.GetChirp(appConfig))
	mux.Handle("DELETE /api/chirps/{id}", handler.DeleteChirp(appConfig))

	// Webhooks routes
	mux.Handle("POST /api/polka/webhooks", handler.PolkaWebhook(appConfig))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
