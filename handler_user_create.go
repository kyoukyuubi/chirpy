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

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
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

	// has the password, handling the errors
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
		return
	}

	// input the user into the database and handle errors
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hash,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert user into database", err)
		return
	}

	// respond to the POST request with the newly created user
	userStruct := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, userStruct)
}