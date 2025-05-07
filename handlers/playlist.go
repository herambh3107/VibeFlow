package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"MCP-project/services"
	"MCP-project/utils"

	"github.com/abadojack/whatlanggo"
	"github.com/zmb3/spotify/v2"
)

type MoodRequest struct {
	Mood     string `json:"mood"`
	Language string `json:"language"`
}

type PlaylistResponse struct {
	Tracks []string `json:"tracks"`
}

// Function to check if the track name matches the desired language
func isLanguage(text, desiredLang string) bool {
	// Detect the language of the text
	lang := whatlanggo.Detect(text)
	fmt.Printf("Track: %s, Detected Language: %s, Desired Language: %s\n", text, lang.Lang.String(), desiredLang)

	// Compare strings case-insensitively
	return strings.EqualFold(lang.Lang.String(), desiredLang)
}
func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req MoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Log the incoming request
	fmt.Printf("Received request: Mood=%s, Language=%s\n", req.Mood, req.Language)

	// Validate the provided language
	if req.Language == "" {
		fmt.Println("Language not provided")
		http.Error(w, "Language not provided", http.StatusBadRequest)
		return
	}

	// Check mood-to-genre mapping
	genres, ok := utils.MoodToGenres[req.Mood]
	if !ok {
		fmt.Printf("Mood not recognized: %s\n", req.Mood)
		http.Error(w, "Mood not recognized", http.StatusBadRequest)
		return
	}
	fmt.Printf("Genres for mood %s: %v\n", req.Mood, genres)

	var trackLinks []string

	// Search for tracks by genre
	for _, genre := range genres {
		results, err := services.SpotifyClient.Search(r.Context(), genre, spotify.SearchTypeTrack)
		if err != nil {
			fmt.Printf("Spotify search failed for genre %s: %v\n", genre, err)
			http.Error(w, "Spotify search failed", http.StatusInternalServerError)
			return
		}

		// Log the number of tracks returned
		fmt.Printf("Genre: %s, Tracks found: %d\n", genre, len(results.Tracks.Tracks))

		for i, track := range results.Tracks.Tracks {
			if i >= 3 {
				break
			}

			// Check if the track name matches the desired language
			if isLanguage(track.Name, req.Language) {
				trackLinks = append(trackLinks, track.ExternalURLs["spotify"])
				fmt.Printf("Added track: %s, URL: %s\n", track.Name, track.ExternalURLs["spotify"])
			}
		}
	}

	// Log the final track links
	fmt.Printf("Final track links: %v\n", trackLinks)

	resp := PlaylistResponse{Tracks: trackLinks}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
