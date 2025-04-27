package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
		Password string
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	// decode the response and handle errors
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// make expires_in_seconds optional
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds < 3600{
		params.ExpiresInSeconds = 3600
	}
	
	// select from the database, handle the errors
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user from database", err)
		return
	}

	// check if the password matches
	err = auth.CheckPasswordHash(user.HashedPassword.String, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid login", err)
		return
	}

	// generate a JWT token
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(params.ExpiresInSeconds) * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate token", err)
		return
	}

	// respond with the user on success
	userStruct := struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, userStruct)
}