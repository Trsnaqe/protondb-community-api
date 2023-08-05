package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/trsnaqe/protondb-api/pkg/server"
	"github.com/trsnaqe/protondb-api/pkg/services/background_services"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

var updateInterval = 24 * 30 * time.Hour // Change this to your desired update interval

func main() {
	if os.Getenv("PORT") == "" {
		// Load .env file if not running on Heroku
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}
	err := storage.ConnectDB()
	if err != nil {
		panic(err)
	}
	background_services.SetUpdateInterval(updateInterval)
	if updateInterval > 0 { // Don't start the goroutine if updateInterval is zero
		go background_services.ProcessReportsBackground(updateInterval) // Pass the update interval to the background process
	}
	defer storage.CloseDB()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := server.NewServer()
	server.Run(":" + port)
}
