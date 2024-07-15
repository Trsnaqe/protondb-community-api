package background_services

import (
	"log"
	"os"
	"time"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessReportsBackground(updateInterval time.Duration) {
	processStatus, err := storage.GetLastProcessStatus()
	if err != nil {
		log.Fatalf("Failed to get process status: %v", err)
	}

	if processStatus == nil {
		processStatus = &models.ProcessStatus{
			LastProcessedFile: "reports_oct31_2019.tar.gz",
			LastProcessedTime: primitive.NewDateTimeFromTime(time.Now()),
		}
		err = storage.CreateProcessStatus(processStatus)
		if err != nil {
			log.Fatalf("Failed to create process status: %v", err)
		}
	}

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		jsonData, newLastProcessedFile, extractedDir, err := GetLatestProcessedReportFile(processStatus.LastProcessedFile)
		if err != nil {
			log.Fatal(err)
		}

		if newLastProcessedFile != processStatus.LastProcessedFile {
			err = ProcessReportFile(jsonData)
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(extractedDir)
			if err != nil {
				log.Println("Error removing directory:", err)
			}

			// Update the last processed file in the database
			processStatus.LastProcessedFile = newLastProcessedFile
			processStatus.LastProcessedTime = primitive.NewDateTimeFromTime(time.Now())
			err = storage.UpdateProcessStatus(processStatus)
			if err != nil {
				log.Fatalf("Failed to update process status: %v", err)
			}

			SetLastTickTime(time.Now())

		}
		if processStatus.LastProcessedFile == newLastProcessedFile {
			SetLastTickTime(time.Now())
		}

	}
}
