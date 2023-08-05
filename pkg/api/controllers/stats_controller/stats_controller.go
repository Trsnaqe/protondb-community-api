package stats_controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/trsnaqe/protondb-api/pkg/services/stats_service"
)

// Endpoint to retrieve stats of the API.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := stats_service.GetStats()
	if err != nil {
		log.Println("Error getting stats:", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
