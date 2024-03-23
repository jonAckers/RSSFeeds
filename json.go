package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithJson writes a given JSON object to the http response body
// as well as a given http code to the header.
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError writes an error message to the  http response body
// as well as a given http code to the header.
func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with error:", msg)
	}
	respondWithJson(w, code, map[string]string{"error": msg})
}
