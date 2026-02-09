package main

import (
	"net/http"

	"github.com/ardamertdedeoglu/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	ref_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token found", err)
	}

	err = cfg.db.RevokeToken(r.Context(), ref_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
