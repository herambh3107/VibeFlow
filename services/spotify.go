package services

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

var SpotifyClient *spotify.Client

func InitSpotify() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Set up client credentials for Spotify OAuth2
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     "https://accounts.spotify.com/api/token", // Direct URL for Spotify token
	}

	// Get token using client credentials flow
	_, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("Error getting Spotify token: %v", err)
	}

	// Use the token to create a Spotify client
	httpClient := config.Client(context.Background())
	SpotifyClient = spotify.New(httpClient)
}
