package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/database"
)

// createFeedFollow adds a new entry to the Feed-Follow
// table of the database, using the provided user and feed id.
func (cfg *apiConfig) createFeedFollow(ctx context.Context, userId, feedId uuid.UUID) (database.FeedFollow, error) {
	return cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		FeedID:    feedId,
	})
}

// handleFeedFollowsCreate handles requests to the Create Feed Follows endpoint.
// Feed follows are created using the authenticated user id, and a feed id
// provided in the request body.
func (cfg *apiConfig) handleFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newFeedFollow, err := cfg.createFeedFollow(r.Context(), user.ID, uuid.MustParse(params.FeedID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not follow feed")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedFollowToFeedFollow(newFeedFollow))
}

// handleFeedFollowsDelete handles requests to the Delete Feed Follows endpoint.
// The feed follow to delete is indicated by the authenticated user id and the
// feed follow id provided in the URL.
func (cfg *apiConfig) handleFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedId := r.PathValue("feedFollowID")

	err := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     uuid.MustParse(feedId),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not unfollow feed")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleFeedFollowsGetAll handles requests to the Get Feed Follows endpoint.
// It returns all feed follows matching the authenticated user id.
func (cfg *apiConfig) handleFeedFollowsGetAll(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch feed follows")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}
