package games_controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/services/games_service"
)

// Endpoint to retrieve all games.
func GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	//games_service.GetAllGamesHandler(w, r)
	message := "Status Code 503: Service Unavailable.\n\n" +
		"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
		"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(message))

}

// Endpoint to retrieve a game by gameId.
func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	games_service.GetGameByIDHandler(w, r)
}

func GetGameSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the game ID from the request URL parameters
	params := mux.Vars(r)
	appID := params["gameId"]

	// Fetch the game summary data using the game service
	summary, err := games_service.GetGameSummary(appID)
	if err != nil {
		log.Printf("Error getting game summary: %v", err)
		http.Error(w, "Failed to retrieve game summary", http.StatusInternalServerError)
		return
	}

	// Set the response Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the summary data and send it in the response
	err = json.NewEncoder(w).Encode(summary)
	if err != nil {
		log.Printf("Error encoding game summary response: %v", err)
		http.Error(w, "Failed to encode game summary response", http.StatusInternalServerError)
		return
	}
}
