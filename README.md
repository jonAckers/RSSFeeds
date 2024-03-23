# RSSFeeds

## Description

This project contains a webserver that allows users to subscribe to RSSFeeds, and continuously checks for new posts on these feeds and publishing them to users.

## Features

- Sign-up with an username.
- Make requests using an API key.
- Subscribe to RSS feeds.
- Get notified of posts from the subscribed to feeds.

## Technologies Used

- Go
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev) for automatic SQL/Go interface generation

## Installation and Usage

1. Clone the repository: `git clone git@github.com:jonAckers/RSSFeeds.git`
2. Build the project: `go build`
3. Run the project: `./rssfeeds`
