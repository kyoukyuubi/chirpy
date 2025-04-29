package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kyoukyuubi/chirpy/internal/auth"
	"github.com/kyoukyuubi/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetRefreshToken(w http.ResponseWriter, r *http.Request) {
	// get token from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found in header", err)
		return
	}

	// check if token is in the database
	tokenData, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "code not found in database", err)
		return
	}

	// check if token is expired
	if tokenData.RevokedAt.Valid{
		respondWithError(w, http.StatusUnauthorized, "code invalid", err)
		return
	}
	if time.Now().After(tokenData.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "code invalid", err)
		return
	}

	// generate access token
	tokenJWT, err := auth.MakeJWT(tokenData.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate JWT", err)
		return
	}

	// return the token
	returnStruct := struct {
		JWT string `json:"token"`
	}{
		JWT: tokenJWT,
	}
	respondWithJSON(w, http.StatusOK, returnStruct)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get token from header", err)
	}

	err = cfg.dbQueries.RevokeToken(r.Context(), database.RevokeTokenParams{
		RevokedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token: token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil) 
}