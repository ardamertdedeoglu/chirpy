package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ardamertdedeoglu/chirpy/internal/auth"
	"github.com/ardamertdedeoglu/chirpy/internal/database"
	"github.com/google/uuid"
)

type returnVals struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}

func convertChirp(chirp database.Chirp) returnVals {
	return returnVals{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	profaneWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token", err)
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: user_id,
	})

	respondWithJSON(w, http.StatusCreated, convertChirp(chirp))
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps", err)
		return
	}

	returnChips := make([]returnVals, 0)
	for _, chirp := range chirps {
		returnChips = append(returnChips, convertChirp(chirp))
	}

	respondWithJSON(w, http.StatusOK, returnChips)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp_id_str := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(chirp_id_str)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing chirp ID", err)
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp with given id not found", err)
	}

	respondWithJSON(w, http.StatusOK, convertChirp(chirp))
}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token", err)
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
	}

	chirp_id_str := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(chirp_id_str)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing chirp ID", err)
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp with given id not found", err)
	}

	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "incorrect user", err)
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp_id,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
