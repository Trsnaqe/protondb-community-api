package api

import (
	"github.com/gorilla/mux"
	gamesCtrl "github.com/trsnaqe/protondb-api/pkg/api/controllers/games_controller"
	infoCtrl "github.com/trsnaqe/protondb-api/pkg/api/controllers/info_controller"
	reportsCtrl "github.com/trsnaqe/protondb-api/pkg/api/controllers/reports_controller"
	statsCtrl "github.com/trsnaqe/protondb-api/pkg/api/controllers/stats_controller"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", infoCtrl.ListAPIEndpointsHandler).Methods("GET")
	r.HandleFunc("/api", infoCtrl.ListAPIEndpointsHandler).Methods("GET")
	r.HandleFunc("/api/games", gamesCtrl.GetAllGamesHandler).Methods("GET")
	r.HandleFunc("/api/games/{gameId}", gamesCtrl.GetGameByAppIDHandler).Methods("GET")
	r.HandleFunc("/api/games/{gameId}/summary", gamesCtrl.GetGameSummaryHandler).Methods("GET")
	r.HandleFunc("/api/reports", reportsCtrl.GetReportsHandler).Methods("GET")
	r.HandleFunc("/api/reports/{gameId}", reportsCtrl.GetReportsByGameIDHandler).Methods("GET")
	r.HandleFunc("/api/stats", statsCtrl.StatsHandler).Methods("GET")

	r.HandleFunc("/api/v2/games", gamesCtrl.GetGameByQueryHandler).Methods("GET")
	r.HandleFunc("/api/v2/reports", reportsCtrl.GetReportsByQueryHandler).Methods("GET")
}
