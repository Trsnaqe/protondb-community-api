package reports_service

import (
	"log"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/services/games_service"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

func GetReports(versioned bool) ([]interface{}, error) {
	cursor, err := storage.GetAllReports()
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var reports []interface{}
	for cursor.Next(nil) {
		var report models.Report
		if err := cursor.Decode(&report); err != nil {
			return nil, err
		}

		if versioned {
			reports = append(reports, report)
		} else {
			reports = append(reports, report.Data)
		}
	}

	return reports, nil
}

func GetReportsByGameID(gameID string, versioned bool, version string) ([]interface{}, error) {
	reports, err := storage.GetReportsByGameID(gameID, version)
	if err != nil {
		return nil, err
	}
	var data []interface{}
	if versioned {
		for _, report := range reports {
			data = append(data, report)
		}
	} else {
		for _, report := range reports {
			data = append(data, report.Data)
		}
	}

	return data, nil
}

// search by title and get its reports with versioned version etc
func GetReportsByTitleSearch(title string, versioned bool, version string) ([]interface{}, error) {
	games, err := games_service.SearchGameByTitle(title)
	if err != nil {
		return nil, err
	}

	var reports []interface{}
	for _, game := range games {
		gameReports, err := GetReportsByGameID(game.AppID, versioned, version)
		if err != nil {
			return nil, err
		}
		reports = append(reports, gameReports...)
	}

	return reports, nil
}

func CreateNewReport(report map[string]interface{}, game *models.Game, reportVersion string) error {
	newReport := &models.Report{
		Data:          report,
		ReportVersion: reportVersion,
	}
	createdReport, err := storage.CreateReport(newReport, game)
	if err != nil {
		log.Println("Error creating report:", err)
		return err
	}
	err = games_service.AddReportToGame(game, createdReport)
	if err != nil {
		log.Println("Error adding report to game:", err)
		return err
	}

	return nil
}
