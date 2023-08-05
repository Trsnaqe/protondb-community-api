package games_service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

func GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cursor, err := storage.GetAllGames()
	if err != nil {
		http.Error(w, "Failed to retrieve games", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	encoder := json.NewEncoder(w)
	w.Write([]byte("["))

	first := true
	var game models.Game
	for cursor.Next(r.Context()) {
		if err := cursor.Decode(&game); err != nil {
			http.Error(w, "Failed to decode game", http.StatusInternalServerError)
			return
		}

		if !first {
			w.Write([]byte(","))
		}
		first = false

		if err := encoder.Encode(game); err != nil {
			http.Error(w, "Failed to encode game", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("]"))
}

func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	gameID := params["gameId"]

	// Get the game by its ID
	game, err := storage.GetGameByAppID(gameID)
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

func AddReportToGame(game *models.Game, report *models.Report) error {
	err := storage.InsertReport(game, report.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetGameSummary returns the game summary data for the given appId.
func GetGameSummary(appID string) (*models.GameSummary, error) {
	apiURL := fmt.Sprintf("https://www.protondb.com/api/v1/reports/summaries/%s.json", appID)
	// Make the HTTP GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch game summary. Status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into the GameSummary struct
	var summary models.GameSummary
	err = json.NewDecoder(resp.Body).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
func GetOrCreateGame(appID string, title *string) (*models.Game, error) {
	game, err := storage.GetGameByAppID(appID)
	if err != nil {
		log.Println("Error getting game by app ID:", err)
		return nil, err
	}

	if game == nil {
		log.Println("No existing game found. Creating a new one...")
		game = models.NewGame(appID, title)
		err = storage.CreateGame(game)
		if err != nil {
			log.Println("Error creating game:", err)
			return nil, err
		}
	}

	return game, nil
}
func UpdateGameTitle(game *models.Game, title *string) error {
	if game.Title == nil || *game.Title != *title {
		game.Title = title

		// Dereference the pointer and pass the string value to ChangeTitle
		err := storage.ChangeTitle(game.ID, *title)
		if err != nil {
			log.Println("Error updating game:", err)
			return err
		}
	}
	return nil
}
