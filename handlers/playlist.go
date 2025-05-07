package handlers

import (
	"encoding/json"
	"net/http"

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

	// Check if the detected language matches the desired language
	return lang.Lang.String() == desiredLang
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req MoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the provided language
	if req.Language == "" {
		http.Error(w, "Language not provided", http.StatusBadRequest)
		return
	}

	genres, ok := utils.MoodToGenres[req.Mood]
	if !ok {
		http.Error(w, "Mood not recognized", http.StatusBadRequest)
		return
	}

	var trackLinks []string

	// Search for tracks by genre
	for _, genre := range genres {
		results, err := services.SpotifyClient.Search(r.Context(), genre, spotify.SearchTypeTrack)
		if err != nil {
			http.Error(w, "Spotify search failed", http.StatusInternalServerError)
			return
		}

		for i, track := range results.Tracks.Tracks {
			if i >= 3 {
				break
			}

			// Check if the track name matches the desired language
			if isLanguage(track.Name, req.Language) {
				trackLinks = append(trackLinks, track.ExternalURLs["spotify"])
			}
		}
	}

	resp := PlaylistResponse{Tracks: trackLinks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
