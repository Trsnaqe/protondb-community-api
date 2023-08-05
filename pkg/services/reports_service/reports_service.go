package reports_service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/services/games_service"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

func GetStreamOfReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cursor, err := storage.GetAllReports()
	if err != nil {
		http.Error(w, "Failed to retrieve reports", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	encoder := json.NewEncoder(w)
	w.Write([]byte("[")) // Start of JSON array

	first := true
	var report models.Report
	for cursor.Next(r.Context()) {
		if err := cursor.Decode(&report); err != nil {
			http.Error(w, "Failed to decode report", http.StatusInternalServerError)
			return
		}

		if !first {
			w.Write([]byte(",")) // Add a comma before each report, except the first one
		}
		first = false

		// Check if versioned is required
		versioned := r.URL.Query().Get("versioned")
		if versioned == "true" || versioned == "1" {
			// If it does, encode the entire report
			if err := encoder.Encode(report); err != nil {
				http.Error(w, "Failed to encode report", http.StatusInternalServerError)
				return
			}
		} else {
			// If it does not, encode only the "Data" part of the report
			if err := encoder.Encode(report.Data); err != nil {
				http.Error(w, "Failed to encode report data", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Write([]byte("]")) // End of JSON array
}

func GetReportsByGameIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	gameID := params["gameId"]

	reports, err := storage.GetReportsByGameID(gameID)
	if err != nil {
		if err.Error() == "game not found" {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve reports: %v", err), http.StatusInternalServerError)
		return
	}

	versioned := r.URL.Query().Get("versioned")
	if versioned == "true" || versioned == "1" {
		json.NewEncoder(w).Encode(reports)
	} else {
		var data []interface{}
		for _, report := range reports {
			data = append(data, report.Data)
		}
		json.NewEncoder(w).Encode(data)
	}
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
