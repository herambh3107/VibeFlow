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
	_ = godotenv.Load()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     "https://accounts.spotify.com/api/token",
	}

	_, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("Spotify token error: %v", err)
	}

	httpClient := config.Client(context.Background())
	SpotifyClient = spotify.New(httpClient)
}
