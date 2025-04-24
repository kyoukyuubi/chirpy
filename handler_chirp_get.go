package main

import (
	"net/http"

	"github.com/google/uuid"
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

func (cfg *apiConfig) handlerChirpSelect(w http.ResponseWriter, r *http.Request) {
	// parse the id (string) into an id (uuid)
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	// select the chirp
	chirp, err := cfg.dbQueries.GetChirps(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	// insert into custom struct for better control and send response
	chirpStruct := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_ID: chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirpStruct)
}