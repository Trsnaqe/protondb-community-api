package games_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/services/games_service"
)

// Endpoint to retrieve all games.
func GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	/*
		games, err := games_service.GetAllGames()
		if err != nil {
			http.Error(w, "Failed to retrieve games", http.StatusInternalServerError)
			return
		}
		if(len(games) == 0) {
			http.Error(w, "No games found", http.StatusNotFound)
			return
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(games)
		if err != nil {
			http.Error(w, "Failed to encode games", http.StatusInternalServerError)
		}
	*/
	message := "Status Code 503: Service Unavailable.\n\n" +
		"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
		"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(message))
}

// Endpoint to search games by title.
func SearchGameByTitleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("title")
	if query == "" {
		http.Error(w, "Title query parameter is required", http.StatusBadRequest)
		return
	}

	games, err := games_service.SearchGameByTitle(query)
	if err != nil {
		http.Error(w, "Failed to search games by title", http.StatusInternalServerError)
		return
	}
	if len(games) == 0 {
		http.Error(w, "No games found matching the query", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(games)
}

// Endpoint to retrieve a game by gameId.
func GetGameByAppIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	gameID := params["gameId"]

	game, err := games_service.GetGameByAppID(gameID)
	if err != nil {
		if err.Error() == "game not found" {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve game", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(game)
}

func GetGameSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	appID := params["gameId"]

	summary, err := games_service.GetGameSummary(appID)
	if err != nil {
		log.Printf("Error getting game summary: %v", err)
		http.Error(w, "Failed to retrieve game summary", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(summary)
	if err != nil {
		log.Printf("Error encoding game summary response: %v", err)
		http.Error(w, "Failed to encode game summary response", http.StatusInternalServerError)
	}
}

//v2 endpoints

func GetGameByQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var appId, title string

	for key, values := range r.URL.Query() {
		lowerKey := strings.ToLower(key)
		switch lowerKey {
		case "appid", "app_id", "gameid", "game_id":
			appId = strings.ToLower(values[0])
		case "title":
			title = strings.ToLower(values[0])
		}
	}

	if appId == "" && title == "" {
		/*
			games, err := games_service.GetAllGames()
			if err != nil {
				http.Error(w, "Failed to retrieve games", http.StatusInternalServerError)
				return
			}

			encoder := json.NewEncoder(w)
			err = encoder.Encode(games)
			if err != nil {
				http.Error(w, "Failed to encode games", http.StatusInternalServerError)
			}
		*/
		message := "Status Code 503: Service Unavailable.\n\n" +
			"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
			"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(message))
		return
	}

	games, err := games_service.GetGameByQuery(appId, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(games) == 0 {
		http.Error(w, "No games found matching the query", http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(games)
	if err != nil {
		http.Error(w, "Failed to encode games", http.StatusInternalServerError)
	}
}
