package stats_service

import (
	"github.com/trsnaqe/protondb-api/pkg/services/background_services"
	"github.com/trsnaqe/protondb-api/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStats() (map[string]interface{}, error) {
	totalGameCount, err := storage.GetTotalGamesCount()
	if err != nil {
		return nil, err
	}

	totalReportsCount, err := storage.GetTotalReportsCount()
	if err != nil {
		return nil, err
	}

	lastProcessedData, err := storage.GetLastProcessedData()
	if err != nil {
		return nil, err
	}
	lastFile := ""
	lastDate := primitive.DateTime(0)
	if lastProcessedData != nil {
		lastFile = lastProcessedData.LastProcessedFile
		lastDate = lastProcessedData.LastProcessedTime
	}

	updateInterval := background_services.GetUpdateInterval()

	// Calculate the time remaining for the next update
	timeRemaining := background_services.TimeRemainingForNextUpdate(updateInterval)

	// Convert the time remaining into a human-readable format
	totalTimeRemainingStr := background_services.TimeInDaysFormat(timeRemaining)

	stats := map[string]interface{}{
		"totalGameCount":            totalGameCount,
		"totalReportsCount":         totalReportsCount,
		"lastProcessedFile":         lastFile,
		"lastProcessedDate":         lastDate,
		"timeRemainingToNextUpdate": totalTimeRemainingStr,
	}

	return stats, nil
}
