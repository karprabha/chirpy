package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karprabha/chirpy/internal/auth"
	"github.com/karprabha/chirpy/internal/config"
	"github.com/karprabha/chirpy/internal/database"
)

func Login(cfg *config.Config) http.HandlerFunc {
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

		user, err := cfg.Queries.GetUserByEmail(r.Context(), p.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		err = auth.CheckPasswordHash(p.Password, user.HashedPassword)
		if err != nil {
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		expiration := 1 * time.Hour

		token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expiration)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		createRefreshTokenParams := database.CreateRefreshTokenParams{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		}

		err = cfg.Queries.CreateRefreshToken(r.Context(), createRefreshTokenParams)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		type response struct {
			ID           uuid.UUID `json:"id"`
			Email        string    `json:"email"`
			IsChirpyRed  bool      `json:"is_chirpy_red"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
		}

		res := response{
			ID:           user.ID,
			Email:        user.Email,
			IsChirpyRed:  user.IsChirpyRed,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Token:        token,
			RefreshToken: refreshToken,
		}

		data, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
