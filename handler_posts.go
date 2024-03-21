package main

import (
	"net/http"
	"strconv"

	"github.com/jonackers/rssfeeds/internal/database"
)

func (cfg *apiConfig) handlePostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if reqLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = reqLimit
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
						UserID: user.ID,
						Limit: int32(limit),
					})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch posts")
		return
	}

	respondWithJson(w, http.StatusOK, databasePostsToPosts(posts))
}
