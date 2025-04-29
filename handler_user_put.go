package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/chirpy/internal/auth"
	"github.com/kyoukyuubi/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser (w http.ResponseWriter, r *http.Request) {
	// request struct params
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// get the access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	// validate the token
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}

	// get the request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	// has the password
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coulnd't hash password", err)
		return
	}

	// update the user
	user, err := cfg.dbQueries.UpdateUserWithID(r.Context(), database.UpdateUserWithIDParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPass,
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update user", err)
		return
	}

	// Make return struct
	updatedUser := struct {
		User_id uuid.UUID `json:"id"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Email string `json:"email"`
	}{
		User_id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	// respond with the new user minus pass
	respondWithJSON(w, http.StatusOK, updatedUser)
}