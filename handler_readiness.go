package main

import "net/http"

// handleReadiness handles a verification endpoint to
// indicate if the server is running.
func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}
