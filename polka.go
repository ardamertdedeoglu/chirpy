package main

import (
	"encoding/json"
	"net/http"

	"github.com/ardamertdedeoglu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "auth header not found", err)
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusForbidden, "invalid api key", err)
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
	}

	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse id", err)
	}

	err = cfg.db.UpgradeUser(r.Context(), user_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
