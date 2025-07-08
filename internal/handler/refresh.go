package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/karprabha/chirpy/internal/auth"
	"github.com/karprabha/chirpy/internal/config"
)

func Refresh(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rt, err := cfg.Queries.GetRefreshToken(r.Context(), refreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if rt.ExpiresAt.Before(time.Now()) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		if rt.RevokedAt.Valid {
			http.Error(w, "refresh token has been revoked", http.StatusUnauthorized)
			return
		}

		expiration := 1 * time.Hour

		token, err := auth.MakeJWT(rt.UserID, cfg.JWTSecret, expiration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type response struct {
			Token string `json:"token"`
		}

		res := response{
			Token: token,
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
