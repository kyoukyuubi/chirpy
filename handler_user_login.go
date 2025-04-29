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

	// generate a JWT token
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate token", err)
		return
	}

	// generate the refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot generate refresh token", err)
		return
	}

	// store the token in the database
	err = cfg.dbQueries.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams{
		Token: refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't insert refresh token", err)
		return
	}

	// respond with the user on success
	userStruct := struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
		Token: token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, userStruct)
}