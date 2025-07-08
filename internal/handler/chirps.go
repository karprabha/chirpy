package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body,omitempty"`
		Error       string `json:"error,omitempty"`
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
	respond(w, http.StatusOK, response{CleanedBody: clean})
}

func respond(w http.ResponseWriter, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
