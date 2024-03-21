package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jonackers/rssfeeds/internal/database"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}


func fetchFeed(feedUrl string) (*RSSFeed, error) {
	client := http.Client {
					Timeout: 10 * time.Second,
				}

	resp, err := client.Get(feedUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}


func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Failed to mark feed '%s' as fetched: %v", feed.Name, err)
		return
	}

	data, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed '%s': %v", feed.Name, err)
	}

	for _, item := range data.Channel.Item {
		log.Printf("Found post: %s", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(data.Channel.Item))
}
