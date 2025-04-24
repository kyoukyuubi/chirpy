package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kyoukyuubi/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_ID   uuid.UUID `json:"user_id"`
}


func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	// request struct params
	type parameters struct {
		Body string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	// get request and decode it, handling errors
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	// validate the body
	cleaned, err := chirpsValidate(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't validate chirp", err)
		return
	}

	// insert chirp into database
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body: cleaned,
		UserID: params.User_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// respond with the nerly created chirp
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body: chirp.Body,
		User_ID: chirp.UserID,
	})
}

func chirpsValidate(body string) (string, error) {
const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", fmt.Errorf("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)

	return cleaned, nil
}

func getCleanedBody(str string, badWords map[string]struct{}) string {
	words := strings.Split(str, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}