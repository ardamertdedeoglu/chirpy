package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ardamertdedeoglu/chirpy/internal/auth"
	"github.com/ardamertdedeoglu/chirpy/internal/database"
)

type parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserVals struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

func convertUser(user database.User) UserVals {
	return UserVals{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	}
}

func (cfg *apiConfig) handleCreateUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, convertUser(user))
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Email    string
		Password string
	}

	type UserCredentials struct {
		Id           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		fmt.Println("user not able")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't check hash", err)
	}

	expire := time.Duration(3600 * time.Second)

	token, err := auth.MakeJWT(user.ID, cfg.secret, expire)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create jwt token", err)
	}

	ref_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: ref_token})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error storing refresh token", err)
	}

	if valid {
		respondWithJSON(w, http.StatusOK, UserCredentials{
			Id:           user.ID.String(),
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
			Email:        user.Email,
			Token:        token,
			RefreshToken: ref_token,
		})
	} else {
		fmt.Println("not valid")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}
}

func (cfg *apiConfig) handleInfoChange(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token", err)
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
	}

	new_hashed_pass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
	}

	new_user, err := cfg.db.UpdateInfo(r.Context(), database.UpdateInfoParams{
		Email:          params.Email,
		HashedPassword: new_hashed_pass,
		ID:             user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update info", err)
	}

	respondWithJSON(w, http.StatusOK, convertUser(new_user))
}
