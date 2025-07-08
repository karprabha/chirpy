package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karprabha/chirpy/internal/auth"
	"github.com/karprabha/chirpy/internal/config"
	"github.com/karprabha/chirpy/internal/database"
)

func PolkaWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		if apiKey != cfg.PolkaKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		type params struct {
			Event string `json:"event"`
			Data  struct {
				UserID uuid.UUID `json:"user_id"`
			} `json:"data"`
		}

		var p params
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
			return
		}

		if p.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		user, err := cfg.Queries.GetUserByID(r.Context(), p.Data.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		updateUserIsChirpyRedParams := database.UpdateUserIsChirpyRedParams{
			ID:          user.ID,
			IsChirpyRed: true,
		}

		user, err = cfg.Queries.UpdateUserIsChirpyRed(r.Context(), updateUserIsChirpyRedParams)
		if err != nil {
			http.Error(w, "Failed to upgrade user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
