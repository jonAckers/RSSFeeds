package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeader = errors.New("no authorisation error included")
var ErrMalformedToken = errors.New("malformed authoriation header")

// GetApiKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", ErrMalformedToken
	}

	return splitAuth[1], nil
}

