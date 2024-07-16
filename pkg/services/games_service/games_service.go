package games_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

func GetAllGames() ([]models.Game, error) {
	cursor, err := storage.GetAllGames()
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var games []models.Game
	for cursor.Next(nil) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func SearchGameByTitle(title string, precision float64) ([]models.Game, error) {
	cursor, err := storage.SearchGameByTitle(title, precision)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var games []models.Game
	for cursor.Next(nil) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func GetGameByAppID(gameID string) (*models.Game, error) {
	game, err := storage.GetGameByAppID(gameID)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func GetGameByQuery(appId string, title string, precision float64) ([]models.Game, error) {
	ctx := context.TODO()

	if appId != "" {
		appId = strings.ToLower(appId)
		game, err := storage.GetGameByAppID(appId)
		if err == nil {
			return []models.Game{*game}, nil
		}
		return nil, fmt.Errorf("no game found with appId: %s", appId)
	}

	if title != "" {
		title = strings.ToLower(title)
		if len(title) < 5 {
			return nil, fmt.Errorf("title parameter must be at least 5 characters long")
		}
		cursor, err := storage.SearchGameByTitle(title, precision)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var games []models.Game
		for cursor.Next(ctx) {
			var game models.Game
			if err := cursor.Decode(&game); err != nil {
				return nil, err
			}
			games = append(games, game)
		}
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return games, nil
	}

	return nil, fmt.Errorf("no valid query parameters provided")
}
func AddReportToGame(game *models.Game, report *models.Report) error {
	return storage.InsertReport(game, report.ID)
}

func GetGameSummary(appID string) (*models.GameSummary, error) {
	apiURL := fmt.Sprintf("https://www.protondb.com/api/v1/reports/summaries/%s.json", appID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch game summary. Status code: %d", resp.StatusCode)
	}

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

		err := storage.ChangeTitle(game.ID, *title)
		if err != nil {
			log.Println("Error updating game:", err)
			return err
		}
	}
	return nil
}
