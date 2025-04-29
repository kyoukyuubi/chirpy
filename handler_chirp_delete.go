package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// get the token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get token", err)
		return
	}

	// validate token
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
	}

	// get the uuid of the chirp
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	// select the chirp
	chirp, err := cfg.dbQueries.GetChirps(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	// validate owner of chirp
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "not autherized to delete chirp, it doesn't belong to you", err)
		return
	}

	// delete chirp
	err = cfg.dbQueries.DeleteChirpFromID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	// respond if succesfull
	respondWithJSON(w, http.StatusNoContent, nil)
}