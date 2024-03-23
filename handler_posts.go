package main

import (
	"net/http"
	"strconv"

	"github.com/jonackers/rssfeeds/internal/database"
)

// handlePostsGet handles get requests to the posts endpoint.
// It returns all of the posts relevant to the authenticated user.
// The number of posts returned can be selected in the query.
// Otherwise at most 10 posts will be returned.
func (cfg *apiConfig) handlePostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if reqLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = reqLimit
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch posts")
		return
	}

	respondWithJson(w, http.StatusOK, databasePostsToPosts(posts))
}
