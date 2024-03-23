package main

import "net/http"

// handleError handles a verification endpoint to test
// if the serve reports errors correctly.
func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "An error occurred")
}
