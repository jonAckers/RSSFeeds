package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jonackers/rssfeeds/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", handleReadiness)
	mux.HandleFunc("GET /v1/err", handleError)

	mux.HandleFunc("POST /v1/users", apiCfg.handleUsersCreate)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handleUsersGetByApiKey))

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handleFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handleFeedsGetAll)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handleFeedFollowsCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleFeedFollowsDelete))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handleFeedFollowsGetAll))

	mux.HandleFunc("GET /v1/posts", apiCfg.middlewareAuth(apiCfg.handlePostsGet))

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	const collectionConcurrency = 10
	go startScraping(dbQueries, collectionConcurrency, time.Minute)

	log.Println("Server is listening on port" + port + "...")
	log.Fatal(srv.ListenAndServe())
}
