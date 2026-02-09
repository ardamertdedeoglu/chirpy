package main

import (
	"net/http"
	"time"

	"github.com/ardamertdedeoglu/chirpy/internal/auth"
)

type Token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	ref_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer token not found", err)
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), ref_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, 60*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create jwt token", err)
	}

	respondWithJSON(w, http.StatusOK, Token{
		Token: token,
	})
}
