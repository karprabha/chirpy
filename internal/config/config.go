package config

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/karprabha/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type Config struct {
	DB             *sql.DB
	Queries        *database.Queries
	FileServerHits atomic.Int32
	Platform       string
	JWTSecret      string
}

func New() *Config {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	return &Config{
		DB:             db,
		Queries:        database.New(db),
		FileServerHits: atomic.Int32{},
		Platform:       platform,
		JWTSecret:      jwtSecret,
	}
}
