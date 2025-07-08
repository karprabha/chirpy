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

	mux.Handle("GET /admin/metrics", handler.AdminMetrics(appConfig))
	mux.Handle("POST /admin/reset", handler.AdminReset(appConfig))

	mux.Handle("GET /api/healthz", http.HandlerFunc(handler.Healthz))
	mux.Handle("POST /api/users", handler.CreateUser(appConfig))
	mux.Handle("POST /api/chirps", handler.CreateChirp(appConfig))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
