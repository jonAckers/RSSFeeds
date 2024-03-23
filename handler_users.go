package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/database"
)

// handleUsersCreate handles create requests to the users endpoint.
// the name of the user to create is provided in the request body.
func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJson(w, http.StatusOK, databaseUserToUser(newUser))
}

// handleUsersGetByApiKey returns the instance of the user
// that was authenticated using the api key in the request header.
func (cfg *apiConfig) handleUsersGetByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, databaseUserToUser(user))
}
