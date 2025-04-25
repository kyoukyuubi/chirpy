package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoukyuubi/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
		Password string
	}

	// decode the response and handle errors
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
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

	// respond with the user on success
	userStruct := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, userStruct)
}