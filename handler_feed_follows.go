package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonackers/rssfeeds/internal/database"
)


func (cfg *apiConfig) createFeedFollow(ctx context.Context, userId, feedId uuid.UUID) (database.FeedFollow, error) {
	return cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				UserID: userId,
				FeedID: feedId,
			})
}


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


func (cfg *apiConfig) handleGetAllFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch feed follows")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}
