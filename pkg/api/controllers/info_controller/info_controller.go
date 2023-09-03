package info_controller

import "net/http"

// listAPIEndpointsHandler lists all available endpoints on localhost:8080/api
func ListAPIEndpointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	endpoints := []string{
		"/api/games (GET): Get all games*",
		"/api/games/{gameId} (GET): Get a game by gameId",
		"/api/games/{gameId}/summary (GET): Get tiers by gameId, fetched from protondb directly",
		"/api/reports (GET): Retrieve reports, add ?versioned=true for versioned data*",
		"/api/reports/{gameId} (GET): Get reports by gameId, add ?versioned=true for versioned data",
		"/api/stats (GET): Get stats of the API",
	}
	response := "Available endpoints in the protondb.solidet.com:\n\n"

	for _, endpoint := range endpoints {
		response += endpoint + "\n"
	}

	response += "\n*Had to disable the /api/games and /api/reports endpoints because the dataset is large and it costs a lot to leave those endpoints open. I want to enable them again in the future, but I need to find support to help me host this.\n\n"

	response += "As of 04/09/2023 I had to disable the deployment and replaced it with a sample database due to budget reasons; please see github for further information. \n\n"
	
	openSourceLink := "You can find the source code for this project on GitHub:\nhttps://github.com/Trsnaqe/protondb-community-api\n\n"

	supportMessage := "If you find this service useful and would like to support the project, consider buying me a coffee. Your support helps cover server costs and keeps this service running:\nhttps://www.buymeacoffee.com/trsnaqe\n\n"

	response += openSourceLink + supportMessage

	w.Write([]byte(response))
}
