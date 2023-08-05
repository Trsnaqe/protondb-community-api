package background_services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"github.com/trsnaqe/protondb-api/pkg/services/games_service"
	"github.com/trsnaqe/protondb-api/pkg/services/reports_service"
	"github.com/trsnaqe/protondb-api/pkg/storage"
)

func processReports(reports interface{}) error {
	switch reports := reports.(type) {
	case []models.ReportFormatV1:
		for _, report := range reports {
			err := processReport(report, report.AppID, report.Title, "V1")
			if err != nil {
				return err
			}
		}
	case []models.ReportFormatV2:
		for _, report := range reports {
			err := processReport(report, fmt.Sprint(report.App.Steam.AppID), report.App.Title, "V2")
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid report type")
	}
	return nil
}

func processReport(report interface{}, appID, title, version string) error {
	reportMap := make(map[string]interface{})
	j, _ := json.Marshal(report)
	json.Unmarshal(j, &reportMap)

	if isScientificNotation(appID) {
		appIDInt, err := convertScientificNotation(appID)
		if err != nil {
			return err
		}
		appID = fmt.Sprint(appIDInt)
	}
	game, err := games_service.GetOrCreateGame(appID, &title)
	if err != nil {
		return err
	}

	err = games_service.UpdateGameTitle(game, &title)
	if err != nil {
		return err
	}

	// If report is of V2 type, compare it before insertion
	if version == "V2" {
		v2report, ok := report.(models.ReportFormatV2)
		if ok {
			if storage.CompareReport(v2report) {
				log.Println("Report already exists in the database. Skipping...")
				return nil
			}
		} else {
			return fmt.Errorf("error casting report to V2 type")
		}
	}

	err = reports_service.CreateNewReport(reportMap, game, version) // Pass "v1" or "v2" as the ReportVersion
	if err != nil {
		return err
	}

	return nil
}

func ProcessReportFile(file []byte) error {
	log.Println("Starting to process report file...")

	var v1Reports []models.ReportFormatV1
	var v2Reports []models.ReportFormatV2

	err := json.Unmarshal(file, &v2Reports)
	if err != nil {
		log.Println("Error while unmarshalling JSON into ReportFormatV2:", err)
		err = json.Unmarshal(file, &v1Reports)
		if err != nil {
			log.Println("Error while unmarshalling JSON into ReportFormatV1:", err)
			return err
		}
		err = processReports(v1Reports)
	} else {
		// Check if there is any existing V2 data in the database
		v2Count, err := storage.CountV2Reports()
		if err != nil {
			log.Println("Error while checking V2 report count in the database:", err)
			return err
		}

		if v2Count > 0 {
			// Process V2 reports in reverse
			log.Println("Processing V2 reports in reverse...")
			for i := len(v2Reports) - 1; i >= 0; i-- {
				err = processReports(v2Reports[i : i+1]) // Process one report at a time
				if err != nil {
					return err
				}
			}
		} else {
			// Process V2 reports in normal order
			log.Println("Processing V2 reports in normal...")
			err = processReports(v2Reports)
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	log.Println("Finished processing report file.")
	return nil
}
