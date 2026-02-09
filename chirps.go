package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ardamertdedeoglu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	profaneWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}
	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserId    string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	var res_body strings.Builder

	words := strings.SplitSeq(params.Body, " ")
	for word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			res_body.WriteString("****" + " ")
		} else {
			res_body.WriteString(word + " ")
		}
	}
	cleaned_body := strings.TrimSpace(res_body.String())
	user_id, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse uuid", err)
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: user_id,
	})

	respondWithJSON(w, http.StatusCreated, returnVals{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	})
}
