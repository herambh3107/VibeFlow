package main

import (
	"log"
	"net/http"

	"MCP-project/handlers"
	"MCP-project/services"
)

func main() {
	services.InitSpotify()

	http.HandleFunc("/get-playlist", handlers.GetPlaylist)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
