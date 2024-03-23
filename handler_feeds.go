package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/database"
)

// handleFeedsCreate handles requests to the create feed endpoint.
// Requests must come from authenticated users, and the feeds are requested using a
// Name and URL provided in the request body.
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
		Feed   Feed       `json:"feed"`
		Follow FeedFollow `json:"feed_follow"`
	}{
		Feed:   databaseFeedToFeed(newFeed),
		Follow: databaseFeedFollowToFeedFollow(newFollow),
	}

	respondWithJson(w, http.StatusOK, response)
}

// handleFeedsGetAll handles requests to the get feeds endpoint.
// It returns all of the feeds in the database.
func (cfg *apiConfig) handleFeedsGetAll(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get feeds")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
