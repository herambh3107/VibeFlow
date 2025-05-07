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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatalf("Missing SPOTIFY_ID or SPOTIFY_SECRET in environment")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://accounts.spotify.com/api/token",
	}

	_, err = config.Token(context.Background())
	if err != nil {
		log.Fatalf("Spotify token error: %v", err)
	}

	httpClient := config.Client(context.Background())
	SpotifyClient = spotify.New(httpClient)

	log.Println("✅ Spotify client initialized successfully")
}
