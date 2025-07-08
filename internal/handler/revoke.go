package handler

import (
	"net/http"

	"github.com/karprabha/chirpy/internal/auth"
	"github.com/karprabha/chirpy/internal/config"
)

func Revoke(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rt, err := cfg.Queries.GetRefreshToken(r.Context(), refreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		err = cfg.Queries.RevokeRefreshToken(r.Context(), rt.Token)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
