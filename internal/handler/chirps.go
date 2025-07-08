package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/karprabha/chirpy/internal/config"
	"github.com/karprabha/chirpy/internal/database"
)

func respond(w http.ResponseWriter, code int, res any) {
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func CreateChirp(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}
		type response struct {
			ID        uuid.UUID `json:"id"`
			Body      string    `json:"body,omitempty"`
			UserID    uuid.UUID `json:"user_id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Error     string    `json:"error,omitempty"`
		}

		var params parameters
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if len(params.Body) > 140 {
			respond(w, http.StatusBadRequest, response{Error: "Chirp is too long"})
			return
		}

		profanity := map[string]bool{
			"kerfuffle": true, "sharbert": true, "fornax": true,
		}
		words := strings.Fields(params.Body)
		for i, w := range words {
			if profanity[strings.ToLower(w)] {
				words[i] = "****"
			}
		}
		clean := strings.Join(words, " ")

		createChirpParams := database.CreateChirpParams{
			Body:   clean,
			UserID: params.UserID,
		}

		chirp, err := cfg.Queries.CreateChirp(r.Context(), createChirpParams)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		respond(w, http.StatusCreated, response{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}
}

func GetChirps(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.Queries.GetChirps(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		type chirp struct {
			ID        uuid.UUID `json:"id"`
			Body      string    `json:"body,omitempty"`
			UserID    uuid.UUID `json:"user_id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Error     string    `json:"error,omitempty"`
		}

		response := make([]chirp, len(chirps))
		for i, c := range chirps {
			response[i] = chirp{
				ID:        c.ID,
				Body:      c.Body,
				UserID:    c.UserID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			}
		}

		respond(w, http.StatusOK, response)
	}
}
