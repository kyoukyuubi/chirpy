package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpSelectAll(w http.ResponseWriter, r *http.Request) {
	// get all chirps, handle errors
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database", err)
		return
	}

	// construct a slice of structs for better control
	returnChirp := []Chirp{}
	for _, chirp := range chirps {
		returnChirp = append(returnChirp, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_ID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, returnChirp)
}