package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karprabha/chirpy/internal/auth"
	"github.com/karprabha/chirpy/internal/config"
	"github.com/karprabha/chirpy/internal/database"
)

func CreateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var p params
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if p.Email == "" || p.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(p.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createUserParams := database.CreateUserParams{
			Email:          p.Email,
			HashedPassword: hashedPassword,
		}

		user, err := cfg.Queries.CreateUser(r.Context(), createUserParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type response struct {
			ID          uuid.UUID `json:"id"`
			Email       string    `json:"email"`
			IsChirpyRed bool      `json:"is_chirpy_red"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		}

		resp := response{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func UpdateUser(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		type params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var p params
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Server error", http.StatusBadRequest)
			return
		}

		if p.Email == "" || p.Password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(p.Password)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		updateUserParams := database.UpdateUserParams{
			ID:             userID,
			Email:          p.Email,
			HashedPassword: hashedPassword,
		}

		user, err := cfg.Queries.UpdateUser(r.Context(), updateUserParams)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		type response struct {
			ID          uuid.UUID `json:"id"`
			Email       string    `json:"email"`
			IsChirpyRed bool      `json:"is_chirpy_red"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		}

		resp := response{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
