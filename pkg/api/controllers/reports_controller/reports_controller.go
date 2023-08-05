package reports_controller

import (
	"net/http"

	"github.com/trsnaqe/protondb-api/pkg/services/reports_service"
)

// Endpoint to retrieve reports.
func GetReportsHandler(w http.ResponseWriter, r *http.Request) {
	//reports_service.GetStreamOfReports(w, r)
	message := "Status Code 503: Service Unavailable.\n\n" +
		"This endpoint is currently unavailable as the server cannot handle the request. Please try again later or consider supporting the project by buying me a coffee. Your support helps keep this service running.\n\n" +
		"Support the project by buying me a coffee at: https://www.buymeacoffee.com/trsnaqe"

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(message))
}

// Endpoint to retrieve reports by gameId.
func GetReportsByGameIDHandler(w http.ResponseWriter, r *http.Request) {
	reports_service.GetReportsByGameIDHandler(w, r)

}
