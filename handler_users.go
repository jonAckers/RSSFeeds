package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/auth"
	"github.com/jonackers/rssfeeds/internal/database"
)


func (cfg apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
						ID: 	   uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Name:      params.Name,
					})

	if err != nil {
		respondWithError(w, 500, "Could not create user")
		return
	}

	respondWithJson(w, 200, newUser)
}


func (cfg apiConfig) handleUsersGetByApiKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User not found")
		return
	}

	respondWithJson(w, 200, user)
}
