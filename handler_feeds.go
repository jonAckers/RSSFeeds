package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/database"
)

func (cfg *apiConfig) handleFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newFeed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Name:      params.Name,
						Url:       params.Url,
						UserID:    user.ID,
					})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create feed")
		return
	}

	newFollow, err := cfg.createFeedFollow(r.Context(), user.ID, newFeed.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not follow new feed")
		return
	}

	response := struct {
		Feed   database.Feed       `json:"feed"`
		Follow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:   newFeed,
		Follow: newFollow,
	}

	respondWithJson(w, http.StatusOK, response)
}


func (cfg *apiConfig) handleFeedsGetAll(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get feeds")
		return
	}

	respondWithJson(w, http.StatusOK, feeds)
}
