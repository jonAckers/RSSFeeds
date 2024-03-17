package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", handleReadiness)
	mux.HandleFunc("GET /v1/err", handleError)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Println("Server is listening on port" + port + "...")
	log.Fatal(srv.ListenAndServe())
}
