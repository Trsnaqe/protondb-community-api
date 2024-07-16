package reports_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/constants"
	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/services/reports_service"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

// Endpoint to retrieve reports.
func GetReportsHandler(w http.ResponseWriter, r *http.Request) {

	//replace with function below to activate GetStreamOfReports
	//GetStreamOfReports(w, r)

	message := "Status Code 503: Service Unavailable.\n\n" +
		"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
		"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(message))
}

func GetStreamOfReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cursor, err := storage.GetAllReports()
	if err != nil {
		http.Error(w, "Failed to retrieve reports", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	encoder := json.NewEncoder(w)
	w.Write([]byte("["))

	first := true
	var report models.Report
	for cursor.Next(r.Context()) {
		if err := cursor.Decode(&report); err != nil {
			http.Error(w, "Failed to decode report", http.StatusInternalServerError)
			return
		}

		if !first {
			w.Write([]byte(","))
		}
		first = false

		versioned := r.URL.Query().Get("versioned")
		if versioned == "true" || versioned == "1" {
			if err := encoder.Encode(report); err != nil {
				http.Error(w, "Failed to encode report", http.StatusInternalServerError)
				return
			}
		} else {
			if err := encoder.Encode(report.Data); err != nil {
				http.Error(w, "Failed to encode report data", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Write([]byte("]"))

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// v1 implementation, so no version filtering support
func GetReportsByGameIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	gameID := params["gameId"]

	reports, err := storage.GetReportsByGameID(gameID, "")
	if err != nil {
		if err.Error() == "game not found" {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve reports: %v", err), http.StatusInternalServerError)
		return
	}
	if len(reports) == 0 {
		http.Error(w, "No reports found for the game", http.StatusNotFound)
		return
	}

	versioned := strings.ToLower(r.URL.Query().Get("versioned"))
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

// Endpoint to retrieve reports by gameId.
func GetReportsByQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reports []interface{}
	var err error

	var appId, version, title string
	var versioned bool
	precision := constants.DEFAULT_SEARCH_PRECISION

	queryParams := r.URL.Query()
	for key, values := range queryParams {
		lowerKey := strings.ToLower(key)
		switch lowerKey {
		case "appid", "app_id", "gameid", "game_id":
			appId = strings.ToLower(values[0])
		case "title":
			title = strings.ToLower(values[0])
		case "versioned":
			versioned = values[0] == "true" || values[0] == "1"
		case "version":
			switch values[0] {
			case "1":
				version = "V1"
			case "2":
				version = "V2"
			}
		case "precision":
			parsedPrecision, err := strconv.ParseFloat(values[0], 32)
			if err != nil {
				http.Error(w, "Invalid precision value", http.StatusBadRequest)
				return
			}
			if parsedPrecision < 0 || parsedPrecision > 2 {
				http.Error(w, "Precision value must be between 0 and 2", http.StatusBadRequest)
				return
			}
			precision = float64(parsedPrecision)
		}
	}
	if appId != "" {
		reports, err = reports_service.GetReportsByGameID(appId, versioned, version)
		if err != nil {
			if err.Error() == "game not found" {
				http.Error(w, "Game not found", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to retrieve reports: %v", err), http.StatusInternalServerError)
			return
		}
	} else if title != "" {

		reports, err = reports_service.GetReportsByTitleSearch(title, versioned, version, precision)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve reports: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		/*
			to activate de-comment this function and remove the code below
				GetStreamOfReports(w, r)
		*/
		message := "Status Code 503: Service Unavailable.\n\n" +
			"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
			"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"

		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(message))
		return
	}

	if len(reports) == 0 {
		http.Error(w, "No reports found matching the query", http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(reports)
	if err != nil {
		http.Error(w, "Failed to encode reports", http.StatusInternalServerError)
	}
}
