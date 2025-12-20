package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/azmtbek/chirpy/internal/auth"
	"github.com/google/uuid"
)

type UserWithToken struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		UserWithToken
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	expirationTime := getExpirationTime(params.ExpiresInSeconds)

	jwt, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couln't generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		UserWithToken: UserWithToken{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     jwt,
		},
	})
}

func getExpirationTime(userInput int) time.Duration {
	expiration := 60 * 60 // 1 hour
	if userInput != 0 {
		expiration = min(userInput, expiration)
	}

	return time.Duration(expiration) * time.Second
}
