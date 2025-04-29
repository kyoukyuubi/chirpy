package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	// set the expected struct
	type paramters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// get the request
	decoder := json.NewDecoder(r.Body)
	params := paramters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode params", err)
		return
	}

	// check for correct event
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// update the user
	_, err = cfg.dbQueries.UpgradeUser(r.Context(), database.UpgradeUserParams{
		IsChirpyRed: true,
		ID: params.Data.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// respond with 204 on success
	w.WriteHeader(http.StatusNoContent)
}